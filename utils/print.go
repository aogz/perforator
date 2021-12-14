package utils

import (
	"fmt"
	"os"

	"github.com/google/go-github/v40/github"
)

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

// ClearPrint clears console and prints the line
func ClearPrint(text string) {
	if os.Getenv("PERFORATOR_DEBUG") != "true" {
		clearConsole()
	}
	fmt.Println(text)
}

func ClearPrintPRInfo(i int, limit int, pr *github.PullRequest) {
	ClearPrint(
		fmt.Sprintf(
			"%d/%d Processing PR #%d created at %s by %s",
			i+1, limit, *pr.Number, *pr.CreatedAt, *pr.User.Login,
		),
	)
}

func ClearPrintIssueInfo(i int, limit int, issue *github.Issue) {
	ClearPrint(
		fmt.Sprintf(
			"%d/%d Processing issue #%d created at %s by %s",
			i+1, limit, *issue.Number, *issue.CreatedAt, *issue.User.Login,
		),
	)
}
