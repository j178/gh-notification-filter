package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is required")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opts := &github.NotificationListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		log.Println("Fetching notifications...")
		nt, resp, err := client.Activity.ListNotifications(ctx, opts)
		if err != nil {
			log.Printf("Fetching error: %v", err)
		}
		log.Printf("Fetched %d notifications\n", len(nt))
		if len(nt) == 0 {
			break
		}
		for _, n := range nt {
			if strings.Contains(strings.ToLower(n.Subject.GetTitle()), "typo") {
				log.Printf("Marking notification `%s` as read, url: %s", n.Subject.GetTitle(), n.GetURL())
				_, err := client.Activity.MarkThreadRead(ctx, n.GetID())
				if err != nil {
					log.Printf("Makring error: %v", err)
				}
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
}
