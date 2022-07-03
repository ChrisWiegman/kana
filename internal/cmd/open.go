package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/setup"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

func newOpenCommand(controller *docker.Controller) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open the current site in your browser.",
		Run: func(cmd *cobra.Command, args []string) {
			runOpen(cmd, args, controller)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setup.SetupApp(controller)
		},
	}

	return cmd

}

func runOpen(cmd *cobra.Command, args []string, controller *docker.Controller) {

	site := site.NewSite(controller.Config)

	err := site.OpenSite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
