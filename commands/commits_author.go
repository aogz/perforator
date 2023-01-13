package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/aogz/perforator/gh"
	"github.com/aogz/perforator/utils"
)

// CommitsAuthor displays list of commits done in the previous day
func CommitsAuthor(args utils.DefaultArgs, daysAgo int) {
	client := gh.GetClient()

	timeDelta := -24 * time.Hour * time.Duration(daysAgo)
	sinceWithTime := time.Now().Add(timeDelta)
	since := time.Date(sinceWithTime.Year(), sinceWithTime.Month(), sinceWithTime.Day(), 0, 0, 0, 0, time.Local)
	untilWithTime := time.Now()
	until := time.Date(untilWithTime.Year(), untilWithTime.Month(), untilWithTime.Day(), 23, 59, 59, 59, time.Local)

	for i, contributor := range args.Contributors {
		commits, err := gh.GetCommits(client, args.Owner, args.RepoName, contributor, since, until)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("\n%d: Commits from %s (%d %s)\n", i+1, contributor, since.Day(), since.Month())
		for i, commit := range commits {
			message := strings.ReplaceAll(*commit.Commit.Message, "\n", " ")
			fmt.Printf("\t%d: %s\n", i+1, message)
			fmt.Printf("\t\thttps://github.com/%s/%s/commit/%s\n", args.Owner, args.RepoName, *commit.SHA)
		}
		fmt.Println(" ")
	}
}
