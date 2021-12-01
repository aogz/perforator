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

func addDefaultArgs(subcommand *flag.FlagSet) (*string, *int) {
	repo := subcommand.String("repo", "", "Repository in format owner/repo, e.g. facebook/react")
	limit := subcommand.Int("limit", 10, "Limit to last X PRs")
	subcommand.Parse(os.Args[2:])

	return repo, limit
}

func main() {
	utils.AddHelp()
	utils.ValidateArgs()

	command := os.Args[1]
	switch command {
	case rejectionRate:
		prRejectionRateCmd := flag.NewFlagSet(rejectionRate, flag.ExitOnError)
		repo, limit := addDefaultArgs(prRejectionRateCmd)
		owner, repoName := utils.ParseRepo(*repo)
		commands.RejectionRate(owner, repoName, *limit)
	case reviewTime:
		prReviewTimeCmd := flag.NewFlagSet(reviewTime, flag.ExitOnError)
		groupBy := prReviewTimeCmd.String("group-by", "reviewer", "Criteria to group by. Accepted values: author or reviewer")
		repo, limit := addDefaultArgs(prReviewTimeCmd)
		owner, repoName := utils.ParseRepo(*repo)
		commands.ReviewTime(owner, repoName, *limit, *groupBy)
	default:
		utils.PrintHelp()
	}
}
