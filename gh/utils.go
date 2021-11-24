package gh

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/aogz/perforator/utils"
	"github.com/google/go-github/v39/github"
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

func GetPullRequestTimeline(client *github.Client, owner string, repo string, prNumber int) ([]*github.Timeline, error) {
	timeline, _, err := client.Issues.ListIssueTimeline(context.Background(), owner, repo, prNumber, nil)
	return timeline, err
}

func GetPullRequestReviews(client *github.Client, owner string, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	reviews, _, err := client.PullRequests.ListReviews(context.Background(), owner, repo, prNumber, nil)
	return reviews, err
}

func GetPullRequestReviewStartTime(client *github.Client, owner string, repo string, prNumber int) (time.Time, error) {
	pr, _, _ := client.PullRequests.Get(context.Background(), owner, repo, prNumber)
	timeline, err := GetPullRequestTimeline(client, owner, repo, prNumber)
	if err != nil {
		panic(err)
	}

	isReadyForReview := true
	reviewStartTime := *pr.CreatedAt
	for _, event := range timeline {
		eventType := *event.Event
		switch eventType {
		case "ready_for_review":
			isReadyForReview = true
			reviewStartTime = *event.CreatedAt
		case "convert_to_draft":
			isReadyForReview = false
		}

		fmt.Println(*event.Event, *event.CreatedAt)
	}

	if !isReadyForReview {
		return reviewStartTime, errors.New("PR is not ready for review")
	}

	return reviewStartTime, nil
}
