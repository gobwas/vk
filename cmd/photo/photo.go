package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gobwas/vk"
	"github.com/gobwas/vk/cli"
	"github.com/gobwas/vk/internal/httputil"
	"github.com/gobwas/vk/internal/syncutil"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func getDefaultDest(suffix string) (path string) {
	user, err := user.Current()
	if err == nil {
		path = user.HomeDir
	}
	if path == "" {
		path = "/tmp"
	}
	return path + "/" + suffix
}

var (
	clientID = flag.String(
		"client_id", "",
		"application id",
	)
	clientSecret = flag.String(
		"client_secret", "",
		"application secret",
	)
	ownerID = flag.String(
		"owner_id", "",
		"albums owner id (empty for your id)",
	)
	dest = flag.String(
		"dest", getDefaultDest("vkphoto"),
		"destination root dir for photos",
	)
	parallelism = flag.Int(
		"parallelism", 32,
		"number of parallel downloadings",
	)
)

func main() {
	flag.Parse()

	if *clientID == "" || *clientSecret == "" {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options]\n\n",
			os.Args[0],
		)
		flag.CommandLine.SetOutput(os.Stderr)
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()

	app := vk.App{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Scope:        vk.ScopePhotos,
	}

	access, err := cli.Authorize(ctx, app)
	if err != nil {
		log.Fatal(err)
	}

	if *ownerID == "" {
		*ownerID = strconv.Itoa(access.UserID)
	}
	*dest += "/" + *ownerID

	albums, err := getAlbums(ctx, access, *ownerID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "ready to store photos at %s\n", *dest)

	bars := sync.Map{}
	progress := mpb.New()

	var wg sync.WaitGroup
	work := make(chan PhotoFromAlbum, 100)
	for i := 0; i < *parallelism; i++ {
		wg.Add(1)
		go processPhotoFromAlbum(ctx, &wg, &bars, *dest, work)
	}

	subctx, cancel := context.WithCancel(ctx)
	defer cancel()
	limiter := syncutil.NewLimiter(subctx, time.Second, 3)

	maxWidth := maxAlbumTitleWidth(albums)
	for _, album := range albums {
		// TODO: could put photos directly to a channel. Will work for large
		// albums.
		photos, err := getPhotos(ctx, limiter, access, *ownerID, strconv.Itoa(album.ID))
		if err != nil {
			log.Fatal(err)
		}
		if len(photos) == 0 {
			continue
		}
		// Prepare directory for this album.
		if err := os.MkdirAll(getAlbumDir(*dest, album), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		bar := progress.AddBar(int64(len(photos)),
			// Prepending decorators
			mpb.PrependDecorators(
				// StaticName decorator with minWidth and no extra config
				// If you need to change name while rendering, use DynamicName
				decor.StaticName(album.Title, maxWidth, decor.DidentRight),
			),
			// Appending decorators
			mpb.AppendDecorators(
				// Percentage decorator with minWidth and no extra config
				decor.Counters("%s/%s", 0, 0, 0),
			),
		)

		bars.Store(album.ID, bar)

		for _, photo := range photos {
			work <- PhotoFromAlbum{photo, album}
		}
	}

	close(work)
	wg.Wait()
	progress.Stop()
}

func maxAlbumTitleWidth(albums []vk.PhotoAlbum) int {
	var max int
	for _, album := range albums {
		if n := len(album.Title); n > max {
			max = n
		}
	}
	return max
}

type PhotoFromAlbum struct {
	Photo vk.Photo
	Album vk.PhotoAlbum
}

func processPhotoFromAlbum(ctx context.Context, wg *sync.WaitGroup, bars *sync.Map, dest string, work <-chan PhotoFromAlbum) {
	defer wg.Done()
	for pa := range work {
		if err := download(ctx, dest, pa.Photo, pa.Album); err != nil {
			log.Fatal(err)
		}

		bar, _ := bars.Load(pa.Album.ID)
		bar.(*mpb.Bar).Increment()
	}
}

func getAlbumDir(root string, album vk.PhotoAlbum) string {
	albumID := album.Title
	if albumID == "" {
		albumID = strconv.Itoa(album.ID)
	}
	return filepath.Clean(fmt.Sprintf("%s/%s", root, albumID))
}

func download(ctx context.Context, dir string, photo vk.Photo, album vk.PhotoAlbum) (err error) {
	largest := getLargestSize(photo.Sizes)

	photoID := strconv.Itoa(photo.ID)
	ext := path.Ext(largest.Src)

	filepath := filepath.Clean(fmt.Sprintf("%s/%s%s", getAlbumDir(dir, album), photoID, ext))

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("GET", largest.Src, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := httputil.CheckResponseStatus(resp); err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)

	return err
}

func getAlbums(ctx context.Context, access *vk.AccessToken, ownerID string) (as []vk.PhotoAlbum, err error) {
	bts, err := vk.Request(ctx, "photos.getAlbums",
		vk.WithAccessToken(access),
		vk.WithParam("owner_id", ownerID),
	)
	if err != nil {
		return nil, err
	}
	var response vk.Response
	if err := response.UnmarshalJSON(bts); err != nil {
		return nil, err
	}
	var albums vk.PhotoAlbums
	if err := albums.UnmarshalJSON(response.Body); err != nil {
		return nil, err
	}
	return albums.Items, nil
}

func getPhotos(ctx context.Context, limiter *syncutil.Limiter, access *vk.AccessToken, ownerID, albumID string) (ps []vk.Photo, err error) {
	const (
		maxCount        = 1000
		maxCountStr     = "1000"
		defaultCoolDown = 50 * time.Millisecond
	)
	var (
		cooldown = defaultCoolDown
		offset   = 0
	)
	for {
		var (
			bts []byte
			err error
		)
		limiter.Do(func() {
			bts, err = vk.Request(ctx, "photos.get",
				vk.WithAccessToken(access),
				vk.WithParam("owner_id", ownerID),
				vk.WithParam("album_id", albumID),
				vk.WithParam("offset", strconv.Itoa(offset)),
				vk.WithParam("count", maxCountStr),
				vk.WithParam("rev", "0"),         // Chronological order.
				vk.WithParam("photo_sizes", "1"), // Special sizes format.
			)
		})
		if err != nil {
			return ps, err
		}
		var response vk.Response
		if err := response.UnmarshalJSON(bts); err != nil {
			return ps, err
		}
		if err := response.Error(); err != nil {
			if err.Temporary() {
				time.Sleep(cooldown)
				cooldown *= 2
				continue
			}
			return ps, err
		}
		cooldown = defaultCoolDown

		var photos vk.Photos
		if err := photos.UnmarshalJSON(response.Body); err != nil {
			return ps, err
		}
		ps = append(ps, photos.Items...)
		if photos.Count < maxCount {
			break
		}
		offset += photos.Count
	}

	return ps, nil
}

func getLargestSize(sizes []vk.PhotoSize) (max vk.PhotoSize) {
	// Range from the end of sizes cause there is a nice chance that 'w' type
	// is the last one.
	for i := len(sizes) - 1; i >= 0; i-- {
		size := sizes[i]
		if size.Type == vk.SizeW {
			// Largest possible size.
			return size
		}
		if max.Type.Less(size.Type) {
			max = size
		}
	}
	return max
}
