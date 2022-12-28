package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

// Bold outputs the requested text as bold
func Bold(output string) string {
	return aurora.Bold(output).String()
}

// Error displays the error message and a panic if needed.
func Error(err error, debugMode bool) {
	fmt.Fprintf(os.Stderr, "%s %s\n", aurora.Bold(aurora.Red("[Error]")), err)

	if debugMode {
		Println("")
		panic(err)
	}

	os.Exit(1)
}

// Println is a temporary wrapper on fmt.Println
func Println(output string) {
	fmt.Println(output)
}

// PromptConfirm asks the user to confirm output.
func PromptConfirm(promptText string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", promptText, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

// Success displays a formatted success message on successful completion of the command
func Success(output string) {
	fmt.Printf("%s %s\n", aurora.Bold(aurora.Green("[Success]")), output)
}

// Warn displays a formatted warning message
func Warn(output string) {
	fmt.Printf("%s %s\n", aurora.Bold(aurora.Yellow("[Warning]")), output)
}
