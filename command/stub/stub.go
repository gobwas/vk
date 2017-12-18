package stub

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
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
	Token        string
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

	panic(access.Token)

	return 0
}

func (c *Command) Authorize(ctx context.Context) (*vk.AccessToken, error) {
	if u := c.config.Token; u != "" {
		return vk.TokenFromURL(u)
	}
	app := vk.App{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        vk.ScopeMessages,
	}
	return vkcli.AuthorizeStandalone(ctx, app)
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
	return "stub command"
}

func (c *Command) Help() string {
	return strings.Join([]string{
		"Usage: stub [options]",
		c.flagDefaults(),
	}, "\n")
}
