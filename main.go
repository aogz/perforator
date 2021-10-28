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
)

func main() {
	prRejectionRate := flag.NewFlagSet(rejectionRate, flag.ExitOnError)
	repo := prRejectionRate.String("repo", "", "Repository in format owner/repo, e.g. facebook/react")
	limit := prRejectionRate.Int("limit", 10, "Limit to last X PRs (Max 100)")

	if len(os.Args) < 2 {
		fmt.Println("Please use the following format: `$ perforator rejection-rate -repo=facebook/react -limit 100`")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case rejectionRate:
		prRejectionRate.Parse(os.Args[2:])
		parts := strings.Split(*repo, "/")
		if len(parts) != 2 {
			fmt.Println("Invaid repo format, should be owner/repo")
			os.Exit(1)
		}

		owner := parts[0]
		repoName := parts[1]
		commands.RejectionRate(owner, repoName, *limit)
	}
}
