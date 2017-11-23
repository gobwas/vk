package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gobwas/vk"
)

var (
	clientID     = flag.String("client_id", "", "application id")
	clientSecret = flag.String("client_secret", "", "application secret")
)

func main() {
	flag.Parse()

	if *clientID == "" || *clientSecret == "" {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options]\n\n",
			os.Args[0],
		)
		flag.CommandLine.SetOutput(os.Stderr)
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()

	auth := vk.Auth{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		State:        "OK",
		Scope:        vk.ScopePhotos,
	}

	access, err := auth.Authorize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	bts, err := vk.Request(ctx, "photos.getAlbums",
		vk.WithAccess(access),
		vk.WithParam("owner_id", strconv.Itoa(access.UserID)),
	)
	if err != nil {
		log.Fatal(err)
	}

	os.Stderr.Write(bts)
}
