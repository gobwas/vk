package photos

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gobwas/vk"
	vkcli "github.com/gobwas/vk/cli"
	"github.com/gobwas/vk/internal/download"
	"github.com/gobwas/vk/internal/logutil"
	"github.com/mitchellh/cli"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func CLI(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return New(ui), nil
	}
}

type Config struct {
	ClientID     string
	ClientSecret string
	Dest         string
	OwnerID      int
	Parallelism  int
}

func (c *Config) ExportTo(flag *flag.FlagSet) {
	flag.StringVar(&c.ClientID,
		"client_id", "",
		"application id",
	)
	flag.StringVar(&c.ClientSecret,
		"client_secret", "",
		"application secret",
	)
	flag.IntVar(&c.OwnerID,
		"owner_id", 0,
		"albums owner id (empty for your id)",
	)
	flag.StringVar(&c.Dest,
		"dest", download.GetDefaultDest("photos"),
		"destination root dir for photos",
	)
	flag.IntVar(&c.Parallelism,
		"parallelism", 64,
		"number of parallel downloads",
	)
}

type Command struct {
	ui     cli.Ui
	flag   *flag.FlagSet
	config *Config
}

func New(ui cli.Ui) *Command {
	flag := flag.NewFlagSet("", flag.ContinueOnError)
	flag.Usage = func() {}

	c := new(Config)
	c.ExportTo(flag)

	return &Command{
		ui:     ui,
		flag:   flag,
		config: c,
	}
}

func (c *Command) Run(args []string) int {
	if err := c.flag.Parse(args); err != nil {
		return cli.RunResultHelp
	}

	ctx := context.Background()

	app := vk.App{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        vk.ScopePhotos | vk.ScopeWall,
	}
	access, err := vkcli.AuthorizeStandalone(ctx, app)
	if err != nil {
		c.errorf("authorize error: %v", err)
		return 1
	}

	ownerID := c.config.OwnerID
	if ownerID == 0 {
		ownerID = access.UserID
	}

	dest := c.config.Dest
	dest += "/" + strconv.Itoa(ownerID)

	albums, err := getAlbums(ctx, access, ownerID)
	if err != nil {
		log.Fatal(err)
	}
	albums = append(albums,
		wallAlbum,
		savedAlbum,
		profileAlbum,
		tagsAlbum,
	)

	c.logf("ready to store photos at %s", dest)

	ringLogger := logutil.NewRingLogger(24)
	log.SetOutput(ringLogger)
	log.SetFlags(0)

	bars := sync.Map{}
	progress := mpb.New(
		mpb.Output(os.Stderr),
		mpb.OutputInterceptors(ringLogger.Interceptor()),
	)

	var wg sync.WaitGroup
	work := make(chan PhotoFromAlbum, 100)
	for i := 0; i < c.config.Parallelism; i++ {
		wg.Add(1)
		go processPhotoFromAlbum(ctx, &wg, &bars, dest, work)
	}

	maxWidth := maxAlbumTitleWidth(albums)
	for _, album := range albums {
		var photos []vk.Photo
		if album.ID == -4 {
			photos, err = getUserTaggedPhotos(ctx, access, ownerID)
		} else {
			photos, err = getAlbumPhotos(ctx, access, ownerID, album.ID)
		}
		if err != nil {
			log.Printf(
				"get photos for album %q (%d) error: %v",
				album.Title, album.ID, err,
			)
			continue
		}
		if len(photos) == 0 {
			continue
		}
		// Prepare directory for this album.
		if err := os.MkdirAll(appendAlbumDir(dest, album), os.ModePerm); err != nil {
			panic(err)
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

	return 0
}

func (c *Command) errorf(f string, args ...interface{}) {
	c.ui.Error(fmt.Sprintf(f, args...))
}

func (c *Command) logf(f string, args ...interface{}) {
	c.ui.Info(fmt.Sprintf(f, args...))
}

func (c *Command) flagDefaults() string {
	var buf bytes.Buffer
	c.flag.SetOutput(&buf)
	c.flag.PrintDefaults()
	c.flag.SetOutput(os.Stderr)
	return buf.String()
}

func (c *Command) Synopsis() string {
	return "photo command"
}

func (c *Command) Help() string {
	return strings.Join([]string{
		"Usage: photo [options]",
		c.flagDefaults(),
	}, "\n")
}

var (
	wallAlbum = vk.PhotoAlbum{
		ID:    -1,
		Title: "wall",
	}
	savedAlbum = vk.PhotoAlbum{
		ID:    -2,
		Title: "saved",
	}
	profileAlbum = vk.PhotoAlbum{
		ID:    -3,
		Title: "profile",
	}
	tagsAlbum = vk.PhotoAlbum{
		ID:    -4,
		Title: "tags",
	}
)

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
		largest := download.GetLargestSize(pa.Photo.Sizes)
		if err := download.Photo(ctx, appendAlbumDir(dest, pa.Album), pa.Photo, largest); err != nil {
			log.Printf(
				"download %s (from %q album) error: %v",
				largest.Src, pa.Album.Title, err,
			)
		}
		bar, _ := bars.Load(pa.Album.ID)
		bar.(*mpb.Bar).Increment()
	}
}

func appendAlbumDir(root string, album vk.PhotoAlbum) string {
	albumID := album.Title
	if albumID == "" {
		albumID = strconv.Itoa(album.ID)
	}
	return filepath.Clean(fmt.Sprintf("%s/%s", root, albumID))
}

func getAlbums(ctx context.Context, access *vk.AccessToken, ownerID int) (as []vk.PhotoAlbum, err error) {
	bts, err := vk.Request(ctx, "photos.getAlbums",
		vk.WithAccessToken(access),
		vk.WithNumber("owner_id", ownerID),
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

func getUserTaggedPhotos(ctx context.Context, access *vk.AccessToken, userID int) (ps []vk.Photo, err error) {
	var list vk.Photos
	it := vk.Iterator{
		Method: "photos.getUserPhotos",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("user_id", userID),
			vk.WithNumber("sort", 1),        // Chronological order.
			vk.WithNumber("photo_sizes", 1), // Special sizes format.
			vk.WithNumber("count", 1000),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Photos{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		ps = append(ps, list.Items...)
	}
	return ps, it.Err()
}

func getAlbumPhotos(ctx context.Context, access *vk.AccessToken, ownerID, albumID int) (ps []vk.Photo, err error) {
	var album string
	switch albumID {
	case -1:
		album = "wall"
	case -2:
		album = "saved"
	case -3:
		album = "profile"
	default:
		album = strconv.Itoa(albumID)
	}
	var list vk.Photos
	it := vk.Iterator{
		Method: "photos.get",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("owner_id", ownerID),
			vk.WithParam("album_id", album),
			vk.WithNumber("rev", 0),         // Chronological order.
			vk.WithNumber("photo_sizes", 1), // Special sizes format.
			vk.WithNumber("count", 1000),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Photos{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		ps = append(ps, list.Items...)
	}
	return ps, it.Err()
}
