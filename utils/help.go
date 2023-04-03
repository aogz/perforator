package utils

import (
	"flag"
	"fmt"
	"os"
)

const HELP_TEXT = `Perforator is a CLI tool to track and analyze performance and statistics based on the information from GitHub Pull Requests.

Usage:
	perforator [command] [flags]

Commands:
	rejection-rate	Calculate the aggregated rejection rate of pull requests grouped by contibutor.
	review-time		Calculate the aggregated review time of pull requests grouped by author/reviewer.
	issue-author	Calculate number of issues created by author.
	issue-labels	Calculate number of issues grouped by labels.
	commits			Returns commits made by a specified contributor.

Flags:
	--help: Prints help information.
`

func PrintHelp() {
	fmt.Println(HELP_TEXT)
	os.Exit(0)
}

func AddHelp() {
	helpCmd := flag.Bool("help", false, "Prints help information")
	hCmd := flag.Bool("h", false, "Prints help information")
	flag.Parse()
	if *helpCmd || *hCmd {
		PrintHelp()
	}
}

func ValidateArgs() {
	if len(os.Args) < 2 {
		PrintHelp()
	}
}
