package settings

// The following are the default settings for Kana.
var (
	rootKey              = "kana.root.key"
	rootCert             = "kana.root.pem"
	siteCert             = "kana.site.pem"
	siteKey              = "kana.site.key"
	domain               = "sites.kana.li"
	configFolderName     = ".config/kana"
	php                  = "8.1"
	mariadb              = "11"
	siteType             = "site"
	xdebug               = false
	mailpit              = false
	local                = false
	ssl                  = false
	wpdebug              = false
	removeDefaultPlugins = false
	adminUsername        = "admin"
	adminPassword        = "password"
	adminEmail           = "admin@sites.kana.li"
	imageUpdateDays      = 7
	activate             = true
	databaseClient       = "phpmyadmin"
	multisite            = "none"
)

// Default permissions for all new files and folders.
const (
	defaultDirPermissions  = 0750
	defaultFilePermissions = 0644
)

var validPHPVersions = []string{
	"7.4",
	"8.0",
	"8.1",
	"8.2",
	"8.4",
}

var validMariaDBVersions = []string{
	"10",
	"11",
}

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
