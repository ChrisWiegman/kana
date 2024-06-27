package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ChrisWiegman/kana/internal/console"

	"github.com/aquasecurity/table"
)

// ListSettings Lists all settings for the config command.
func ListSettings(settings *Settings, consoleOutput *console.Console) {
	if consoleOutput.JSON {
		printJSONSettings(settings)
		return
	}

	settingsTable := table.New(os.Stdout)
	globalSettings := new(Options)
	localSettings := new(Options)

	updateSettingsFromViper(settings.global, globalSettings)
	updateSettingsFromViper(settings.local, localSettings)

	localPlugins := []string{}
	globalPlugins := []string{}

	for _, plugin := range localSettings.Plugins {
		localPlugins = append(localPlugins, consoleOutput.Bold(plugin))
	}

	for _, plugin := range globalSettings.Plugins {
		globalPlugins = append(globalPlugins, consoleOutput.Bold(plugin))
	}

	settingsTable.SetHeaders("Setting", "Global Value", "Local Value")

	settingsTable.AddRow("activate",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.Activate)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.Activate)))
	settingsTable.AddRow("admin.email", consoleOutput.Bold(globalSettings.AdminEmail), consoleOutput.Bold(localSettings.AdminEmail))
	settingsTable.AddRow("admin.password",
		consoleOutput.Bold(globalSettings.AdminPassword),
		consoleOutput.Bold(localSettings.AdminPassword))
	settingsTable.AddRow("admnin.username",
		consoleOutput.Bold(globalSettings.AdminUsername),
		consoleOutput.Bold(localSettings.AdminUsername))
	settingsTable.AddRow("automaticLogin",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.AutomaticLogin)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.AutomaticLogin)))
	settingsTable.AddRow("database",
		consoleOutput.Bold(globalSettings.Database),
		consoleOutput.Bold(localSettings.Database))
	settingsTable.AddRow("databaseClient",
		consoleOutput.Bold(globalSettings.DatabaseClient),
		consoleOutput.Bold(localSettings.DatabaseClient))
	settingsTable.AddRow("databaseVersion",
		consoleOutput.Bold(globalSettings.DatabaseVersion),
		consoleOutput.Bold(localSettings.DatabaseVersion))
	settingsTable.AddRow("environment", consoleOutput.Bold(globalSettings.Environment), consoleOutput.Bold(localSettings.Environment))
	settingsTable.AddRow("mailpit",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.Mailpit)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.Mailpit)))
	settingsTable.AddRow("multisite", consoleOutput.Bold(globalSettings.Multisite), consoleOutput.Bold(localSettings.Multisite))
	settingsTable.AddRow("php", consoleOutput.Bold(globalSettings.PHP), consoleOutput.Bold(localSettings.PHP))
	settingsTable.AddRow("plugins",
		consoleOutput.Bold(strings.Join(globalPlugins, "\n")),
		consoleOutput.Bold(strings.Join(localPlugins, "\n")))
	settingsTable.AddRow("removeDefaultPlugins",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.RemoveDefaultPlugins)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.RemoveDefaultPlugins)))
	settingsTable.AddRow("ssl",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.SSL)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.SSL)))
	settingsTable.AddRow("scriptDebug",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.ScriptDebug)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.ScriptDebug)))
	settingsTable.AddRow("theme",
		consoleOutput.Bold(globalSettings.Theme),
		consoleOutput.Bold(localSettings.Theme))
	settingsTable.AddRow("type", consoleOutput.Bold(globalSettings.Type), consoleOutput.Bold(localSettings.Type))
	settingsTable.AddRow("updateInterval", consoleOutput.Bold(strconv.FormatInt(globalSettings.UpdateInterval, 10)), "")
	settingsTable.AddRow("wpdebug",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.WPDebug)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.WPDebug)))
	settingsTable.AddRow("xdebug",
		consoleOutput.Bold(strconv.FormatBool(globalSettings.Xdebug)),
		consoleOutput.Bold(strconv.FormatBool(localSettings.Xdebug)))

	settingsTable.Render()
}

func (s *Settings) PrintSingleSetting(name string, consoleOutput *console.Console) {
	value, err := s.GetGlobalSetting(name)
	if err != nil {
		consoleOutput.Error(err)
	}
	if consoleOutput.JSON {
		type JSONSetting struct {
			Setting, Value string
		}

		setting := JSONSetting{
			Setting: name,
			Value:   value,
		}

		str, _ := json.Marshal(setting)

		fmt.Println(string(str))
	} else {
		consoleOutput.Println(value)
	}
}

// printJSONSettings Prints out all settings in JSON format.
func printJSONSettings(settings *Settings) {
	type JSONSettings struct {
		Global, Local map[string]interface{}
	}

	jsonSettings := JSONSettings{
		Global: settings.global.AllSettings(),
		Local:  settings.local.AllSettings(),
	}

	str, _ := json.Marshal(jsonSettings)

	fmt.Print(string(str))
}
