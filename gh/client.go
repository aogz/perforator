package gh

import (
	"context"
	"os"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

// GetClient ...
func GetClient() *github.Client {
	token := &oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")}
	tokenSource := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(context.Background(), tokenSource)
	return github.NewClient(tc)
}
