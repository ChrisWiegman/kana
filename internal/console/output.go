package console

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

// Error displays the error message and a panic if needed.
func Error(err error, debugMode bool) {

	fmt.Fprintf(os.Stderr, "%s %s\n", aurora.Bold(aurora.Red("[Error]")), err)

	if debugMode {
		Println("")
		panic(err)
	}

	os.Exit(1)
}

func Println(output string) {
	fmt.Println(output)
}

func Success(output string) {
	fmt.Printf("%s %s\n", aurora.Bold(aurora.Green("[Success]")), output)
}
