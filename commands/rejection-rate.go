package commands

import (
	"fmt"

	"github.com/aogz/perforator/gh"
)

// RejectedStatus is status that marks PR as rejected
const RejectedStatus = "CHANGES_REQUESTED"

// UserRejectionRateStatistic is used to store data needed for calculating rejection rate per user
type UserRejectionRateStatistic struct {
	Username string
	Total    int
	Rejected int
	Rate     float64
}

// RejectionRate shows percentage of rejected PRs
func RejectionRate(owner string, repo string, limit int) {
	stats := getStatsByUser(owner, repo, limit)
	fmt.Println("\nResults:")
	if stats != nil {
		sumRate := 0.0
		ratesCount := 0
		for username, userStats := range stats {
			rate := calculateRejectionRatePerUser(userStats)
			fmt.Printf("@%s's rejection rate is: %.2f%% (%d/%d)\n", username, rate, userStats.Rejected, userStats.Total)
			sumRate += rate
			ratesCount += 1
		}

		aggregatedRate := float64(sumRate) / float64(ratesCount)
		fmt.Printf("\nAggregated rejection rate is: %.2f%%\n", aggregatedRate)
	}
}

func getStatsByUser(owner string, repo string, limit int) map[string]UserRejectionRateStatistic {
	client := gh.GetClient()

	prs, err := gh.GetPRs(client, owner, repo, limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	statsByUser := map[string]UserRejectionRateStatistic{}
	for i, pr := range prs {
		username := *pr.User.Login
		prNumber := *pr.Number
		fmt.Printf("%d/%d Processing PR #%d created by %s\n", i+1, limit, prNumber, username)

		if _, ok := statsByUser[username]; !ok {
			statsByUser[username] = UserRejectionRateStatistic{
				Username: username,
			}
		}

		userStats := statsByUser[username]
		userStats.Total += 1
		reviews, err := gh.GetPullRequestReviews(client, owner, repo, prNumber)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		for _, review := range reviews {
			fmt.Printf("\t%s %s\n", *review.User.Login, *review.State)
			if *review.State == RejectedStatus {
				userStats.Rejected += 1
				break
			}
		}
		statsByUser[username] = userStats
	}

	return statsByUser
}

func calculateRejectionRatePerUser(rejectionStats UserRejectionRateStatistic) float64 {
	return float64(rejectionStats.Rejected) / float64(rejectionStats.Total) * 100
}
