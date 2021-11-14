package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aogz/perforator/commands"
)

const (
	rejectionRate string = "rejection-rate"
	reviewTime    string = "review-time"
)

func parseRepo(repo string) (string, string) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		fmt.Println("Invaid repo format, should be owner/repo. Type --help for more info.")
		os.Exit(1)
	}

	return parts[0], parts[1]
}

func addDefaultArgs(subcommand *flag.FlagSet) (*string, *int) {
	repo := subcommand.String("repo", "", "Repository in format owner/repo, e.g. facebook/react")
	limit := subcommand.Int("limit", 10, "Limit to last X PRs (Max 100)")

	return repo, limit
}

func main() {
	var repo *string
	var limit *int

	helpCmd := flag.Bool("help", false, "Help text")
	flag.Parse()

	prRejectionRateCmd := flag.NewFlagSet(rejectionRate, flag.ExitOnError)
	prReviewTimeCmd := flag.NewFlagSet(reviewTime, flag.ExitOnError)
	groupBy := prReviewTimeCmd.String("group-by", "reviewer", "Criteria to group by. Accepted values: author or reviewer, default reviewer")

	if *helpCmd {
		fmt.Println("Usage: perforator [command] [--help] [flags]")
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
		prRejectionRateCmd.Parse(os.Args[2:])
		owner, repoName := parseRepo(*repo)
		commands.RejectionRate(owner, repoName, *limit)
	case reviewTime:
		repo, limit = addDefaultArgs(prReviewTimeCmd)
		prReviewTimeCmd.Parse(os.Args[2:])
		owner, repoName := parseRepo(*repo)
		commands.ReviewTime(owner, repoName, *limit, *groupBy)
	default:
		fmt.Println("Invalid command. Please use: `$ perforator --help`")
	}
}
