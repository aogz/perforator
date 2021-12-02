package commands

import (
	"fmt"
	"time"

	"github.com/aogz/perforator/gh"
	"github.com/aogz/perforator/utils"
	"github.com/google/go-github/v40/github"
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
		stats = calculateReviewTimeByAuthor(client, stats, prs, limit)
	} else {
		stats = calculateReviewTimeByReviewer(client, stats, prs, limit)
	}

	calculateAggregatedResultsPerUser(stats)
}

func calculateReviewTimeByAuthor(client *github.Client, stats map[string][]time.Duration, prs []*github.PullRequest, limit int) map[string][]time.Duration {
	for i, pr := range prs {
		if pr.MergedAt == nil {
			fmt.Println("PR is not merged yet, skipping..")
			continue
		}
		username := *pr.User.Login
		utils.ClearPrint(fmt.Sprintf("%d/%d Processing PR #%d created at %s by %s", i+1, limit, *pr.Number, *pr.CreatedAt, username))
		inReviewTime := calculatePRInReviewTime(client, pr)
		if inReviewTime > 0 {
			stats[username] = append(stats[username], inReviewTime)
		}
	}

	return stats
}

func calculateReviewTimeByReviewer(client *github.Client, stats map[string][]time.Duration, prs []*github.PullRequest, limit int) map[string][]time.Duration {
	for i, pr := range prs {
		utils.ClearPrint(fmt.Sprintf("%d/%d Processing PR #%d created at %s by %s", i+1, limit, *pr.Number, *pr.CreatedAt, *pr.User.Login))
		if pr.MergedAt == nil {
			fmt.Println("PR is not merged yet, skipping..")
			continue
		}

		timeline, err := gh.GetPullRequestTimeline(client, pr)
		if err != nil {
			panic(err)
		}

		reviewPerUser := make(map[string]map[string]interface{})
		for _, event := range timeline {
			eventType := *event.Event

			switch eventType {
			case "review_requested":
				eventCreatedAt := *event.CreatedAt
				reviewer := *event.Reviewer.Login
				fmt.Println("Review requested from:", reviewer)
				if _, ok := reviewPerUser[reviewer]; !ok {
					reviewPerUser[reviewer] = map[string]interface{}{
						"inReviewTime": time.Duration(0),
					}
				}
				reviewPerUser[reviewer]["previousReviewPeriodStartTime"] = eventCreatedAt
			case "review_request_removed", "reviewed":
				var eventTime time.Time
				var reviewer string
				if eventType == "reviewed" {
					reviewer = *event.User.Login
					eventTime = *event.SubmittedAt
					if _, ok := reviewPerUser[reviewer]; !ok {
						fmt.Println("Skipping review from, not requested:", reviewer)
						continue
					}
					reviewPerUser[reviewer]["reviewed"] = true
				} else {
					eventTime = *event.CreatedAt
					reviewer = *event.Reviewer.Login
				}
				timeSinceLastEvent := eventTime.Sub(reviewPerUser[reviewer]["previousReviewPeriodStartTime"].(time.Time))
				currentDuration := reviewPerUser[reviewer]["inReviewTime"].(time.Duration)
				reviewPerUser[reviewer]["inReviewTime"] = currentDuration + timeSinceLastEvent
				reviewPerUser[reviewer]["previousReviewPeriodStartTime"] = eventTime
			}
		}

		for reviewer, reviewerStats := range reviewPerUser {
			fmt.Printf("\tProcessing %s's review", reviewer)
			if isReviewed, ok := reviewPerUser[reviewer]["reviewed"]; !ok || !isReviewed.(bool) {
				mergetAt := *pr.MergedAt
				lastEvent := reviewPerUser[reviewer]["previousReviewPeriodStartTime"].(time.Time)
				currentDuration := reviewPerUser[reviewer]["inReviewTime"].(time.Duration)
				reviewerStats["inReviewTime"] = currentDuration + mergetAt.Sub(lastEvent)
			}
			stats[reviewer] = append(stats[reviewer], reviewerStats["inReviewTime"].(time.Duration))
		}
	}

	return stats
}

func calculateAggregatedResultsPerUser(stats map[string][]time.Duration) float64 {
	utils.ClearPrint("-----")
	var sumAggregatedDurationInHours float64
	for username, durations := range stats {
		averageReviewTime := calculateAverageReviewTimePerUser(durations)
		sumAggregatedDurationInHours += averageReviewTime
		fmt.Printf("@%s's average review time is: %.2f hours (%d PRs)\n", username, averageReviewTime, len(durations))
	}
	averageAggregatedDurationInHours := calculateAggregatedAverageReviewTime(sumAggregatedDurationInHours, len(stats))
	fmt.Printf("\nAggregated PR review duration is: %.2f hours (From %d Devs)\n", averageAggregatedDurationInHours, len(stats))
	return averageAggregatedDurationInHours
}

func calculatePRInReviewTime(client *github.Client, pr *github.PullRequest) time.Duration {
	var inReviewTime time.Duration
	previousReviewPeriodStartTime := *pr.CreatedAt
	timeline, err := gh.GetPullRequestTimeline(client, pr)
	if err != nil {
		panic(err.Error())
	}

	for _, event := range timeline {
		eventType := *event.Event
		switch eventType {
		case "ready_for_review":
			previousReviewPeriodStartTime = *event.CreatedAt
		case "convert_to_draft":
			eventCreatedAt := *event.CreatedAt
			inReviewTime += eventCreatedAt.Sub(previousReviewPeriodStartTime)
			previousReviewPeriodStartTime = *event.CreatedAt
		}
	}

	mergedAt := *pr.MergedAt
	inReviewTime = mergedAt.Sub(previousReviewPeriodStartTime)
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
