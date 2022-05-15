package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a new environment in the local folder.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the Environment.")
	},
}
