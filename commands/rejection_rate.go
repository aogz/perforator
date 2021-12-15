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
func RejectionRate(args utils.DefaultArgs) {
	stats := getStatsByUser(args)
	utils.ClearPrint("-----")
	if stats != nil {
		aggregatedRate := calculateAggregatedRate(stats)
		fmt.Printf("\nAggregated rejection rate is: %.2f%% (From %d Devs)\n", aggregatedRate, len(stats))
	}
}

// Helper methods
func getStatsByUser(args utils.DefaultArgs) map[string]UserRejectionRateStatistic {
	client := gh.GetClient()

	prs, err := gh.GetPRs(client, args.Owner, args.RepoName, args.Limit, args.Skip)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	statsByUser := map[string]UserRejectionRateStatistic{}
	for i, pr := range prs {
		username := *pr.User.Login
		if len(args.Contributors) > 0 && !utils.Contains(args.Contributors, username) {
			fmt.Printf("\t%s is not in the list of contributors (--only param), skipping..\n", username)
			continue
		}
		prNumber := *pr.Number

		utils.ClearPrintPRInfo(i, args.Limit, pr)
		tryCreateUserStats(statsByUser, username)
		userStats := statsByUser[username]
		userStats.Total += 1
		reviews, err := gh.GetPullRequestReviews(client, args.Owner, args.RepoName, prNumber)
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
