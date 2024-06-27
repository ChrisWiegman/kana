package settings

import (
	"strings"

	"github.com/spf13/cobra"
)

// processStartFlags Process the start flags and save them to the settings object.
func processStartFlags(cmd *cobra.Command, flags *StartFlags, settings *Settings) {
	if cmd.Flags().Lookup("activate").Changed {
		settings.settings.Activate = flags.Activate
	}

	if cmd.Flags().Lookup("database").Changed {
		settings.settings.Database = flags.Database
	}

	if cmd.Flags().Lookup("environment").Changed {
		settings.settings.Environment = flags.Environment
	}

	if cmd.Flags().Lookup("mailpit").Changed {
		settings.settings.Mailpit = flags.Mailpit
	}

	if cmd.Flags().Lookup("multisite").Changed {
		settings.settings.Multisite = flags.Multisite
	}

	if cmd.Flags().Lookup("plugins").Changed {
		settings.settings.Plugins = strings.Split(flags.Plugins, ",")
	}

	if cmd.Flags().Lookup("remove-default-plugins").Changed {
		settings.settings.RemoveDefaultPlugins = flags.RemoveDefaultPlugins
	}

	if cmd.Flags().Lookup("scriptdebug").Changed {
		settings.settings.ScriptDebug = flags.ScriptDebug
	}

	if cmd.Flags().Lookup("ssl").Changed {
		settings.settings.SSL = flags.SSL
	}

	if cmd.Flags().Lookup("theme").Changed {
		settings.settings.Theme = flags.Theme
	}

	if cmd.Flags().Lookup("type").Changed {
		settings.settings.Type = flags.Type
	}

	if cmd.Flags().Lookup("wpdebug").Changed {
		settings.settings.WPDebug = flags.WPDebug
	}

	if cmd.Flags().Lookup("xdebug").Changed {
		settings.settings.Xdebug = flags.Xdebug
	}
}
