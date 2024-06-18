package options

// The following are the default settings for Kana.

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
