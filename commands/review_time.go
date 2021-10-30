package commands

import (
	"fmt"
	"time"

	"github.com/aogz/perforator/gh"
	"github.com/google/go-github/v39/github"
)

// ReviewTime shows average review time per pr author or reviewer
func ReviewTime(owner string, repo string, limit int, groupBy string) {
	client := gh.GetClient()
	prs, err := gh.GetPRs(client, owner, repo, limit)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stats := map[string][]time.Duration{}
	if groupBy == "author" {
		stats = calculateReviewTimeByAuthor(stats, prs)
	} else {
		stats = calculateReviewTimeByReviewer(stats, prs, client)
	}

	calculateAggregatedResultsPerUser(stats)
}

func calculateReviewTimeByAuthor(stats map[string][]time.Duration, prs []*github.PullRequest) map[string][]time.Duration {
	for i, pr := range prs {
		username := *pr.User.Login
		fmt.Printf("%d Processing PR #%d created by %s\n", i+1, *pr.Number, username)

		inReviewTime := calculatePRInReviewTime(pr)
		if inReviewTime > 0 {
			stats[username] = append(stats[username], inReviewTime)
		}
	}

	return stats
}

func calculateReviewTimeByReviewer(stats map[string][]time.Duration, prs []*github.PullRequest, client *github.Client) map[string][]time.Duration {
	for i, pr := range prs {
		fmt.Printf("%d Processing PR #%d created by %s\n", i+1, *pr.Number, *pr.User.Login)
		reviews, err := gh.GetPullRequestReviews(client, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Number)
		if err != nil {
			fmt.Println(err.Error())
			return stats
		}

		for _, review := range reviews {
			inReviewTime := calculateUserReviewTime(pr, review)
			stats[*review.User.Login] = append(stats[*review.User.Login], inReviewTime)
		}
	}

	return stats
}

func calculateAggregatedResultsPerUser(stats map[string][]time.Duration) float64 {
	fmt.Println("-----")
	var sumAggregatedDurationInHours float64
	for username, durations := range stats {
		averageReviewTime := calculateAverageReviewTimePerUser(durations)
		sumAggregatedDurationInHours += averageReviewTime
		fmt.Printf("@%s's average review time is: %.2f hours\n", username, averageReviewTime)
	}
	averageAggregatedDurationInHours := calculateAggregatedAverageReviewTime(sumAggregatedDurationInHours, len(stats))
	fmt.Printf("\nAggregated PR review duration is: %.2f hours (From %d Devs)\n", averageAggregatedDurationInHours, len(stats))
	return averageAggregatedDurationInHours
}

func calculateUserReviewTime(pr *github.PullRequest, review *github.PullRequestReview) time.Duration {
	reviewedAt := *review.SubmittedAt
	return reviewedAt.Sub(*pr.CreatedAt)
}

func calculatePRInReviewTime(pr *github.PullRequest) time.Duration {
	var inReviewTime time.Duration
	createdAt := *pr.CreatedAt
	if pr.MergedAt != nil {
		mergedAt := *pr.MergedAt
		inReviewTime = mergedAt.Sub(createdAt)
	}

	return inReviewTime
}

// calculateAverageReviewTimePerUser returns average time in seconds
func calculateAverageReviewTimePerUser(durations []time.Duration) float64 {
	var durationSum time.Duration
	for _, duration := range durations {
		durationSum += duration
	}

	return durationSum.Hours() / float64(len(durations))
}

func calculateAggregatedAverageReviewTime(sumAggregatedDurationInHours float64, count int) float64 {
	return sumAggregatedDurationInHours / float64(count)
}
