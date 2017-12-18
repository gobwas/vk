package main

import (
	"log"
	"os"

	"github.com/gobwas/vk/command/fave"
	"github.com/gobwas/vk/command/friends"
	"github.com/gobwas/vk/command/messages"
	"github.com/gobwas/vk/command/photos"
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
		"stub":     stub.CLI(&ui),
		"posts":    posts.CLI(&ui),
		"photos":   photos.CLI(&ui),
		"friends":  friends.CLI(&ui),
		"messages": messages.CLI(&ui),
		"fave":     fave.CLI(&ui),
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
