package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

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

	var list vk.Posts
	it := vk.Iterator{
		Method: "wall.get",
		Options: vk.QueryOptions(
			vk.WithAccessToken(access),
			vk.WithNumber("owner_id", access.UserID),
			vk.WithNumber("count", 100),
			vk.WithParam("filter", "owner"),
		),
		Parse: func(p []byte) (int, error) {
			list = vk.Posts{} // Reset.
			err := list.UnmarshalJSON(p)
			return len(list.Items), err
		},
	}
	for it.Next(ctx) {
		for _, post := range list.Items {
			if n := len(post.CopyHistory); n > 0 {
				if !*force {
					action, err := cli.Prompt(ctx, fmt.Sprintf(
						"delete post from %s %s ? ",
						time.Unix(int64(post.Date), 0).Format(time.RFC3339),
						homePage(access, post),
					))
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
	if err := it.Err(); err != nil {
		log.Fatal(err)
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
