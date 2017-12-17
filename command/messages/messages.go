package messages

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/briandowns/spinner"
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
	All          bool
	Parallelism  int
	Delete       bool
	Save         bool
	TokenURL     string
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
		"dest", download.GetDefaultDest("messages"),
		"destination root dir for saved chats",
	)
	flag.BoolVar(&c.All,
		"all", false,
		"download all dialogs from current user",
	)
	flag.IntVar(&c.Parallelism,
		"parallelism", 16,
		"number of simultaneous downloads",
	)
	flag.BoolVar(&c.Delete,
		"delete", false,
		"delete all dialogs",
	)
	flag.BoolVar(&c.Save,
		"save", false,
		"save all dialogs",
	)
	flag.StringVar(&c.TokenURL,
		"token", "",
		"token url copied from browser",
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

	var (
		access *vk.AccessToken
		err    error
	)
	if u := c.config.TokenURL; u != "" {
		access, err = vk.TokenFromURL(u)
	} else {
		app := vk.App{
			ClientID:     c.config.ClientID,
			ClientSecret: c.config.ClientSecret,
			Scope:        vk.ScopeMessages,
		}
		access, err = vkcli.AuthorizeStandalone(ctx, app)
	}
	if err != nil {
		c.errorf("authorize error: %v", err)
		return 1
	}

	lim := rate.NewLimiter(rate.Every(time.Second/3), 3)

	ds, err := getDialogs(ctx, access, lim)
	if err != nil {
		log.Fatal(err)
	}

	users, err := getUsers(ctx, access, lim, append([]int{access.UserID}, usersFromDialogs(ds)...))
	if err != nil {
		log.Fatal(err)
	}

	me := users[0]
	users = users[1:]
	if len(users) == 0 {
		return 0
	}

	var (
		progress *mpb.Progress
		bar      *mpb.Bar
		sem      chan struct{}
	)
	if c.config.All {
		sem = make(chan struct{}, c.config.Parallelism)

		ringLogger := logutil.NewRingLogger(c.config.Parallelism)
		log.SetOutput(ringLogger)
		log.SetFlags(0)

		progress = mpb.New(
			mpb.Output(os.Stderr),
			mpb.OutputInterceptors(
				func(w io.Writer) {
					w.Write([]byte{'\n'})
				},
				ringLogger.Interceptor(),
			),
		)
		barTitle := "dialogs"
		bar = progress.AddBar(int64(len(users)),
			// Prepending decorators
			mpb.PrependDecorators(
				// StaticName decorator with minWidth and no extra config
				// If you need to change name while rendering, use DynamicName
				decor.StaticName(barTitle, len(barTitle), decor.DidentRight),
			),
			// Appending decorators
			mpb.AppendDecorators(
				// Percentage decorator with minWidth and no extra config
				decor.Counters("%s/%s", 0, 0, 0),
			),
		)

		fmt.Println()
	}
	for _, user := range users {
		if !c.config.All {
			if c.config.Save {
				action, err := vkcli.AskRune(ctx, fmt.Sprintf(
					"save chat with %s %s (%s)? ",
					user.LastName, user.FirstName, user.Domain,
				))
				if err != nil {
					log.Fatal(err)
				}
				if action == 'y' {
					destDir := appendUserDir(c.config.Dest, user)

					fmt.Printf("\tsaving messages at '%s' ", destDir)
					s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
					s.Start()

					err := c.saveChat(ctx, access, lim, destDir, me, user)
					s.Stop()
					if err != nil {
						fmt.Printf("error: %v", err)
					}
					fmt.Print("\n")
				}
			}
			if c.config.Delete {
				action, err := vkcli.AskRune(ctx, fmt.Sprintf(
					"delete chat with %s %s (%s)? ",
					user.LastName, user.FirstName, user.Domain,
				))
				if err != nil {
					log.Fatal(err)
				}
				if action == 'y' {
					fmt.Printf(
						"\tdeleting chat with %s %s (%s) ",
						user.LastName, user.FirstName, user.Domain,
					)
					s := spinner.New(spinner.CharSets[42], 100*time.Millisecond)
					s.Start()

					err := c.deleteChat(ctx, access, lim, user)
					s.Stop()
					if err != nil {
						fmt.Printf("error: %v", err)
					}
					fmt.Print("\n")
				}
			}
			continue
		}

		sem <- struct{}{}

		destDir := appendUserDir(c.config.Dest, user)
		user := user // For closure.
		go func() {
			defer func() {
				bar.Increment()
				<-sem
			}()
			if c.config.Save {
				if err := c.saveChat(ctx, access, lim, destDir, me, user); err != nil {
					log.Printf(
						"error saving messages from %s %s (%s): %v",
						user.FirstName, user.LastName, user.Domain, err,
					)
				} else {
					log.Printf(
						"messages from %s %s (%s) are stored at '%s'",
						user.FirstName, user.LastName, user.Domain, destDir,
					)
				}
			}
			if c.config.Delete {
				if err := c.deleteChat(ctx, access, lim, user); err != nil {
					log.Printf(
						"delete messages from %s %s (%s) error: %v",
						user.FirstName, user.LastName, user.Domain, err,
					)
				} else {
					log.Printf(
						"deleted messages from %s %s (%s)",
						user.FirstName, user.LastName, user.Domain,
					)
				}
			}
		}()
	}
	// Wait all workers are done.
	for i := 0; i < c.config.Parallelism; i++ {
		sem <- struct{}{}
	}
	if progress != nil {
		progress.Stop()
	}

	return 0
}

func (c *Command) deleteChat(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter, user vk.User) error {
	var list vk.Messages
	it := vk.Iterator{
		Method: "messages.getHistory",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("user_id", user.ID),
			vk.WithNumber("count", 200),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Messages{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
		Limiter: lim,
	}
	ids := make([]int, 0, 200)
	for it.Next(ctx) {
		ids = ids[:0]
		for _, message := range list.Items {
			ids = append(ids, message.ID)
		}
	retry:
		if err := lim.Wait(ctx); err != nil {
			return err
		}
		bts, err := vk.Request(ctx, "messages.delete",
			vk.WithAccessToken(access),
			vk.WithNumbers("message_ids", ids...),
		)
		if err == nil {
			_, err = vk.StripResponse(bts)
		}
		if vk.TemporaryError(err) {
			goto retry
		}
		if err != nil {
			return err
		}
	}
	if err := it.Err(); err != nil {
		return err
	}

retryd:
	if err := lim.Wait(ctx); err != nil {
		return err
	}
	bts, err := vk.Request(ctx, "messages.deleteDialog",
		vk.WithAccessToken(access),
		vk.WithNumber("user_id", user.ID),
		vk.WithNumber("count", 10000),
	)
	if err == nil {
		_, err = vk.StripResponse(bts)
	}
	if vk.TemporaryError(err) {
		goto retryd
	}
	return err
}

func (c *Command) saveChat(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter, userDir string, me, user vk.User) error {
	type TemplateData struct {
		Messages <-chan vk.Message
		User     vk.User
		Me       vk.User
	}

	var (
		once     sync.Once
		raw      *os.File
		rawBuf   *bufio.Writer
		html     *os.File
		htmlBuf  *bufio.Writer
		messages chan vk.Message
	)
	init := func() (ok bool) {
		once.Do(func() {
			ok = true
			err := os.MkdirAll(userDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
			raw, err = os.Create(filepath.Clean(userDir + "/raw.json"))
			if err != nil {
				panic(err)
			}
			rawBuf = bufio.NewWriter(raw)

			html, err = os.Create(filepath.Clean(userDir + "/index.html"))
			if err != nil {
				panic(err)
			}
			htmlBuf = bufio.NewWriter(html)

			messages = make(chan vk.Message, 10)
			go func() {
				t.Execute(htmlBuf, TemplateData{
					Messages: messages,
					User:     user,
					Me:       me,
				})
			}()
		})
		return
	}

	var (
		list vk.Messages
		bts  []byte
	)
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
			bts = p
			list = vk.Messages{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
		Limiter: lim,
	}
	for it.Next(ctx) {
		if init() {
			defer func() {
				close(messages)

				rawBuf.Flush()
				raw.Close()

				htmlBuf.Flush()
				html.Close()
			}()
		}

		rawBuf.Write(bts)
		rawBuf.WriteByte('\n')

		for _, message := range list.Items {
			messages <- message
			for _, attach := range message.Attachments {
				if attach.Type != "photo" {
					continue
				}
				size := download.GetLargestSize(attach.Photo.Sizes)
				err := download.Photo(ctx, userDir, attach.Photo, size)
				if err != nil {
					log.Printf(
						"download %s %s (%s) attachement photo %s error: %v",
						user.FirstName, user.LastName, user.Domain, size.Src, err,
					)
				}
			}
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

func getUsers(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter, userIDs []int) ([]vk.User, error) {
	users := make([]vk.User, 0, len(userIDs))

	for i := 0; i < len(userIDs); i += 1000 {
		sub := userIDs[i:]
		if len(sub) > 1000 {
			sub = sub[:1000]
		}
		if err := lim.Wait(ctx); err != nil {
			return users, err
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

func getDialogs(ctx context.Context, access *vk.AccessToken, lim *rate.Limiter) (ret []vk.Dialog, err error) {
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
		Limiter: lim,
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
