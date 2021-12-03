package utils

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type DefaultArgs struct {
	Repo  string
	Limit int
	Skip  int
}

// ParseRepo parses repo string and terminates the script if the format is incalid. Returns author and repo strings.
func ParseRepo(repo string) (string, string) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		fmt.Println("Invaid repo format, should be owner/repo. Type --help for more info.")
		os.Exit(1)
	}

	return parts[0], parts[1]
}

func AddDefaultArgs(subcommand *flag.FlagSet) DefaultArgs {
	repo := subcommand.String("repo", "", "Repository in format owner/repo, e.g. facebook/react")
	limit := subcommand.Int("limit", 10, "Limit to last X PRs")
	skip := subcommand.Int("skip", 0, "Skip first X PRs")
	subcommand.Parse(os.Args[2:])

	return DefaultArgs{
		Repo:  *repo,
		Limit: *limit,
		Skip:  *skip,
	}
}
