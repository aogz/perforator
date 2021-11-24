package main

import (
	"flag"
	"fmt"
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
	var repo *string
	var limit *int

	helpCmd := flag.Bool("help", false, "Help text")
	flag.Parse()

	prRejectionRateCmd := flag.NewFlagSet(rejectionRate, flag.ExitOnError)
	prReviewTimeCmd := flag.NewFlagSet(reviewTime, flag.ExitOnError)

	if *helpCmd {
		fmt.Println("Usage: `$ perforator [command] [--help] [flags]`")
		fmt.Printf("Commands: %s, %s\n", rejectionRate, reviewTime)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("Invalid format. Please use: `$ perforator --help`")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case rejectionRate:
		repo, limit = addDefaultArgs(prRejectionRateCmd)
		owner, repoName := utils.ParseRepo(*repo)
		commands.RejectionRate(owner, repoName, *limit)
	case reviewTime:
		groupBy := prReviewTimeCmd.String("group-by", "reviewer", "Criteria to group by. Accepted values: author or reviewer")
		repo, limit = addDefaultArgs(prReviewTimeCmd)
		owner, repoName := utils.ParseRepo(*repo)
		commands.ReviewTime(owner, repoName, *limit, *groupBy)
	default:
		fmt.Println("Invalid command")
		fmt.Println("Use `$ perforator --help` to get more help")
	}
}
