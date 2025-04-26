package settings

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
		name:         "isStartCommand",
		defaultValue: "false",
		settingType:  "bool",
	},
	{
		name:         "activate",
		defaultValue: "true",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			ShortName: "a",
			Usage:     "Activate the current plugin or theme (only works when used with the 'plugin' or 'theme' flags).",
		},
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
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			Usage: "Select the database server you wish to use with your installation.",
		},
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
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			Usage: "Sets the WP_ENVIRONMENT_TYPE for the site.",
		},
	},
	{
		name:         "mailpit",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			ShortName: "m",
			Usage:     "Enable Mailpit when starting the container.",
		},
	},
	{
		name:         "multisite",
		defaultValue: "none",
		settingType:  "string",
		validValues: []string{
			"none",
			"subdomain",
			"subdirectory"},
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			NoOptDefValue: "subdomain",
			Usage:         "Creates your new site as a multisite installation.",
		},
	},
	{
		name:         "php",
		defaultValue: "8.4",
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
		hasStartFlag: true,
		startFlag: StartFlag{
			Usage: "Installs and activates the specified plugins. Multiple plugins should be separated by commas",
		},
	},
	{
		name:         "removeDefaultPlugins",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			ShortName: "r",
			Usage:     "If true will remove the default plugins installed with WordPress (Akismet and Hello Dolly) when starting a site.",
		},
	},
	{
		name:         "scriptDebug",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			ShortName: "c",
			Usage:     "Enable SCRIPT_DEBUG when starting the WordPress site.",
		},
	},
	{
		name:         "ssl",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			ShortName: "s",
			Usage:     "Whether the site should default to SSL (https) or not.",
		},
	},
	{
		name:         "theme",
		defaultValue: "",
		settingType:  "string",
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			Usage: "Installs and activates a theme when starting a WordPress site.",
		},
	},
	{
		name:         "type",
		defaultValue: "site",
		settingType:  "string",
		validValues: []string{
			"site",
			"plugin",
			"theme"},
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			Usage: "Set the type of the installation, `site`, `plugin` or `theme`.",
		},
	},
	{
		name:         "typeDetected",
		defaultValue: "false",
		settingType:  "bool",
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
		hasStartFlag: true,
		startFlag: StartFlag{
			ShortName: "d",
			Usage:     "Enable WP_Debug when starting the WordPress site.",
		},
	},
	{
		name:         "xdebug",
		defaultValue: "false",
		settingType:  "bool",
		hasLocal:     true,
		hasGlobal:    true,
		hasStartFlag: true,
		startFlag: StartFlag{
			ShortName: "x",
			Usage:     "Enable Xdebug when starting the WordPress site.",
		},
	},
}

const (
	certOS                 = "darwin"
	configFolderName       = ".config/kana"
	defaultDirPermissions  = 0750
	defaultFilePermissions = 0644
	domain                 = "sites.kana.sh"
	mariadbVersion         = "11"
	mysqlVersion           = "9"
	rootCert               = "kana.root.pem"
	rootKey                = "kana.root.key"
	siteCert               = "kana.site.pem"
	siteKey                = "kana.site.key"
)
