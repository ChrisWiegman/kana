package main

//go:generate go run scripts/generate.go

import "github.com/ChrisWiegman/kana-cli/internal/cmd"

func main() {
	cmd.Execute()
}
