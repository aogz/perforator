package commands

import (
	"fmt"

	"github.com/aogz/perforator/gh"
	"github.com/aogz/perforator/utils"
)

// IssueLabels shows number of issues created by each author
func IssueLabels(args utils.DefaultArgs) {
	client := gh.GetClient()
	issues, err := gh.GetIssuesByRepo(client, args.Owner, args.RepoName, args.Limit, args.Skip)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	stats := map[string]int{}
	for i, issue := range issues {
		author := *issue.User.Login
		if len(args.Contributors) > 0 && !utils.Contains(args.Contributors, author) {
			fmt.Printf("\t%s is not in the list of contributors (--only param), skipping..\n", author)
			continue
		}

		utils.ClearPrintIssueInfo(i, args.Limit, issue)
		for _, label := range issue.Labels {
			stats[*label.Name]++
		}
	}

	utils.ClearPrint("-----")
	for label, count := range stats {
		fmt.Printf("Issues with %s label: %d\n", label, count)
	}
}
