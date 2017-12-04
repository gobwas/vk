package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	force = flag.Bool(
		"force", false,
		"delete posts without prompt",
	)
)

func main() {
	flag.Parse()

	ctx := context.Background()

	app := vk.App{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Scope:        vk.ScopeWall,
	}
	access, err := cli.AuthorizeStandalone(ctx, app)
	if err != nil {
		log.Fatal(err)
	}

	posts, err := getPosts(ctx, access)
	if err != nil {
		log.Fatal(err)
	}
	for _, post := range posts {
		if n := len(post.CopyHistory); n > 0 {
			if !*force {
				action, err := cli.Prompt(ctx, fmt.Sprintf("delete %s ? ", homePage(access, post)))
				if err != nil {
					log.Fatal(err)
				}
				if action != "y" {
					continue
				}
			}
			if err := deletePost(ctx, access, strconv.Itoa(post.ID)); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func homePage(access *vk.AccessToken, post vk.Post) string {
	return "https://vk.com/wall" + strconv.Itoa(access.UserID) + "_" + strconv.Itoa(post.ID)
}

func deletePost(ctx context.Context, access *vk.AccessToken, postID string) error {
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

func getPosts(ctx context.Context, access *vk.AccessToken) ([]vk.Post, error) {
	bts, err := vk.Request(ctx, "wall.get",
		vk.WithAccessToken(access),
		vk.WithParam("owner_id", strconv.Itoa(access.UserID)),
		vk.WithParam("filter", "owner"),
	)
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
