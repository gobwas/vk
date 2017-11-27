package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gobwas/vk"
	"github.com/gobwas/vk/cli"
)

var (
	clientID = flag.String(
		"client_id", "",
		"application id",
	)
	clientSecret = flag.String(
		"client_secret", "",
		"application secret",
	)
	ownerID = flag.String(
		"owner_id", "",
		"albums owner id (empty for your id)",
	)
	parallelism = flag.Int(
		"parallelism", 32,
		"number of parallel downloadings",
	)
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

	app := vk.App{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Scope:        vk.ScopePhotos,
	}

	access, err := cli.Authorize(ctx, app)
	if err != nil {
		log.Fatal(err)
	}

	if *ownerID == "" {
		*ownerID = strconv.Itoa(access.UserID)
	}

	posts, err := getPosts(ctx, access, *ownerID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(posts)
}

func getPosts(ctx context.Context, access *vk.AccessToken, ownerID string) ([]vk.Post, error) {
	bts, err := vk.Request(ctx, "wall.get",
		vk.WithAccessToken(access),
		vk.WithParam("owner_id", ownerID),
	)
	if err != nil {
		return nil, err
	}
	log.Println(string(bts))
	var response vk.Response
	if err := response.UnmarshalJSON(bts); err != nil {
		return nil, err
	}
	if err := response.Error(); err != nil {
		return nil, err
	}
	var posts vk.Posts
	if err := posts.UnmarshalJSON(response.Body); err != nil {
		return nil, err
	}
	return posts.Items, nil
}
