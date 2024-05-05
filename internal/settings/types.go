package settings

import (
	"os"

	"github.com/spf13/viper"
)

// Settings Individual Settings for use throughout the app lifecycle.
type Settings struct {
	Local, Xdebug, Mailpit, SSL, WPDebug, Activate, ScriptDebug bool
	RemoveDefaultPlugins                                        bool
	IsNewSite, IsNamedSite                                      bool
	AdminLogin                                                  bool
	ImageUpdateDays                                             int
	AdminEmail, AdminPassword, AdminUsername                    string
	AppDirectory, SiteDirectory, WorkingDirectory               string
	AppDomain, SiteDomain                                       string
	Name                                                        string
	PHP                                                         string
	MariaDB                                                     string
	RootCert, RootKey, SiteCert, SiteKey                        string
	URL, Protocol                                               string
	Type                                                        string
	DatabaseClient                                              string
	Multisite                                                   string
	Environment                                                 string
	Version                                                     string
	Theme                                                       string
	Plugins                                                     []string
	global                                                      *viper.Viper
	local                                                       *viper.Viper
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
	Mailpit, Xdebug, SSL, WPDebug, Activate, ScriptDebug bool
	AdminEmail, AdminPassword, AdminUsername             string
	RemoveDefaultPlugins                                 bool
	AdminLogin                                           bool
	Type, DatabaseClient, Multisite, Environment         string
	Plugins                                              []string
	Theme                                                string
}

type File struct {
	Name, Template, LocalPath string
	Permissions               os.FileMode
}

type KanaPluginVars struct {
	SiteName, Version string
}
