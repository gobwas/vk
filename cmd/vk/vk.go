package main

import (
	"log"
	"os"

	"github.com/gobwas/vk/command/posts"
	"github.com/gobwas/vk/command/stub"
	"github.com/mitchellh/cli"
)

var (
	name    = "vk"
	version = "0.0.0"
)

func main() {
	ui := cli.BasicUi{
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c := cli.NewCLI(name, version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"stub":  stub.CLI(&ui),
		"posts": posts.CLI(&ui),
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
