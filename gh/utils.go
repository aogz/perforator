package gh

import (
	"context"
	"fmt"
	"math"

	"github.com/aogz/perforator/utils"
	"github.com/google/go-github/v40/github"
)

const MAX_PER_PAGE = 100

// GetPRs returns a list of tickets filtered by label for specified repo
func GetPRs(client *github.Client, owner string, repo string, limit int) ([]*github.PullRequest, error) {
	options := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: limit,
		},
	}

	var pullRequestList []*github.PullRequest
	if limit > MAX_PER_PAGE {
		pages := int(math.Ceil(float64(limit) / float64(MAX_PER_PAGE)))
		for page := 1; page <= pages; page++ {
			utils.ClearPrint(fmt.Sprintf("Retrieving pull requests.. %d/%d", page, pages))
			options.ListOptions.Page = page
			if page == pages {
				options.ListOptions.PerPage = limit - page*MAX_PER_PAGE - MAX_PER_PAGE
			} else {
				options.ListOptions.PerPage = MAX_PER_PAGE
			}
			prs, _, err := client.PullRequests.List(context.Background(), owner, repo, options)
			if err != nil {
				return pullRequestList, err
			}
			pullRequestList = append(pullRequestList, prs...)
		}
	} else {
		utils.ClearPrint("Retrieving pull requests..")
		prs, _, err := client.PullRequests.List(context.Background(), owner, repo, options)
		if err != nil {
			return prs, err
		}
		pullRequestList = append(pullRequestList, prs...)
	}

	return pullRequestList, nil
}

// GetPullRequestTimeline ...
func GetPullRequestTimeline(client *github.Client, pr *github.PullRequest) ([]*github.Timeline, error) {
	options := &github.ListOptions{
		PerPage: MAX_PER_PAGE,
	}
	timeline, _, err := client.Issues.ListIssueTimeline(context.Background(), *pr.Head.Repo.Owner.Login, *pr.Head.Repo.Name, *pr.Number, options)
	return timeline, err
}

func GetPullRequestReviews(client *github.Client, owner string, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	options := &github.ListOptions{
		PerPage: MAX_PER_PAGE,
	}
	reviews, _, err := client.PullRequests.ListReviews(context.Background(), owner, repo, prNumber, options)
	return reviews, err
}

func GetPullRequest(client *github.Client, owner string, repo string, prNumber int) (*github.PullRequest, error) {
	pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, prNumber)
	return pr, err
}
