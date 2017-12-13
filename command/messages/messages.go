package messages

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
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
	ClientID     string
	ClientSecret string
	Dest         string
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
	flag.StringVar(&c.Dest,
		"dest", download.GetDefaultDest("vkmessages"),
		"destination root dir for saved chats",
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
		Scope:        vk.ScopeMessages,
	}
	access, err := vkcli.AuthorizeStandalone(ctx, app)
	if err != nil {
		c.errorf("authorize error: %v", err)
		return 1
	}

	ds, err := getDialogs(ctx, access)
	if err != nil {
		log.Fatal(err)
	}

	users, err := getUsers(ctx, access, usersFromDialogs(ds))
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		action, err := vkcli.AskRune(ctx, fmt.Sprintf(
			"save chat for %s %s (%s)? ",
			user.LastName, user.FirstName, user.Domain,
		))
		if err != nil {
			log.Fatal(err)
		}
		if action != 'y' {
			continue
		}
		if err := c.saveChat(ctx, access, user); err != nil {
			log.Fatal(err)
		}
	}

	return 0
}

func (c *Command) saveChat(ctx context.Context, access *vk.AccessToken, user vk.User) error {
	userDir := appendUserDir(c.config.Dest, user)

	var (
		once     sync.Once
		file     *os.File
		bw       *bufio.Writer
		progress *spinner.Spinner
	)
	init := func() (ok bool) {
		once.Do(func() {
			ok = true
			err := os.MkdirAll(userDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			file, err = os.Create(filepath.Clean(userDir + "/history.json"))
			if err != nil {
				panic(err)
			}

			fmt.Printf("\tsaving messages to '%s' ", userDir)
			progress = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			progress.Start()

			bw = bufio.NewWriter(file)
		})
		return
	}

	var list vk.Messages
	it := vk.Iterator{
		Method: "messages.getHistory",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("user_id", user.ID),
			vk.WithNumber("count", 200),
			vk.WithNumber("rev", 1),          // Reverse chronological order.
			vk.WithParam("photo_sizes", "1"), // Special sizes format.
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Messages{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		if init() {
			defer func() {
				bw.Flush()
				file.Close()
				progress.Stop()
				fmt.Printf("\n")
			}()
		}

		for _, message := range list.Items {
			for _, attach := range message.Attachments {
				if attach.Type != "photo" {
					continue
				}
				size := download.GetLargestSize(attach.Photo.Sizes)
				err := download.Photo(ctx, userDir, attach.Photo, size)
				if err != nil {
					log.Printf("download attachement photo %s error: %v", size.Src, err)
				}
			}
			bts, err := message.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			bw.Write(bts)
			bw.WriteByte('\n')
		}
	}
	if err := it.Err(); err != nil {
		return err
	}

	return nil
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
	return "messages command"
}

func (c *Command) Help() string {
	return strings.Join([]string{
		"Usage: messages [options]",
		c.flagDefaults(),
	}, "\n")
}

func getUsers(ctx context.Context, access *vk.AccessToken, userIDs []int) ([]vk.User, error) {
	users := make([]vk.User, 0, len(userIDs))

	for i := 0; i < len(userIDs); i += 1000 {
		sub := userIDs[i:]
		if len(sub) > 1000 {
			sub = sub[:1000]
		}
		bts, err := vk.Request(ctx, "users.get",
			vk.WithAccessToken(access),
			vk.WithNumbers("user_ids", sub...),
			vk.WithStrings("fields", "domain"),
		)
		if err == nil {
			bts, err = vk.StripResponse(bts)
		}
		if err != nil {
			return nil, err
		}

		// Need to hack up response.
		bts = append([]byte(`{"items":`), bts...)
		bts = append(bts, '}')

		var list vk.Users
		if err = list.UnmarshalJSON(bts); err != nil {
			return nil, err
		}
		users = append(users, list.Items...)
	}

	return users, nil
}

func getDialogs(ctx context.Context, access *vk.AccessToken) (ret []vk.Dialog, err error) {
	var list vk.Dialogs
	it := vk.Iterator{
		Method: "messages.getDialogs",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("count", 200),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Dialogs{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		for _, d := range list.Items {
			ret = append(ret, d)
		}
	}
	return ret, it.Err()
}

func appendUserDir(dest string, user vk.User) string {
	return filepath.Clean(fmt.Sprintf(
		"%s/%s %s (%s)",
		dest, user.LastName, user.FirstName, user.Domain,
	))
}

func usersFromDialogs(ds []vk.Dialog) []int {
	ret := make([]int, len(ds))
	for i, d := range ds {
		ret[i] = d.Message.UserID
	}
	return ret
}
