package photo

import (
	"bytes"
	"container/ring"
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
	"strings"
	"sync"
	"time"

	"github.com/gobwas/vk"
	vkcli "github.com/gobwas/vk/cli"
	"github.com/gobwas/vk/internal/httputil"
	"github.com/gobwas/vk/internal/syncutil"
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
	OwnerID      string
	Dest         string
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
	flag.StringVar(&c.OwnerID,
		"owner_id", "",
		"albums owner id (empty for your id)",
	)
	flag.StringVar(&c.Dest,
		"dest", getDefaultDest("vkphoto"),
		"destination root dir for photos",
	)
	flag.IntVar(&c.Parallelism,
		"parallelism", 32,
		"number of parallel downloadings",
	)
}

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
		Scope:        vk.ScopePhotos,
	}
	access, err := vkcli.Authorize(ctx, app)
	if err != nil {
		c.errorf("authorize error: %v", err)
		return 1
	}

	ownerID := c.config.OwnerID
	if ownerID == "" {
		ownerID = strconv.Itoa(access.UserID)
	}

	dest := c.config.Dest
	dest += "/" + ownerID

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

	ringLogger := newRingLogger(24)
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

	limiter := syncutil.NewLimiter(time.Second, 3)
	defer limiter.Close()

	photoGetter := &PhotoGetter{
		Access:  access,
		Limiter: limiter,
	}

	maxWidth := maxAlbumTitleWidth(albums)
	for _, album := range albums {
		var photos []vk.Photo
		if album.ID == -4 {
			// Tags album.
			photos, err = photoGetter.GetUserPhotos(ctx, ownerID)
		} else {
			var albumID string
			switch album.ID {
			case -1:
				albumID = "wall"
			case -2:
				albumID = "saved"
			case -3:
				albumID = "profile"
			default:
				albumID = strconv.Itoa(album.ID)
			}
			// TODO: could put photos directly to a channel. Will work for large
			// albums.
			photos, err = photoGetter.GetAlbumPhotos(ctx, ownerID, albumID)
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
		if err := os.MkdirAll(getAlbumDir(dest, album), os.ModePerm); err != nil {
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
		largest := getLargestSize(pa.Photo.Sizes)
		if err := download(ctx, dest, pa.Photo, largest, pa.Album); err != nil {
			log.Printf(
				"download %s (from %q album) error: %v",
				largest.Src, pa.Album.Title, err,
			)
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

func download(ctx context.Context, dir string, photo vk.Photo, size vk.PhotoSize, album vk.PhotoAlbum) (err error) {
	photoID := strconv.Itoa(photo.ID)
	ext := path.Ext(size.Src)

	filepath := filepath.Clean(fmt.Sprintf("%s/%s%s", getAlbumDir(dir, album), photoID, ext))

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("GET", size.Src, nil)
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

type ringLogger struct {
	mu   sync.Mutex
	ring *ring.Ring
}

func newRingLogger(n int) *ringLogger {
	return &ringLogger{
		ring: ring.New(n),
	}
}

func (r *ringLogger) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.ring.Next()
	next.Value = string(p)
	r.ring = r.ring.Move(1)
	return len(p), nil
}

func (r *ringLogger) Interceptor() func(io.Writer) {
	return func(out io.Writer) {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.ring.Do(func(v interface{}) {
			if v == nil {
				return
			}
			out.Write([]byte(v.(string)))
		})
	}
}

type PhotoGetter struct {
	Limiter *syncutil.Limiter
	Access  *vk.AccessToken
}

func (pg *PhotoGetter) GetUserPhotos(ctx context.Context, userID string) ([]vk.Photo, error) {
	return pg.get(ctx, "photos.getUserPhotos",
		vk.WithParam("user_id", userID),
		vk.WithParam("sort", "1"), // Chronological order.
	)
}

func (pg *PhotoGetter) GetAlbumPhotos(ctx context.Context, ownerID, albumID string) ([]vk.Photo, error) {
	return pg.get(ctx, "photos.get",
		vk.WithParam("owner_id", ownerID),
		vk.WithParam("album_id", albumID),
		vk.WithParam("rev", "0"),         // Chronological order.
		vk.WithParam("photo_sizes", "1"), // Special sizes format.
	)
}

func (pg *PhotoGetter) get(ctx context.Context, method string, queryOptions ...vk.QueryOption) (ret []vk.Photo, err error) {
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
		pg.Limiter.Do(func() {
			bts, err = vk.Request(ctx, method,
				vk.WithOptions(queryOptions),
				vk.WithAccessToken(pg.Access),
				vk.WithNumber("offset", offset),
				vk.WithParam("count", maxCountStr),
				vk.WithParam("photo_sizes", "1"), // Special sizes format.
			)
		})
		if err != nil {
			return ret, err
		}
		bts, err = vk.StripResponse(bts)
		if err != nil {
			if vkErr, ok := err.(*vk.Error); ok && vkErr.Temporary() {
				time.Sleep(cooldown)
				cooldown *= 2
				continue
			}
			return ret, err
		}
		cooldown = defaultCoolDown

		var photos vk.Photos
		if err := photos.UnmarshalJSON(bts); err != nil {
			return ret, err
		}

		ret = append(ret, photos.Items...)
		if photos.Count < maxCount {
			// No need to repeat request.
			break
		}
		offset += photos.Count
	}
	return ret, nil
}
