package utils

import (
	"fmt"
	"os"
	"strings"
)

// ParseRepo parses repo string and terminates the script if the format is incalid. Returns author and repo strings.
func ParseRepo(repo string) (string, string) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		fmt.Println("Invaid repo format, should be owner/repo. Type --help for more info.")
		os.Exit(1)
	}

	return parts[0], parts[1]
}
