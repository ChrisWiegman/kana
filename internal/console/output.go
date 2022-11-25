package console

import (
	"fmt"
	"os"
)

func Error(err error, debugMode bool) {

	if debugMode {
		panic(err)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(1)
}

func Println(output string) {
	fmt.Println(output)
}
