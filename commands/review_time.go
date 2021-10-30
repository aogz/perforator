package commands

import (
	"fmt"
	"time"

	"github.com/aogz/perforator/gh"
)

// ReviewTime shows percentage of rejected PRs
func ReviewTime(owner string, repo string, limit int, groupBy string) {
	client := gh.GetClient()
	prs, err := gh.GetPRs(client, owner, repo, limit)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stats := map[string][]time.Duration{}
	for i, pr := range prs {
		username := *pr.User.Login
		prNumber := *pr.Number
		createdAt := *pr.CreatedAt
		fmt.Printf("%d/%d Processing PR #%d created by %s\n", i+1, limit, prNumber, username)

		if groupBy == "author" {
			if pr.MergedAt != nil {
				mergedAt := *pr.MergedAt
				inReviewTime := mergedAt.Sub(createdAt)
				stats[username] = append(stats[username], inReviewTime)
			}
		} else {
			reviews, err := gh.GetPullRequestReviews(client, owner, repo, prNumber)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, review := range reviews {
				reviewedBy := *review.User.Login
				reviewedAt := *review.SubmittedAt
				inReviewTime := reviewedAt.Sub(createdAt)
				stats[reviewedBy] = append(stats[reviewedBy], inReviewTime)
			}
		}
	}

	getResults(stats)
}

func getResults(stats map[string][]time.Duration) {
	fmt.Println("-----")
	var sumAggregatedDurationInHours float64
	for username, durations := range stats {
		var durationSum time.Duration
		for _, duration := range durations {
			durationSum += duration
		}

		aggregatedDuration := durationSum.Seconds() / float64(len(durations))
		aggregatedDurationInHours := aggregatedDuration / 60 / 60
		sumAggregatedDurationInHours += aggregatedDurationInHours
		fmt.Printf("@%s's average review time is: %.2f hours\n", username, aggregatedDurationInHours)
	}
	averageAggregatedDurationInHours := sumAggregatedDurationInHours / float64(len(stats))
	fmt.Printf("\nAggregated PR review duration is: %.2f hours (From %d Devs)\n", averageAggregatedDurationInHours, len(stats))

}
