package settings

import (
	"os"

	"github.com/spf13/viper"
)

// Settings Individual Settings for use throughout the app lifecycle.
type Settings struct {
	global *viper.Viper // The "global" viper holding app-wide settings.
	local  *viper.Viper // The "local" viper holding site-specific settings.

	Activate             bool   // Activate the plugin or theme for the appropriate site type.
	AdminEmail           string // The email address for the admin user.
	AdminPassword        string // The password for the admin user.
	AdminUsername        string // The username for the admin user.
	AppDomain            string // The top-level domain for the app.
	AppDirectory         string
	AutomaticLogin       bool
	DatabaseClient       string
	Environment          string
	ImageUpdateDays      int
	IsNamedSite          bool
	IsNewSite            bool
	Local                bool
	Mailpit              bool
	MariaDB              string
	Multisite            string
	Name                 string
	PHP                  string
	Plugins              []string
	Protocol             string
	RemoveDefaultPlugins bool
	RootCert             string
	RootKey              string
	ScriptDebug          bool
	SiteCert             string
	SiteDirectory        string
	SiteDomain           string
	SiteKey              string
	SSL                  bool
	Theme                string
	Type                 string
	URL                  string
	Version              string
	WorkingDirectory     string
	WPDebug              bool
	Xdebug               bool
}

type StartFlags struct {
	Xdebug               bool
	WPDebug              bool
	Mailpit              bool
	SSL                  bool
	Activate             bool
	RemoveDefaultPlugins bool
	ScriptDebug          bool
	Environment          string
	Multisite            string
	Type                 string
	Theme                string
}

type LocalSettings struct {
	Activate             bool
	AutomaticLogin       bool
	DatabaseClient       string
	Environment          string
	Mailpit              bool
	Multisite            string
	Plugins              []string
	RemoveDefaultPlugins bool
	ScriptDebug          bool
	SSL                  bool
	Theme                string
	Type                 string
	WPDebug              bool
	Xdebug               bool
}

type File struct {
	LocalPath   string
	Name        string
	Permissions os.FileMode
	Template    string
}

type KanaPluginVars struct {
	SiteName string
	Version  string
}
