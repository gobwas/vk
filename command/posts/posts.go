package posts

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	ClientID     string
	ClientSecret string
	Force        bool
	PreviewSize  int
	ForceLimit   int
	ForcePreview bool
	OnlyReposts  bool
	OnlyOthers   bool
	OnlyOwner    bool
	Store        bool
	StoreDir     string
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
	flag.BoolVar(&c.ForcePreview,
		"force_preview", false,
		"show post info even in force mode",
	)
	flag.IntVar(&c.PreviewSize,
		"preview", 140,
		"preview the post text up to this length",
	)
	flag.IntVar(&c.ForceLimit,
		"limit", -1,
		"force limit (-1 for no limit)",
	)
	flag.BoolVar(&c.OnlyReposts,
		"only_reposts", false,
		"remove only reposts",
	)
	flag.BoolVar(&c.OnlyOthers,
		"only_others", false,
		"remove only others posts on the wall",
	)
	flag.BoolVar(&c.OnlyOwner,
		"only_owner", false,
		"remove only owner's posts on the wall",
	)
	flag.BoolVar(&c.Store,
		"store", false,
		"store posts json backup",
	)
	flag.StringVar(&c.StoreDir,
		"store_dir", download.GetDefaultDest("posts"),
		"store posts json backup dir",
	)
}

type Command struct {
	ui     cli.Ui
	flag   *flag.FlagSet
	config *Config
	limit  *rate.Limiter
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
		limit: rate.NewLimiter(
			rate.Every(vk.DefaultRateInterval),
			vk.DefaultRateBurst,
		),
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

	var filter string
	switch {
	case c.config.OnlyOthers:
		filter = "others"
	case c.config.OnlyOwner || c.config.OnlyReposts:
		filter = "owner"
	default:
		filter = "all"
	}

	var (
		backup *os.File
		bbuf   *bufio.Writer
	)
	if c.config.Store {
		destDir := c.config.StoreDir
		err := os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		backup, err = os.Create(filepath.Clean(destDir + "/" + filter + ".backup." + strconv.FormatInt(time.Now().Unix(), 16) + ".json"))
		if err != nil {
			panic(err)
		}
		defer backup.Close()

		bbuf = bufio.NewWriter(backup)
		defer bbuf.Flush()
	}

	var list vk.Posts
	it := vk.Iterator{
		Method: "wall.get",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("owner_id", access.UserID),
			vk.WithNumber("count", 100),
			vk.WithParam("filter", filter),
		),
		Parse: func(p []byte) (int, error) {
			if bbuf != nil {
				bbuf.Write(p)
			}
			list = vk.Posts{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}

	for it.Next(ctx) {
		for _, post := range list.Items {
			if c.config.OnlyReposts && len(post.CopyHistory) == 0 {
				continue
			}
			if !c.config.Force {
				action, err := vkcli.AskRune(ctx, fmt.Sprintf(
					"delete post dated %s: %s (%s)? ",
					time.Unix(int64(post.Date), 0).Format(time.RFC3339),
					c.postPreview(ctx, access, post),
					homePage(access, post),
				))
				if err != nil {
					log.Fatal(err)
				}
				if action != 'y' {
					continue
				}
			}

		retry:
			err = c.deletePost(ctx, access, strconv.Itoa(post.ID))
			if vk.TemporaryError(err) {
				goto retry
			}
			if err != nil {
				log.Fatal(err)
			}
			if c.config.Force {
				if c.config.ForcePreview {
					fmt.Printf(
						"removed post: %s: %s\n",
						time.Unix(int64(post.Date), 0).Format(time.RFC3339),
						c.postPreview(ctx, access, post),
					)
				} else {
					fmt.Printf(
						"removed post: %s\n",
						time.Unix(int64(post.Date), 0).Format(time.RFC3339),
					)
				}
				c.config.ForceLimit--
				if c.config.ForceLimit == 0 {
					// Turn off force.
					c.config.Force = false
				}
			}
		}
	}
	if err := it.Err(); err != nil {
		log.Fatal(err)
	}

	return 0
}

func (c *Command) postPreview(ctx context.Context, access *vk.AccessToken, post vk.Post) (text string) {
	if n := len(post.CopyHistory); n > 0 {
		post = post.CopyHistory[0]
	}
	if m, n := len(post.Text), c.config.PreviewSize; m > 0 {
		if m > n {
			text = post.Text[:n]
		} else {
			text = post.Text
		}
	}
	if err := c.limit.Wait(ctx); err != nil {
		return text
	}
	if post.FromID < 0 {
		owner, _ := getGroup(ctx, access, post.FromID)
		text = fmt.Sprintf(
			"from group %q: %q",
			owner.Name, text,
		)
	} else {
		user, _ := getUser(ctx, access, post.FromID)
		text = fmt.Sprintf(
			"from user %s %s (%s): %q",
			user.FirstName, user.LastName, user.Domain,
			text,
		)
	}

	return text
}

func getGroup(ctx context.Context, access *vk.AccessToken, groupID int) (group vk.Group, err error) {
	bts, err := vk.Request(ctx, "groups.getById",
		vk.WithAccessToken(access),
		vk.WithNumber("group_id", -1*groupID),
	)
	if err == nil {
		bts, err = vk.StripResponse(bts)
	}
	if err != nil {
		return group, err
	}

	// Need to hack up response.
	bts = append([]byte(`{"items":`), bts...)
	bts = append(bts, '}')

	var list vk.Groups
	err = list.UnmarshalJSON(bts)
	if len(list.Items) > 0 {
		group = list.Items[0]
	}
	return group, err
}

func getUser(ctx context.Context, access *vk.AccessToken, userID int) (user vk.User, err error) {
	bts, err := vk.Request(ctx, "users.get",
		vk.WithAccessToken(access),
		vk.WithNumber("user_ids", userID),
		vk.WithStrings("fields", "domain"),
	)
	if err == nil {
		bts, err = vk.StripResponse(bts)
	}
	if err != nil {
		return user, err
	}

	// Need to hack up response.
	bts = append([]byte(`{"items":`), bts...)
	bts = append(bts, '}')

	var list vk.Users
	err = list.UnmarshalJSON(bts)
	if len(list.Items) > 0 {
		user = list.Items[0]
	}
	return user, err
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

func (c *Command) deletePost(ctx context.Context, access *vk.AccessToken, postID string) error {
	if err := c.limit.Wait(ctx); err != nil {
		return err
	}
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
