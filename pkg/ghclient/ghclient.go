package ghclient

import (
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

func NewGhClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}
