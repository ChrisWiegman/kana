package settings

import (
	"encoding/json"
	"fmt"
	"os"
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

	settingsTable.SetHeaders("Setting", "Global Value", "Local Value")

	globalSettings := settings.GetAll("global")
	localSettings := settings.GetAll("local")

	for i := range settings.settings {
		if !settings.settings[i].hasGlobal && !settings.settings[i].hasLocal {
			continue
		}

		globalOutput := fmt.Sprint(globalSettings[settings.settings[i].name])
		localOutput := fmt.Sprint(localSettings[settings.settings[i].name])

		if settings.settings[i].settingType == "slice" { //nolint:goconst
			globalOutput = strings.Join(globalSettings[settings.settings[i].name].([]string), "\n")
			localOutput = strings.Join(localSettings[settings.settings[i].name].([]string), "\n")
		}

		if !settings.settings[i].hasLocal {
			localOutput = ""
		}

		settingsTable.AddRow(settings.settings[i].name,
			consoleOutput.Bold(globalOutput),
			consoleOutput.Bold(localOutput))
	}

	settingsTable.Render()
}

func PrintSingleSetting(name string, kanaSettings *Settings, consoleOutput *console.Console) {
	globalSettings := kanaSettings.GetAll("global")

	if consoleOutput.JSON {
		type JSONSetting struct {
			Setting, Value string
		}

		setting := JSONSetting{
			Setting: name,
			Value:   fmt.Sprint(globalSettings[name]),
		}

		str, _ := json.Marshal(setting)

		fmt.Println(string(str))
	} else {
		consoleOutput.Println(fmt.Sprint(globalSettings[name]))
	}
}

// printJSONSettings Prints out all settings in JSON format.
func printJSONSettings(settings *Settings) {
	type JSONSettings struct {
		Global, Local map[string]interface{}
	}

	globalSettings := settings.GetAll("global")
	localSettings := settings.GetAll("local")

	jsonSettings := JSONSettings{
		Global: globalSettings,
		Local:  localSettings,
	}

	str, _ := json.Marshal(jsonSettings)

	fmt.Print(string(str))
}
