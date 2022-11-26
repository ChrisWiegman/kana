package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Reads all files from the templates directory
// and encodes them as strings literals in templates.go
func main() {

	fs, _ := os.ReadDir("./internal/config/source/")
	out, err := os.Create("./internal/config/constants.go")
	if err != nil {
		os.Exit(1)
	}

	out.Write([]byte("// nolint\npackage config \n\nconst (\n"))

	for _, f := range fs {

		if strings.HasPrefix(f.Name(), ".") {
			continue // Don't include hidden files
		}

		cname := normalize(f.Name())

		out.Write([]byte(cname + " = `"))

		f, err := os.Open("./internal/config/source/" + f.Name())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		io.Copy(out, f)

		out.Write([]byte("`\n"))

	}

	out.Write([]byte(")\n"))
}

func normalize(name string) string {
	return strings.ToUpper(strings.Replace(strings.Replace(name, ".", "_", -1), "-", "_", -1))
}
