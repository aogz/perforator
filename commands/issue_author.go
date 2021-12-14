package commands

import (
	"fmt"

	"github.com/aogz/perforator/gh"
	"github.com/aogz/perforator/utils"
)

// IssueAuthor shows number of issues created by each author
func IssueAuthor(owner string, repo string, limit int, skip int, contributors []string) {
	client := gh.GetClient()
	issues, err := gh.GetIssuesByRepo(client, owner, repo, limit, skip)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stats := map[string]int{}
	for i, issue := range issues {
		author := *issue.User.Login
		if len(contributors) > 0 && !utils.Contains(contributors, author) {
			fmt.Printf("\t%s is not in the list of contributors (--only param), skipping..\n", author)
			continue
		}

		utils.ClearPrintIssueInfo(i, limit, issue)
		stats[author]++
	}

	utils.ClearPrint("-----")
	for author, count := range stats {
		fmt.Printf("Issues created by %s: %d\n", author, count)
	}
}
