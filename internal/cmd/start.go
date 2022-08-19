package cmd

import (
	"fmt"
	"os"

	"github.com/ChrisWiegman/kana/internal/appConfig"
	"github.com/ChrisWiegman/kana/internal/site"
	"github.com/ChrisWiegman/kana/internal/traefik"

	"github.com/spf13/cobra"
)

var flagXdebug bool
var flagLocal bool
var flagIsTheme bool
var flagIsPlugin bool

func newStartCommand(site *site.Site) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a new environment in the local folder.",
		Run: func(cmd *cobra.Command, args []string) {
			runStart(cmd, args, site)
		},
		Args: cobra.NoArgs,
	}

	cmd.Flags().BoolVar(&flagXdebug, "xdebug", false, "Enable Xdebug when starting the container.")
	cmd.Flags().BoolVar(&flagIsPlugin, "plugin", false, "Run the site as a plugin using the current folder as the plugin source.")
	cmd.Flags().BoolVar(&flagIsTheme, "theme", false, "Run the site as a theme using the current folder as the theme source.")
	cmd.Flags().BoolVar(&flagLocal, "local", false, "Installs the WordPress files in your current path at ./wordpress instead of the global app path.")

	return cmd

}

func runStart(cmd *cobra.Command, args []string, kanaSite *site.Site) {

	if flagIsPlugin && flagIsTheme {
		fmt.Println(fmt.Errorf("you have set both the plugin and theme flags. Please choose only one option"))
		os.Exit(1)
	}

	dynamicConfig, err := appConfig.GetDynamicContent(kanaSite.StaticConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	kanaSite.DynamicConfig = dynamicConfig

	if kanaSite.IsSiteRunning() {
		fmt.Println("Site is already running. Please stop your site before running the start command")
		os.Exit(1)
	}

	startFlags := site.SiteFlags{
		Xdebug:   flagXdebug,
		IsTheme:  flagIsTheme,
		IsPlugin: flagIsPlugin,
		Local:    flagLocal,
	}

	kanaSite.ProcessSiteFlags(cmd, startFlags)

	fmt.Printf("Starting development site: %s\n", kanaSite.GetURL(false))

	traefikClient, err := traefik.NewTraefik(kanaSite.StaticConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = traefikClient.StartTraefik()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = kanaSite.StartWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = kanaSite.VerifySite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = kanaSite.InstallWordPress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = kanaSite.InstallXdebug()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = kanaSite.OpenSite()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
