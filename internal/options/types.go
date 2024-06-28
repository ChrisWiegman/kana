package options

import (
	"os"

	"github.com/knadh/koanf/v2"
)

// File represents template information for files that need to be placed by Kana.
type File struct {
	LocalPath   string
	Name        string
	Permissions os.FileMode
	Template    string
}

// PluginVersion represents the name and version of a plugin to allow for better templating.
type PluginVersion struct {
	SiteName string
	Version  string
}

// A collection of all settings values used by Kana.
type Settings struct {
	settings []Setting
	global   *koanf.Koanf
	local    *koanf.Koanf
}

// An individual setting and its associated data.
type Setting struct {
	defaultValue string
	name         string
	settingType  string
	currentValue string
	hasLocal     bool
	hasGlobal    bool
	hasStartFlag bool
	startFlag    StartFlag
	validValues  []string
}

// StartFlag represents the data needed to programmatically create a start flag.
type StartFlag struct {
	ShortName     string
	Usage         string
	NoOptDefValue string
}
