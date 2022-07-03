package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/internal/setup"
	"github.com/ChrisWiegman/kana/internal/traefik"
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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setup.SetupApp(controller)
		},
	}

	return cmd

}

func runStart(cmd *cobra.Command, args []string, controller *docker.Controller) {

	site := wordpress.NewSite(controller)

	fmt.Printf("Starting development site: %s\n", site.GetURL(false))

	err := traefik.StartTraefik(controller)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = site.StartWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = site.VerifySite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = site.InstallWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = site.OpenSite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
