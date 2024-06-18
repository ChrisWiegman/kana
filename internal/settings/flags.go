package settings

import (
	"strings"

	"github.com/spf13/cobra"
)

// ProcessStartFlags Process the start flags and save them to the settings object.
func ProcessStartFlags(cmd *cobra.Command, flags *StartFlags, options *Options) {
	if cmd.Flags().Lookup("xdebug").Changed {
		options.Xdebug = flags.Xdebug
	}

	if cmd.Flags().Lookup("wpdebug").Changed {
		options.WPDebug = flags.WPDebug
	}

	if cmd.Flags().Lookup("scriptdebug").Changed {
		options.ScriptDebug = flags.ScriptDebug
	}

	if cmd.Flags().Lookup("ssl").Changed {
		options.SSL = flags.SSL
	}

	if cmd.Flags().Lookup("mailpit").Changed {
		options.Mailpit = flags.Mailpit
	}

	if cmd.Flags().Lookup("type").Changed {
		options.Type = flags.Type
	}

	if cmd.Flags().Lookup("theme").Changed {
		options.Theme = flags.Theme
	}

	if cmd.Flags().Lookup("multisite").Changed {
		options.Multisite = flags.Multisite
	}

	if cmd.Flags().Lookup("activate").Changed {
		options.Activate = flags.Activate
	}

	if cmd.Flags().Lookup("environment").Changed {
		options.Environment = flags.Environment
	}

	if cmd.Flags().Lookup("plugins").Changed {
		options.Plugins = strings.Split(flags.Plugins, ",")
	}

	if cmd.Flags().Lookup("remove-default-plugins").Changed {
		options.RemoveDefaultPlugins = flags.RemoveDefaultPlugins
	}

	if cmd.Flags().Lookup("database").Changed {
		options.Database = flags.Database
	}
}
