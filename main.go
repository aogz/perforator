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
	issueLabels   string = "issue-labels"
	commitsAuthor string = "commits"
	activity      string = "activity"
	repoActivity  string = "repo-activity"
)

func main() {
	utils.AddHelp()
	utils.ValidateArgs()

	cmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	switch cmd.Name() {
	case rejectionRate:
		args := utils.AddDefaultArgs(cmd)
		commands.RejectionRate(args)
	case reviewTime:
		groupBy := cmd.String("group-by", "reviewer", "Criteria to group by. Accepted values: author or reviewer")
		args := utils.AddDefaultArgs(cmd)
		commands.ReviewTime(args, *groupBy)
	case issueAuthor:
		labels := cmd.String("labels", "", "Comma separated list of labels to filter by")
		state := cmd.String("state", "all", "State of the issues to filter by. Accepted values: all, open, closed")
		args := utils.AddDefaultArgs(cmd)
		commands.IssueAuthor(args, utils.ParseCommaSeparatedValue(*labels), *state)
	case issueLabels:
		labels := cmd.String("labels", "", "Comma separated list of labels to filter by")
		state := cmd.String("state", "all", "State of the issues to filter by. Accepted values: all, open, closed")
		args := utils.AddDefaultArgs(cmd)
		commands.IssueLabels(args, utils.ParseCommaSeparatedValue(*labels), *state)
	case commitsAuthor:
		daysAgo := cmd.Int("days-ago", 1, "Days ago (default: 1 (yesterday))")
		explain := cmd.Bool("explain", false, "Explain commit messages in a human friendly way using Chat GPT")
		args := utils.AddDefaultArgs(cmd)
		commands.CommitsAuthor(args, *daysAgo, *explain)
	case activity:
		sinceDaysAgo := cmd.Int("since-days-ago", 1, "Since days ago (default: 1 (since yesterday))")
		explain := cmd.Bool("explain", false, "Explain messages in a human friendly way using Chat GPT")
		args := utils.AddDefaultArgs(cmd)
		commands.UsersActivity(args, *sinceDaysAgo, *explain)
	case repoActivity:
		sinceDaysAgo := cmd.Int("since-days-ago", 1, "Since days ago (default: 1 (since yesterday))")
		explain := cmd.Bool("explain", false, "Explain messages in a human friendly way using Chat GPT")
		args := utils.AddDefaultArgs(cmd)
		commands.RepositoryActivity(args, *sinceDaysAgo, *explain)
	default:
		utils.PrintHelp()
	}
}
