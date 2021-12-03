package main

import (
	"flag"
	"os"

	"github.com/aogz/perforator/commands"
	"github.com/aogz/perforator/utils"
)

const (
	rejectionRate string = "rejection-rate"
	reviewTime    string = "review-time"
)

func main() {
	utils.AddHelp()
	utils.ValidateArgs()

	command := os.Args[1]
	switch command {
	case rejectionRate:
		prRejectionRateCmd := flag.NewFlagSet(rejectionRate, flag.ExitOnError)
		args := utils.AddDefaultArgs(prRejectionRateCmd)
		owner, repoName := utils.ParseRepo(args.Repo)
		commands.RejectionRate(owner, repoName, args.Limit, args.Skip)
	case reviewTime:
		prReviewTimeCmd := flag.NewFlagSet(reviewTime, flag.ExitOnError)
		groupBy := prReviewTimeCmd.String("group-by", "reviewer", "Criteria to group by. Accepted values: author or reviewer")
		args := utils.AddDefaultArgs(prReviewTimeCmd)
		owner, repoName := utils.ParseRepo(args.Repo)
		commands.ReviewTime(owner, repoName, args.Limit, args.Skip, *groupBy)
	default:
		utils.PrintHelp()
	}
}
