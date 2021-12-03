package commands

import (
	"fmt"

	"github.com/aogz/perforator/gh"
	"github.com/aogz/perforator/utils"
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
func RejectionRate(owner string, repo string, limit int, skip int) {
	stats := getStatsByUser(owner, repo, limit, skip)
	utils.ClearPrint("-----")
	if stats != nil {
		aggregatedRate := calculateAggregatedRate(stats)
		fmt.Printf("\nAggregated rejection rate is: %.2f%% (From %d Devs)\n", aggregatedRate, len(stats))
	}
}

// Helper methods
func getStatsByUser(owner string, repo string, limit int, skip int) map[string]UserRejectionRateStatistic {
	client := gh.GetClient()

	prs, err := gh.GetPRs(client, owner, repo, limit, skip)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	statsByUser := map[string]UserRejectionRateStatistic{}
	for i, pr := range prs {
		username := *pr.User.Login
		prNumber := *pr.Number

		utils.ClearPrintPRInfo(i, limit, pr)
		tryCreateUserStats(statsByUser, username)
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

		userStats.Rate = calculateRejectionRatePerUser(userStats)
		statsByUser[username] = userStats
	}

	return statsByUser
}

func tryCreateUserStats(stats map[string]UserRejectionRateStatistic, username string) {
	if _, ok := stats[username]; !ok {
		stats[username] = UserRejectionRateStatistic{
			Username: username,
		}
	}
}

func calculateRejectionRatePerUser(rejectionStats UserRejectionRateStatistic) float64 {
	return float64(rejectionStats.Rejected) / float64(rejectionStats.Total) * 100
}

func calculateAggregatedRate(stats map[string]UserRejectionRateStatistic) float64 {
	sumRate := 0.0
	for username, userStats := range stats {
		rate := calculateRejectionRatePerUser(userStats)
		fmt.Printf("@%s's rejection rate is: %.2f%% (%d/%d)\n", username, rate, userStats.Rejected, userStats.Total)
		sumRate += rate
	}

	return float64(sumRate) / float64(len(stats))
}
