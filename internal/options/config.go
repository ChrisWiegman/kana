package options

import (
	"encoding/json"
	"fmt"
	"os"

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

	globalSettings := settings.getAll("global")
	localSettings := settings.getAll("local")

	for _, setting := range settings.settings {
		if !setting.hasGlobal && !setting.hasLocal {
			continue
		}

		globalOutput := fmt.Sprint(globalSettings[setting.name])
		localOutput := fmt.Sprint(localSettings[setting.name])

		if !setting.hasGlobal {
			localOutput = ""
		}

		settingsTable.AddRow(setting.name,
			consoleOutput.Bold(globalOutput),
			consoleOutput.Bold(localOutput))
	}

	settingsTable.Render()
}

// func (s *Settings) PrintSingleSetting(name string, consoleOutput *console.Console) {
// 	value, err := s.GetGlobalSetting(name)
// 	if err != nil {
// 		consoleOutput.Error(err)
// 	}
// 	if consoleOutput.JSON {
// 		type JSONSetting struct {
// 			Setting, Value string
// 		}

// 		setting := JSONSetting{
// 			Setting: name,
// 			Value:   value,
// 		}

// 		str, _ := json.Marshal(setting)

// 		fmt.Println(string(str))
// 	} else {
// 		consoleOutput.Println(value)
// 	}
// }

// printJSONSettings Prints out all settings in JSON format.
func printJSONSettings(settings *Settings) {
	type JSONSettings struct {
		Global, Local map[string]interface{}
	}

	globalSettings := settings.getAll("global")
	localSettings := settings.getAll("local")

	jsonSettings := JSONSettings{
		Global: globalSettings,
		Local:  localSettings,
	}

	str, _ := json.Marshal(jsonSettings)

	fmt.Print(string(str))
}
