package utils

import "fmt"

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

// ClearPrint clears console and prints the line
func ClearPrint(text string) {
	clearConsole()
	fmt.Println(text)
}
