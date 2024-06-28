package settings

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func AddStartFlags(cmd *cobra.Command, settings *Settings) {
	if cmd.Use == "start" { //nolint:goconst
		for i := range defaults {
			if defaults[i].hasStartFlag {
				switch defaults[i].settingType {
				case "bool": //nolint:goconst
					boolValue, _ := strconv.ParseBool(defaults[i].defaultValue)

					if defaults[i].startFlag.ShortName != "" {
						cmd.Flags().BoolP(defaults[i].name, defaults[i].startFlag.ShortName, boolValue, defaults[i].startFlag.Usage)
					} else {
						cmd.Flags().Bool(defaults[i].name, boolValue, defaults[i].startFlag.Usage)
					}
				case "slice":
					sliceValue := strings.Split(defaults[i].defaultValue, ",")

					cmd.Flags().StringSlice(defaults[i].name, sliceValue, defaults[i].startFlag.Usage)
				default:
					cmd.Flags().String(defaults[i].name, defaults[i].defaultValue, defaults[i].startFlag.Usage)
				}

				if defaults[i].startFlag.NoOptDefValue != "" {
					cmd.Flags().Lookup(defaults[i].name).NoOptDefVal = defaults[i].startFlag.NoOptDefValue
				}
			}
		}
	}
}

// processStartFlags Process the start flags and save them to the settings object.
func processStartFlags(cmd *cobra.Command, settings *Settings) error {
	if cmd.Use == "start" {
		for i := range settings.settings {
			if settings.settings[i].hasStartFlag && cmd.Flags().Lookup(settings.settings[i].name).Changed {
				err := validateFlags(cmd, settings)
				if err != nil {
					return err
				}

				if settings.settings[i].settingType == "slice" {
					strings.Split(cmd.Flags().Lookup("plugins").Value.String(), ",")
				} else {
					err = settings.Set(settings.settings[i].name, cmd.Flags().Lookup(settings.settings[i].name).Value.String())
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func validateFlags(cmd *cobra.Command, settings *Settings) error {
	if cmd.Flags().Lookup("type").Changed &&
		settings.GetBool("isNamed") {
		return fmt.Errorf("the type flag is not valid with named sties. Only a `site` is valid with named sites")
	}

	return nil
}
