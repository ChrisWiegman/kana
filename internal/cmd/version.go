package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ChrisWiegman/kana-wordpress/internal/console"

	"github.com/spf13/cobra"
)

var Version, Timestamp string

type VersionInfo struct {
	Version, Timestamp string
}

func version(consoleOutput *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays version information for the Kana CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			if consoleOutput.JSON {
				v := VersionInfo{
					Version:   Version,
					Timestamp: Timestamp,
				}

				str, _ := json.Marshal(v)

				fmt.Println(string(str))
			} else {
				consoleOutput.Printf("Version: %s\n", Version)
				consoleOutput.Printf("Build Time: %s\n", Timestamp)
			}
		},
		Args: cobra.NoArgs,
	}

	return cmd
}
