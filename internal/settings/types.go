package settings

import (
	"os"

	"github.com/spf13/viper"
)

// Settings is a complete collection of all settings values used by Kana.
type Settings struct {
	settings    Options
	constants   Constants
	directories Directories
	global      *viper.Viper
	local       *viper.Viper
	site        Site
	SettingsAPI
}

type SettingsAPI interface {
	Get(name string) interface{}
	GetArray(name string) []string
	GetBool(name string) bool
	GetGlobalSetting(name string) (string, error)
	GetInt(name string) int64
	Set(name, value string) error
}

// Certificate represents the certificate and key for a SSL certificate used by Kana.
type Certificate struct {
	Certificate, Key string
}

// Constants are static values that are used throughout the app.
type Constants struct {
	Domain   string // The top-level domain for the app.
	RootCert Certificate
	SiteCert Certificate
	Version  string
}

// Directories represents the directories that Kana uses to work with a site.
type Directories struct {
	App     string
	Site    string
	Working string
}

// File represents template information for files that need to be placed by Kana.
type File struct {
	LocalPath   string
	Name        string
	Permissions os.FileMode
	Template    string
}

type Option struct {
	Name    string
	Value   string
	Type    string
	Default string
}

// Options represents the options that can be configured by a user either globally or through a site's .kana.json.
type Options struct {
	Activate             bool   // Activate the plugin or theme for the appropriate site type.
	AdminEmail           string // The email address for the admin user.
	AdminPassword        string // The password for the admin user.
	AdminUsername        string // The username for the admin user.
	AutomaticLogin       bool
	Database             string
	DatabaseClient       string
	DatabaseVersion      string
	Environment          string
	Mailpit              bool
	Multisite            string
	PHP                  string
	Plugins              []string
	RemoveDefaultPlugins bool
	ScriptDebug          bool
	SSL                  bool
	Theme                string
	Type                 string
	UpdateInterval       int64
	WPDebug              bool
	Xdebug               bool
}

// PluginVersion represents the name and version of a plugin to allow for better templating.
type PluginVersion struct {
	SiteName string
	Version  string
}

// Site represents values specific to running a site within Kana.
type Site struct {
	Name             string
	IsNamed          bool
	IsNew            bool
	TypeIsDetected   bool
	HasLocalSettings bool
}

// StartFlags represents the flags that can be passed to the start command.
type StartFlags struct {
	Xdebug               bool
	WPDebug              bool
	Mailpit              bool
	SSL                  bool
	Activate             bool
	Database             string
	Plugins              string
	RemoveDefaultPlugins bool
	ScriptDebug          bool
	Environment          string
	Multisite            string
	Type                 string
	Theme                string
}
