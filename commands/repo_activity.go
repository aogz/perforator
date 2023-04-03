package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/aogz/perforator/gh"
	"github.com/aogz/perforator/utils"
	"github.com/google/go-github/v40/github"
)

type UserActivity struct {
	PRsMerged     []string
	PRReviews     []string
	IssuesCreated []string
	Commits       []string
}

// RepositoryActivity returns activity for a specific organization
func RepositoryActivity(args utils.DefaultArgs, sinceDaysAgo int, explain bool) {
	client := gh.GetClient()
	activity, err := gh.GetRepositoryActivity(client, args.Owner, args.RepoName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	timeDelta := -24 * time.Hour * time.Duration(sinceDaysAgo)
	sinceWithTime := time.Now().Add(timeDelta)
	since := time.Date(sinceWithTime.Year(), sinceWithTime.Month(), sinceWithTime.Day(), 0, 0, 0, 0, time.Local)

	usersActivity := make(map[string]UserActivity)
	for _, event := range activity {
		if fmt.Sprintf("%s/%s", args.Owner, args.RepoName) != *event.Repo.Name {
			continue
		}

		payload, err := event.ParsePayload()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// check since
		if (*event.CreatedAt).Before(since) {
			break
		}

		if !utils.Contains(args.Contributors, *event.Actor.Login) {
			continue
		}

		if _, ok := usersActivity[*event.Actor.Login]; !ok {
			usersActivity[*event.Actor.Login] = UserActivity{}
		}

		userActivity := usersActivity[*event.Actor.Login]
		if *event.Type == "PullRequestReviewEvent" {
			prReview := fmt.Sprintf("%s '%s'", strings.Title(*payload.(*github.PullRequestReviewEvent).Review.State), *payload.(*github.PullRequestReviewEvent).PullRequest.Title)
			userActivity.PRReviews = append(usersActivity[*event.Actor.Login].PRReviews, prReview)
		} else if *event.Type == "IssuesEvent" {
			issueCreated := fmt.Sprintf("%s '%s'", strings.Title(*payload.(*github.IssuesEvent).Action), *payload.(*github.IssuesEvent).Issue.Title)
			userActivity.IssuesCreated = append(usersActivity[*event.Actor.Login].IssuesCreated, issueCreated)
		} else if *event.Type == "PullRequestEvent" {
			prMerged := fmt.Sprintf("%s '%s'", strings.Title(*payload.(*github.PullRequestEvent).Action), *payload.(*github.PullRequestEvent).PullRequest.Title)
			userActivity.PRsMerged = append(usersActivity[*event.Actor.Login].PRsMerged, prMerged)
		} else if *event.Type == "PushEvent" {
			for _, commit := range payload.(*github.PushEvent).Commits {
				if !strings.Contains(*commit.Message, "Merge pull request") {
					userActivity.Commits = append(usersActivity[*event.Actor.Login].Commits, *commit.Message)
				}
			}
		}

		usersActivity[*event.Actor.Login] = userActivity
	}

	for contributor, activity := range usersActivity {
		fmt.Printf("%s:\n", contributor)
		fmt.Println(utils.PrintList("PRs", activity.PRsMerged))
		fmt.Println(utils.PrintList("PR Reviews", activity.PRReviews))
		fmt.Println(utils.PrintList("Issues Created", activity.IssuesCreated))
		fmt.Println(utils.PrintList("Commits", activity.Commits))
		fmt.Println("")
	}
}
