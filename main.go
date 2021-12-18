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
	default:
		utils.PrintHelp()
	}
}
