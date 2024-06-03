package settings

// The following are the default settings for Kana.
var (
	activate             = true
	adminEmail           = "admin@sites.kana.sh"
	adminPassword        = "password"
	adminUsername        = "admin"
	automaticLogin       = true
	configFolderName     = ".config/kana"
	database             = "mariadb"
	databaseClient       = "phpmyadmin"
	domain               = "sites.kana.sh"
	environment          = "local"
	imageUpdateDays      = 7
	mailpit              = false
	mariadb              = "11"
	multisite            = "none"
	php                  = "8.2"
	plugins              = []string{}
	removeDefaultPlugins = false
	rootCert             = "kana.root.pem"
	rootKey              = "kana.root.key"
	scriptDebug          = false
	siteCert             = "kana.site.pem"
	siteKey              = "kana.site.key"
	siteType             = "site"
	ssl                  = false
	theme                = ""
	wpdebug              = false
	xdebug               = false
)

// Default permissions for all new files and folders.
const (
	defaultDirPermissions  = 0750
	defaultFilePermissions = 0644
)

var validTypes = []string{
	"site",
	"plugin",
	"theme",
}

var validDatabaseClients = []string{
	"phpmyadmin",
	"tableplus",
}

var validDatabases = []string{
	"mariadb",
	"sqlite",
}

var validMultisiteTypes = []string{
	"none",
	"subdomain",
	"subdirectory",
}

var validEnvironmentTypes = []string{
	"local",
	"development",
	"staging",
	"production",
}
