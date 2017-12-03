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
		Scope:        vk.ScopeWall,
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

	for _, post := range posts {
		log.Println(len(post.CopyHistory))
		if n := len(post.CopyHistory); n > 0 {
			log.Println("delete", deletePost(ctx, access, *ownerID, strconv.Itoa(post.ID)))
		}
	}
}

func deletePost(ctx context.Context, access *vk.AccessToken, ownerID, postID string) error {
	bts, err := vk.Request(ctx, "wall.delete",
		vk.WithAccessToken(access),
		vk.WithParam("owner_id", ownerID),
		vk.WithParam("post_id", postID),
	)
	if err == nil {
		_, err = vk.StripResponse(bts)
	}
	return err
}

func getPosts(ctx context.Context, access *vk.AccessToken, ownerID string) ([]vk.Post, error) {
	bts, err := vk.Request(ctx, "wall.get",
		vk.WithAccessToken(access),
		vk.WithParam("owner_id", ownerID),
		vk.WithParam("filter", "owner"),
		vk.WithParam("count", "2"),
	)

	log.Println(string(bts))
	if err != nil {
		return nil, err
	}
	bts, err = vk.StripResponse(bts)
	if err != nil {
		return nil, err
	}
	var posts vk.Posts
	if err := posts.UnmarshalJSON(bts); err != nil {
		return nil, err
	}
	return posts.Items, nil
}
