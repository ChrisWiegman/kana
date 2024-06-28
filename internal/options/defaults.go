package options

var defaults = []Setting{
	{
		name:         "appDirectory",
		defaultValue: "",
		settingType:  "string",
	},
	{
		name:         "siteDirectory",
		defaultValue: "",
		settingType:  "string",
	},
	{
		name:         "workingDirectory",
		defaultValue: "",
		settingType:  "string",
	},
	{
		name:         "name",
		defaultValue: "",
		settingType:  "string",
	},
	{
		name:         "isNamed",
		defaultValue: "false",
		settingType:  "bool",
	},
	{
		name:         "isNew",
		defaultValue: "false",
		settingType:  "bool",
	},
	{
		name:         "activate",
		defaultValue: "true",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "adminEmail",
		defaultValue: "admin@sites.kana.sh",
		settingType:  "string",
		hasGlobal:    true,
	},
	{
		name:         "adminPassword",
		defaultValue: "password",
		settingType:  "string",
		hasGlobal:    true,
	},
	{
		name:         "adminUser",
		defaultValue: "admin",
		settingType:  "string",
		hasGlobal:    true,
	},
	{
		name:         "automaticLogin",
		defaultValue: "true",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "database",
		defaultValue: "mariadb",
		settingType:  "string",
		validValues: []string{
			"mariadb",
			"mysql",
			"sqlite"},
		hasLocal:  true,
		hasGlobal: true,
	},
	{
		name:         "databaseClient",
		defaultValue: "phpmyadmin",
		settingType:  "string",
		validValues: []string{
			"phpmyadmin",
			"tableplus"},
		hasLocal:  true,
		hasGlobal: true,
	},
	{
		name:         "databaseVersion",
		defaultValue: mariadbVersion,
		settingType:  "string",
		hasGlobal:    true,
		hasLocal:     true,
	},
	{
		name:         "environment",
		defaultValue: "local",
		settingType:  "string",
		validValues: []string{
			"local",
			"development",
			"staging",
			"production"},
		hasLocal:  true,
		hasGlobal: true,
	},
	{
		name:         "mailpit",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "multisite",
		defaultValue: "none",
		settingType:  "string",
		validValues: []string{
			"none",
			"subdomain",
			"subdirectory"},
		hasLocal:  true,
		hasGlobal: true,
	},
	{
		name:         "php",
		defaultValue: "8.2",
		settingType:  "string",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "plugins",
		defaultValue: "",
		settingType:  "slice",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "removeDefaultPlugins",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "scriptDebug",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "ssl",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "theme",
		defaultValue: "",
		settingType:  "string",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "type",
		defaultValue: "site",
		settingType:  "string",
		validValues: []string{
			"site",
			"plugin",
			"theme"},
		hasLocal:  true,
		hasGlobal: true,
	},
	{
		name:         "updateInterval",
		defaultValue: "7",
		settingType:  "int",
		hasGlobal:    true,
	},
	{
		name:         "wpdebug",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
	{
		name:         "xdebug",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
	},
}

const (
	configFolderName       = ".config/kana"
	defaultDirPermissions  = 0750
	defaultFilePermissions = 0644
	domain                 = "sites.kana.sh"
	mariadbVersion         = "11"
	mysqlVersion           = "8"
)
