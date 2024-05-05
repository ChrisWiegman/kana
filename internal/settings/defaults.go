package settings

// The following are the default settings for Kana.
var (
	activate             = true
	adminEmail           = "admin@sites.kana.li"
	adminPassword        = "password"
	adminUsername        = "admin"
	configFolderName     = ".config/kana"
	databaseClient       = "phpmyadmin"
	domain               = "sites.kana.li"
	environment          = "local"
	imageUpdateDays      = 7
	mailpit              = false
	mariadb              = "11"
	multisite            = "none"
	php                  = "8.2"
	removeDefaultPlugins = false
	rootCert             = "kana.root.pem"
	rootKey              = "kana.root.key"
	scriptDebug          = false
	siteCert             = "kana.site.pem"
	siteKey              = "kana.site.key"
	siteType             = "site"
	ssl                  = false
	wpdebug              = false
	xdebug               = false
	adminLogin           = true
	theme                = ""
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
