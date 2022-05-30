package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Reads all files from the templates directory
// and encodes them as strings literals in templates.go
func main() {
	fs, _ := ioutil.ReadDir("./internal/setup/source/")
	out, err := os.Create("./internal/setup/constants.go")
	if err != nil {
		fmt.Println(err)
	}
	out.Write([]byte("// nolint\npackage setup \n\nconst (\n"))
	for _, f := range fs {
		if strings.HasPrefix(f.Name(), ".") {
			// Don't include hidden files
			continue
		}
		cname := normalize(f.Name())
		out.Write([]byte(cname + " = `"))
		f, err := os.Open("./internal/setup/source/" + f.Name())
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
