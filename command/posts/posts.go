package posts

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

	"github.com/gobwas/vk"
	vkcli "github.com/gobwas/vk/cli"
	"github.com/mitchellh/cli"
)

func CLI(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return New(ui), nil
	}
}

type Config struct {
	ClientID     string
	ClientSecret string
	Force        bool
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
	flag.BoolVar(&c.Force,
		"force", false,
		"do not ask for deletion",
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
		Scope:        vk.ScopeWall,
	}
	access, err := vkcli.AuthorizeStandalone(ctx, app)
	if err != nil {
		c.errorf("authorize error: %v", err)
		return 1
	}

	var list vk.Posts
	it := vk.Iterator{
		Method: "wall.get",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("owner_id", access.UserID),
			vk.WithNumber("count", 100),
			vk.WithParam("filter", "owner"),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Posts{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		for _, post := range list.Items {
			if n := len(post.CopyHistory); n > 0 {
				if !c.config.Force {
					action, err := vkcli.AskRune(ctx, fmt.Sprintf(
						"delete post from %s %s ? ",
						time.Unix(int64(post.Date), 0).Format(time.RFC3339),
						homePage(access, post),
					))
					if err != nil {
						log.Fatal(err)
					}
					if action != 'y' {
						continue
					}
				}
				if err := deletePost(ctx, access, strconv.Itoa(post.ID)); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	if err := it.Err(); err != nil {
		log.Fatal(err)
	}

	return 0
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
	return "posts command"
}

func (c *Command) Help() string {
	return strings.Join([]string{
		"Usage: posts [options]",
		c.flagDefaults(),
	}, "\n")
}

func homePage(access *vk.AccessToken, post vk.Post) string {
	return "https://vk.com/wall" + strconv.Itoa(access.UserID) + "_" + strconv.Itoa(post.ID)
}

func deletePost(ctx context.Context, access *vk.AccessToken, postID string) error {
	bts, err := vk.Request(ctx, "wall.delete",
		vk.WithAccessToken(access),
		vk.WithParam("owner_id", strconv.Itoa(access.UserID)),
		vk.WithParam("post_id", postID),
	)
	if err == nil {
		_, err = vk.StripResponse(bts)
	}
	return err
}
