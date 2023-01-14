package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/aogz/perforator/gh"
	"github.com/aogz/perforator/openai"
	"github.com/aogz/perforator/utils"
)

// CommitsAuthor displays list of commits done in the previous day
func CommitsAuthor(args utils.DefaultArgs, daysAgo int, explain bool) {
	client := gh.GetClient()

	timeDelta := -24 * time.Hour * time.Duration(daysAgo)
	sinceWithTime := time.Now().Add(timeDelta)
	since := time.Date(sinceWithTime.Year(), sinceWithTime.Month(), sinceWithTime.Day(), 0, 0, 0, 0, time.Local)
	until := sinceWithTime.Add(24 * time.Hour)

	for i, contributor := range args.Contributors {
		commits, err := gh.GetCommits(client, args.Owner, args.RepoName, contributor, since, until)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("\n%d: Commits from %s (%d %s)\n", i+1, contributor, since.Day(), since.Month())

		commitsText := ""
		for i, commit := range commits {
			message := strings.ReplaceAll(*commit.Commit.Message, "\n", " ")
			fmt.Printf("\t%d: %s\n", i+1, message)
			fmt.Printf("\t\thttps://github.com/%s/%s/commit/%s\n", args.Owner, args.RepoName, *commit.SHA)

			commitsText += fmt.Sprintf("%d: %s", i+1, message)
		}

		if explain {
			request := fmt.Sprintf("Using the following commit messages, explain what I did in a human friendly way:\n %s", commitsText)
			response, err := openai.DaVinciRequest(request)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("#######")
				fmt.Println(response)
			}
		}
		fmt.Println(" ")
	}
}
