package cmd

import (
	"fmt"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/traefik"
	"github.com/ChrisWiegman/kana/internal/utilities"
	"github.com/ChrisWiegman/kana/internal/wordpress"
	"github.com/spf13/cobra"
)

func newStartCommand(controller *docker.Controller) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {
			runStart(cmd, args, controller)
		},
	}

	return cmd

}

func runStart(cmd *cobra.Command, args []string, controller *docker.Controller) {

	siteURL := fmt.Sprintf("https://%s.%s/", controller.Config.CurrentDirectory, controller.Config.SiteDomain)

	fmt.Printf("Starting development site: %s\n", siteURL)

	err := traefik.NewTraefik(controller)
	if err != nil {
		fmt.Println(err)
	}

	err = wordpress.NewWordPress(controller)
	if err != nil {
		fmt.Println(err)
	}

	utilities.OpenURL(siteURL)

}
