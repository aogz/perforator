package gh

import (
	"context"

	"github.com/google/go-github/v39/github"
)

// GetTickets returns a list of tickets filtered by label for specified repo
func GetPRs(client *github.Client, owner string, repo string, limit int) ([]*github.PullRequest, error) {
	options := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	prs, _, err := client.PullRequests.List(context.Background(), owner, repo, options)
	return prs, err
}

func GetPullRequestReviews(client *github.Client, owner string, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	reviews, _, err := client.PullRequests.ListReviews(context.Background(), owner, repo, prNumber, nil)
	return reviews, err
}
