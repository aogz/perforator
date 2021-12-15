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
		args := utils.AddDefaultArgs(cmd)
		commands.IssueAuthor(args)
	default:
		utils.PrintHelp()
	}
}
