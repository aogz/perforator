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
	issueAuthor   string = "issue-author"
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
		commands.RejectionRate(owner, repoName, args.Limit, args.Skip, args.Contributors)
	case reviewTime:
		prReviewTimeCmd := flag.NewFlagSet(reviewTime, flag.ExitOnError)
		groupBy := prReviewTimeCmd.String("group-by", "reviewer", "Criteria to group by. Accepted values: author or reviewer")
		args := utils.AddDefaultArgs(prReviewTimeCmd)
		owner, repoName := utils.ParseRepo(args.Repo)
		commands.ReviewTime(owner, repoName, args.Limit, args.Skip, *groupBy, args.Contributors)
	case issueAuthor:
		issueAuthorCmd := flag.NewFlagSet(issueAuthor, flag.ExitOnError)
		args := utils.AddDefaultArgs(issueAuthorCmd)
		owner, repoName := utils.ParseRepo(args.Repo)
		commands.IssueAuthor(owner, repoName, args.Limit, args.Skip, args.Contributors)
	default:
		utils.PrintHelp()
	}
}
