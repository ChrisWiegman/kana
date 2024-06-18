package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChrisWiegman/kana/internal/helpers"

	"github.com/spf13/cobra"
)

func loadSiteOptions(settings *Settings, cmd *cobra.Command) error {
	siteSettings := Site{}

	isNamedSite, siteName, err := getSiteName(settings, cmd)
	if err != nil {
		return err
	}

	siteSettings.Name = siteName
	siteSettings.IsNamed = isNamedSite

	// We can set the site directory here now that we have the correct name.
	settings.directories.Site = filepath.Join(settings.directories.App, "sites", siteName)

	_, err = os.Stat(settings.directories.Site)
	if err != nil && os.IsNotExist(err) {
		if os.IsNotExist(err) {
			siteSettings.IsNew = true
		} else {
			return err
		}
	}

	settings.site = siteSettings

	return nil
}

func getSiteName(settings *Settings, cmd *cobra.Command) (isNamedSite bool, siteName string, err error) {
	siteName = helpers.SanitizeSiteName(filepath.Base(settings.directories.Working))

	isStartCommand := cmd.Use == "start"

	// Don't run this on commands that wouldn't possibly use it.
	if cmd.Use == "config" || cmd.Use == "version" || cmd.Use == "help" {
		return isNamedSite, siteName, nil
	}

	// Validate the flags that could make this an invalid site
	if isStartCommand && cmd.Flags().Lookup("multisite").Changed {
		multisiteValue, err := cmd.Flags().GetString("multisite")
		if !helpers.IsValidString(multisiteValue, validMultisiteTypes) || err != nil {
			return isNamedSite, siteName,
				fmt.Errorf("the multisite type, %s, is not a valid type. You must use either `none`, `subdomain` or `subdirectory`", multisiteValue)
		}
	}

	if isStartCommand && cmd.Flags().Lookup("type").Changed {
		typeValue, err := cmd.Flags().GetString("type")
		if !helpers.IsValidString(typeValue, validTypes) || err != nil {
			return isNamedSite, siteName,
				fmt.Errorf("the type, %s, is not valid. Only a `site` is valid with named sites", typeValue)
		}
	}

	if isStartCommand && cmd.Flags().Lookup("environment").Changed {
		environmentValue, err := cmd.Flags().GetString("environment")
		if !helpers.IsValidString(environmentValue, validEnvironmentTypes) || err != nil {
			return isNamedSite, siteName,
				fmt.Errorf("the environment, %s, is not valid. You must use either `local`, `development`, `staging` or `production`", environmentValue)
		}
	}

	// Process the name flag if set
	if cmd.Flags().Lookup("name").Changed {
		isNamedSite = true

		// Check that we're not using invalid start flags for the start command
		if isStartCommand && cmd.Flags().Lookup("type").Changed {
			typeValue, _ := cmd.Flags().GetString("type")
			if typeValue != "site" {
				return isNamedSite, siteName,
					fmt.Errorf("the type, %s, is not valid. Only a `site` is valid with named sites", typeValue)
			}
		}

		siteName = helpers.SanitizeSiteName(cmd.Flags().Lookup("name").Value.String())
	}

	return isNamedSite, siteName, nil
}
