package fave

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/gobwas/vk"
	vkcli "github.com/gobwas/vk/cli"
	"github.com/gobwas/vk/internal/download"
	"github.com/mitchellh/cli"
)

func CLI(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return New(ui), nil
	}
}

type Config struct {
	ClientID       string
	ClientSecret   string
	DeleteInterval time.Duration
	Token          string
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
	flag.DurationVar(&c.DeleteInterval,
		"interval", 3*time.Second,
		"interval between deletions",
	)
	flag.StringVar(&c.Token,
		"token", "",
		"token retreived before",
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

	access, err := c.Authorize(ctx)
	if err != nil {
		c.errorf("authorize error: %v", err)
		return 1
	}

	lim := vk.DefaultLimiter()
	limDelete := rate.NewLimiter(
		rate.Every(c.config.DeleteInterval),
		1,
	)

	posts := make(chan vk.Post, 10)
	photos := make(chan vk.Photo, 10)
	videos := make(chan vk.Video, 10)
	go func() {
		if err := getPosts(ctx, access, lim, posts); err != nil {
			log.Fatal(err)
		}
		close(posts)

		if err := getPhotos(ctx, access, lim, photos); err != nil {
			log.Fatal(err)
		}
		close(photos)

		if err := getVideos(ctx, access, lim, videos); err != nil {
			log.Fatal(err)
		}
		close(videos)
	}()

	var n int64
	for post := range posts {
		if err := deleteLike(ctx, access, limDelete, post); err != nil {
			log.Fatal(err)
		}
		log.Println("unliked post", homePage(post))
		n++
	}
	log.Printf("successfully unliked %d posts", n)

	n = 0
	for photo := range photos {
		if err := deleteLike(ctx, access, limDelete, photo); err != nil {
			log.Fatal(err)
		}
		log.Println("unliked photo", download.GetLargestSize(photo.Sizes).Src)
		n++
	}
	log.Printf("successfully unliked %d photos", n)

	n = 0
	for video := range videos {
		if err := deleteLike(ctx, access, limDelete, video); err != nil {
			log.Fatal(err)
		}
		log.Printf("unliked video %s %s", video.Player, video.Photo800)
		n++
	}
	log.Printf("successfully unliked %d videos", n)

	return 0
}

func (c *Command) Authorize(ctx context.Context) (*vk.AccessToken, error) {
	if u := c.config.Token; u != "" {
		return vk.TokenFromURL(u)
	}
	app := vk.App{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        vk.ScopeWall | vk.ScopeFriends | vk.ScopeOffline,
	}
	return vkcli.AuthorizeStandalone(ctx, app)
}

func homePage(post vk.Post) string {
	return "https://vk.com/wall" + strconv.Itoa(post.OwnerID) + "_" + strconv.Itoa(post.ID)
}

func deleteLike(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter, v interface{}) error {
	var (
		t       string
		id      int
		ownerID int
	)
	switch x := v.(type) {
	case vk.Post:
		t = "post"
		id = x.ID
		ownerID = x.OwnerID
	case vk.Photo:
		t = "photo"
		id = x.ID
		ownerID = x.OwnerID
	case vk.Video:
		t = "video"
		id = x.ID
		ownerID = x.OwnerID
	default:
		panic("unkown object")
	}

	c := vk.Caller{
		Method: "likes.delete",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithParam("type", t),
			vk.WithNumber("item_id", id),
			vk.WithNumber("owner_id", ownerID),
		),
		Limiter:        lim,
		ResolveCaptcha: vkcli.ResolveCaptcha,
	}

	_, err := c.Call(ctx)
	return err
}

func getVideos(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter, videos chan<- vk.Video) error {
	var list vk.Videos
	it := vk.Iterator{
		Method:  "fave.getVideos",
		Limiter: lim,
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("count", 50),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Videos{}
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		for _, video := range list.Items {
			videos <- video
		}
	}
	return it.Err()
}

func getPhotos(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter, photos chan<- vk.Photo) error {
	var list vk.Photos
	it := vk.Iterator{
		Method:  "fave.getPhotos",
		Limiter: lim,
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("count", 50),
			vk.WithNumber("photo_sizes", 1),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Photos{}
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		for _, photo := range list.Items {
			photos <- photo
		}
	}
	return it.Err()
}

func getPosts(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter, posts chan<- vk.Post) error {
	var list vk.Posts
	it := vk.Iterator{
		Method:  "fave.getPosts",
		Limiter: lim,
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("count", 100),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Posts{}
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		for _, post := range list.Items {
			posts <- post
		}
	}
	return it.Err()
}

func (c *Command) errorf(f string, args ...interface{}) {
	c.ui.Error(fmt.Sprintf(f, args...))
}

func (c *Command) flagDefaults() string {
	var buf bytes.Buffer
	c.flag.SetOutput(&buf)
	c.flag.PrintDefaults()
	c.flag.SetOutput(os.Stderr)
	return buf.String()
}

func (c *Command) Synopsis() string {
	return "fave command"
}

func (c *Command) Help() string {
	return strings.Join([]string{
		"Usage: fave [options]",
		c.flagDefaults(),
	}, "\n")
}
