package friends

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
		"delete friends without prompt",
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

	friends, err := getFriends(ctx, access)
	if err != nil {
		log.Fatal(err)
	}

	for _, friend := range friends {
		if !c.config.Force {
			action, err := vkcli.AskRune(ctx,
				fmt.Sprintf(
					"delete %s %s (%s)? ",
					friend.FirstName, friend.LastName, homePage(friend),
				),
			)
			if err != nil {
				log.Fatal(err)
			}
			if action != 'y' {
				continue
			}
		}
		err := deleteFriend(ctx, access, friend)
		if err != nil {
			log.Fatal(err)
		}
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
	return "friends command"
}

func (c *Command) Help() string {
	return strings.Join([]string{
		"Usage: friends [options]",
		c.flagDefaults(),
	}, "\n")
}

func homePage(f vk.Friend) string {
	s := f.Domain
	if s == "" {
		s = strconv.Itoa(f.ID)
	}
	return "https://vk.com/" + s
}

func deleteFriend(ctx context.Context, access *vk.AccessToken, friend vk.Friend) error {
	bts, err := vk.Request(ctx, "friends.delete",
		vk.WithAccessToken(access),
		vk.WithParam("user_id", strconv.Itoa(friend.ID)),
	)
	if err == nil {
		_, err = vk.StripResponse(bts)
	}
	return err
}

func getFriends(ctx context.Context, access *vk.AccessToken) ([]vk.Friend, error) {
	bts, err := vk.Request(ctx, "friends.get",
		vk.WithAccessToken(access),
		vk.WithParam("user_id", strconv.Itoa(access.UserID)),
		vk.WithParam("fields", "domain"),
	)
	if err != nil {
		return nil, err
	}
	if bts, err = vk.StripResponse(bts); err != nil {
		return nil, err
	}
	var fs vk.Friends
	if err := fs.UnmarshalJSON(bts); err != nil {
		return nil, err
	}
	return fs.Items, nil
}
