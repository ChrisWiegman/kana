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

	fs, _ := os.ReadDir("./internal/appSetup/source/")
	out, err := os.Create("./internal/appSetup/constants.go")
	if err != nil {
		fmt.Println(err)
	}

	out.Write([]byte("// nolint\npackage appSetup \n\nconst (\n"))

	for _, f := range fs {

		if strings.HasPrefix(f.Name(), ".") {
			continue // Don't include hidden files
		}

		cname := normalize(f.Name())

		out.Write([]byte(cname + " = `"))

		f, err := os.Open("./internal/appSetup/source/" + f.Name())
		if err != nil {
			fmt.Println(err)
		}

		io.Copy(out, f)

		out.Write([]byte("`\n"))

	}

	out.Write([]byte(")\n"))
}

func normalize(name string) string {
	return strings.ToUpper(strings.Replace(strings.Replace(name, ".", "_", -1), "-", "_", -1))
}
