package settings

// The following are the default settings for Kana.

var defaults = []Option{
	{
		Name:    "activate",
		Type:    "bool",
		Default: "true",
	},
	{
		Name:    "admin.email",
		Type:    "string",
		Default: "admin@sites.kana.sh",
	},
	{
		Name:    "admin.password",
		Type:    "string",
		Default: "password",
	},
	{
		Name:    "admin.username",
		Type:    "string",
		Default: "admin",
	},
	{
		Name:    "automaticlogin",
		Type:    "bool",
		Default: "true",
	},
	{
		Name:    "database",
		Type:    "string",
		Default: "mariadb",
	},
	{
		Name:    "databaseclient",
		Type:    "string",
		Default: "phpmyadmin",
	},
	{
		Name:    "databaseversion",
		Type:    "string",
		Default: mariadbVersion,
	},
	{
		Name:    "environment",
		Type:    "string",
		Default: "local",
	},
	{
		Name:    "updateinterval",
		Type:    "int",
		Default: "7",
	},
	{
		Name:    "mailpit",
		Type:    "bool",
		Default: "false",
	},
	{
		Name:    "multisite",
		Type:    "string",
		Default: "none",
	},
	{
		Name:    "php",
		Type:    "string",
		Default: "8.2",
	},
	{
		Name:    "plugins",
		Type:    "array",
		Default: "",
	},
	{
		Name:    "removedefaultplugins",
		Type:    "bool",
		Default: "false",
	},
	{
		Name:    "scriptdebug",
		Type:    "bool",
		Default: "false",
	},
	{
		Name:    "ssl",
		Type:    "bool",
		Default: "false",
	},
	{
		Name:    "theme",
		Type:    "string",
		Default: "",
	},
	{
		Name:    "type",
		Type:    "string",
		Default: "site",
	},
	{
		Name:    "wpdebug",
		Type:    "bool",
		Default: "false",
	},
	{
		Name:    "xdebug",
		Type:    "bool",
		Default: "false",
	},
}

var defaultOptions = Options{
	Activate:             true,
	AdminEmail:           "admin@sites.kana.sh",
	AdminPassword:        "password",
	AdminUsername:        "admin",
	AutomaticLogin:       true,
	Database:             "mariadb",
	DatabaseClient:       "phpmyadmin",
	DatabaseVersion:      mariadbVersion,
	Environment:          "local",
	UpdateInterval:       7,
	Mailpit:              false,
	Multisite:            "none",
	PHP:                  "8.2",
	Plugins:              []string{},
	RemoveDefaultPlugins: false,
	ScriptDebug:          false,
	SSL:                  false,
	Theme:                "",
	Type:                 "site",
	WPDebug:              false,
	Xdebug:               false,
}

var appConstants = Constants{
	Domain: "sites.kana.sh",
	RootCert: Certificate{
		Certificate: "kana.root.pem",
		Key:         "kana.root.key",
	},
	SiteCert: Certificate{
		Certificate: "kana.site.pem",
		Key:         "kana.site.key",
	},
}

const (
	configFolderName = ".config/kana"
	mariadbVersion   = "11"
	mysqlVersion     = "8"
)

// Default permissions for all new files and folders.
const (
	defaultDirPermissions  = 0750
	defaultFilePermissions = 0644
)
