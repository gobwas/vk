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
	force = flag.Bool(
		"force", false,
		"delete friends withoud prompt",
	)
)

func main() {
	flag.Parse()

	ctx := context.Background()

	app := vk.App{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Scope:        vk.ScopeFriends,
	}
	access, err := cli.AuthorizeStandalone(ctx, app)
	if err != nil {
		log.Fatal(err)
	}

	friends, err := getFriends(ctx, access)
	if err != nil {
		log.Fatal(err)
	}

	for _, friend := range friends {
		if !*force {
			action, err := cli.Prompt(ctx,
				fmt.Sprintf(
					"delete %s %s (%s)? ",
					friend.FirstName, friend.LastName, homePage(friend),
				),
			)
			if err != nil {
				log.Fatal(err)
			}
			if action != "y" {
				continue
			}
		}
		err := deleteFriend(ctx, access, friend)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func homePage(f vk.Friend) string {
	s := f.Domain
	if s == "" {
		s = strconv.Itoa(f.ID)
	}
	return "https://vk.com/" + s
}

func deleteFriend(ctx context.Context, access *vk.AccessToken, friend vk.Friend) error {
	bts, err := vk.Request(ctx, "friends.delete",
		vk.WithAccessToken(access),
		vk.WithParam("user_id", strconv.Itoa(friend.ID)),
	)
	if err == nil {
		_, err = vk.StripResponse(bts)
	}
	return err
}

func getFriends(ctx context.Context, access *vk.AccessToken) ([]vk.Friend, error) {
	bts, err := vk.Request(ctx, "friends.get",
		vk.WithAccessToken(access),
		vk.WithParam("user_id", strconv.Itoa(access.UserID)),
		vk.WithParam("fields", "domain"),
	)
	if err != nil {
		return nil, err
	}
	if bts, err = vk.StripResponse(bts); err != nil {
		return nil, err
	}
	var fs vk.Friends
	if err := fs.UnmarshalJSON(bts); err != nil {
		return nil, err
	}
	return fs.Items, nil
}
