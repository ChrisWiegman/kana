# Kana

Kana is a CLI (command line) tool for developing WordPress sites, plugins and themes efficiently.

# Why Kana?

I've gone through many different tools to run WordPress locally over the years. All of them are either extremely complex or don't support basic features such as ensuring plugin and theme development can be a first-class experience. I rarely build sites with WordPress and I wanted a tool that will allow me to build the plugins and themes I do work on as efficiently as possible.

# System requirements

## MacOS

- [Docker Desktop](https://www.docker.com)

## Linux

- [Docker Engine](https://docs.docker.com/engine/install/)

Docker Desktop for Linux may work but I have not tested it.

# Installing Kana

There are a few options for installing Kana. You can use [Homebrew](https://brew.sh) (recommended), you can install it from [the "releases" page](https://github.com/ChrisWiegman/kana/releases) here or you can build it manually.

## Install from Homebrew

Installing from [Homebrew](https://brew.sh) is the recommended approach on both Mac and Linux as it allows for automatic updates when needed. To install from Homebrew run the following command:

```
brew install ChrisWiegman/kana/kana
```

Note that, as there are numerous ways to install Docker, I have chosen, at least for now, to not list it as a dependency when installing via Homebrew. You'll want to make sure Docker is already installed or install it with `brew install --cask docker` if you're on Mac (see [this documentation](https://docs.docker.com/engine/install/) if you're in Linux).

## Download from GitHub releases

Simply download the latest release from our [release page](https://github.com/ChrisWiegman/kana/releases) and extract the CLI to a location accessible by your system PATH

**Note for Mac users** I have not signed the download copy so you'll need to manually allow it in your Mac settings if you download it from the releases page. Install it via Homebrew to avoid this step.

## Build manually

You will need [Go](https://go.dev) installed locally to build the application.

1. Clone this repo `git clone https://github.com/ChrisWiegman/kana.git`
2. CD into the repo and run `make install`

Assuming you have Go properly setup with `GOBIN` in your system path, you should now be able to use Kana. Run `kana version` to test.

# Updating Kana

If you build Kana from source you'll need to manually update Kana with a `git pull` and then a fresh build.

If you use Homebrew or install Kana from the releases page Kana will automatically check for updates and warn you of a new version. You can then update via homebrew, if that's how you installed it, or by running `kana update` which will replace the existing binary with the updated version.

# Using Kana

At it's most basic you can start a zero-config Kana site by running `kana start` in your terminal. This will create a new Kana site based on your current directory and open it in your default browser.

Kana relies on [Traefik](https://traefik.io) to map real domains to local sites. You can run as many sites as you need and each will be mapped to a subdomain of _sites.kana.sh_.

## Start

`kana start` will start a kana site based on your current directory and open it in your browser. It will detect if the current directory is a plugin or a theme and start the site as the appropriate type.

To login to the new site use the following:

- _User Name_: **admin**
- _Password_: **password**

Note: these can be changed in the config. Please see below.

### Start options

`--type` Defaults to `site` for developing a WordPress site. Can set to `plugin` map the current directory as a plugin within the created site or `theme` to map the current directory as a theme within the created site.

`--xdebug` will start Xdebug on the site (see below for usage).

`--wpdebug` will enable `WP_DEBUG` on the site.

`--scriptdebug` will enable `SCRIPT_DEBUG` on the site.

`--environment` allows the user to change the `WP_ENVIRONMENT_TYPE` constant. Defaults to `local`. Valid options are `local`, `devlopment`, `staging` and `production`.

`--mailpit` will start an instance of [Mailpit](https://github.com/axllent/mailpit) to allow for email capture and troubleshooting.

`--ssl` will set the site's default URLs to use SSL.

`--name` The name flag allows you to run an arbitrary site from anywhere. For example, if you already started and stopped a site from a directory called _test_ you can run `kana start --name=test` to start that site from anywhere. If you use the `name` flag on a new site it will create that site without a link to any local folder. This can be handy for testing a plugin or other configuration but not that none of the other start flags will apply.

`--multisite` Use the multisite flag to setup a WordPress Multisite installation. The optional `subdomain` and `subdirectory` flags will allow for either type of installation.

`--removedefaultplugins` Will remove the default "Hello Dolly" and Akismet plugins when starting the site. Note this will not restore them if they've been manually removed.

`--theme` Sets the default theme if you do not wish to use the theme bundled with WordPress. Will attempt to download the theme from wordpress.org. Does not work if the site type is set to "theme"

`--plugins` A comma-separated list of plugins to install when starting the site.

`--database` By default Kana uses [MariaDB](https://mariadb.org) for its WordPress database. You can use MySQL or [SQLite](https://www.sqlite.org/index.html) instead by specifying `mysql` or `sqlite` as the database type here.

## Trusting the SSL certificate on Mac

On MacOS, Kana will automatically attempt to add its SSL certificate to the MacOS system Keychain the first time you start a site where SSL is the default. You can manually do this without starting a new site using the `kana trust-ssl` command.

## Importing an existing WordPress database

Kana offers a simple way to import an existing WordPress database. Just use the `kana db import <your database file>` to get started.

If you're coming from a site with a different home address you can specify `--replace-domain=<my old site domain>` to automatically replace it with the appropriate domain for your dev site.

### Example:

`kana db import --replace-domain=chriswiegman.com database.sql` would import the file _database.sql_ from my current directory and rename the old site address, chriswiegman.com, to the current and correct site address to work in Kana.

### Import options

`--replace-domain` The domain of your source site to replace with the appropriate Kana domain
`--preserve` Prevents Kana from dropping any existing database and overwrites what you have. Warning: this may result in unpredictable issues.

### Exporting your Kana database

You can also export the database file your Kana site is using with `kana db export`. By default it will save the file in your default site directory but you can specify a relative path to the file where you would like to export your database if you wish.

> *Note* Currently importang and exporting databases only works with MariaDB databases. [I am working on bringing this functionality to MySQL](https://github.com/docker-library/wordpress/pull/902) and hope to have it available with MySQL soon. I do not anticipate bringing this to SQLite for a while.

## Stop

`kana stop` will stop the current site and, if no other sites are running, will shut down shared containers like Traefik as well.

## List

`kana list` will list all sites known by Kana and their current running status. Any site listed can then be addressed with the `name` flag in other commands.

## Destroy

`kana destroy` will stop and destroy the current site. This is different than `stop` in that `stop` will leave the database and files it creates alone so you can start it again later. Once destroyed a site is irrecoverable.

By default Kana will prompt you to confirm any site you wish to destroy. You can bypass the prompt by adding the `--force` flag to the destroy command.

## Open

`kana open` will open the site in your default browser

`kana open -a` will open the WordPress Dashboard. This will also login the "admin" user unless the `automaticLogin` setting is set to false.

`kana open -t` will open the Traefik dashboard.

By default Kana will open the appropriate WordPress site. To open the database or Mailpit simply append the appropriate flag to the open command ie `kana open --database`.

Note that by default Kana will open the database in [phpMyAdmin](https://www.phpmyadmin.net). You can also tell Kana to open the database in [TablePlus](https://tableplus.com) instead by setting the `databaseClient` configuration setting to `tableplus`.

Currently pphpMyAdmin and TablePlus are the only two clients I've configured. If you would like to use a different client, please [open an issue](https://github.com/ChrisWiegman/kana/issues) and I'd be happy to take a look.

> *Note* Opening the Database directly with Kana doesn't work for SQLite databases. To open a SQLite database directly navigate to `<your-site-folder>/wp-content/database/.ht.sqlite` and open the file directly.

## wp-cli

`kana wp <WP-CLI COMMAND>` will execute a [wp-cli](https://wp-cli.org) command on your site. For example `kana wp plugin list` will list all the plugins on the site and their associated statuses

# Configuring Kana

The above commands will get an individual site up and running but there are a few more options to consider that can be changed for a given site or globally

## Global Config

Kana has a handful of options that apply to all new sites created with the app. You can adjust these with the `config` command as noted below:

`kana config` will list all changeable defaults for a new site. Currently these include the following:

- `activate` **true** - if the project site is set to `theme` or `plugin` this will activate the project on first load
- `adminEmail` __admin@kanasite.localhost__ - the admin email address for the default admin account
- `adminPassword` **password** - the default password used to login to WordPress
- `adminUser` **admin** - the default username used to login to WordPress
- `automaticLogin` **true** - will automatically login the "admin" user when accessing the WordPress dashboard
- `database` **mariadb** - Specify the database server for WordPress, currently either `mariadb`, `mysql` or `sqlite`
- `databaseClient` **phpmyadmin** - the default database client for accessing the database directly (currently `phpmyadmin` and `tableplus` are supported)
- `databaseVersion` **11** - the default database version used for sites. 11 is chosen for the default MariaDB database. You will need to update this if you switch to MySQL.
- `environment` **local** - the default usage of the `environment` start flag
- `mailpit` **false** - the default usage of the `mailpit` start flag
- `multisite` **none** - set to either `subdirectory` or `subdomain` to create the site as the appropriate type of Multisite installation.
- `php` **8.2** - the default PHP version used for new sites (see [https://hub.docker.com/_/wordpress] for all supported versions)
- `removeDefaultPlugins` **false** - removes the default "Hello Dolly" and Akismet plugins when starting a new site. Note this will not restore them if they've already been removed.
- `scriptDebug` **false** - the default usage of the `scriptDebug` wp-config item
- `ssl` **false** - the default usage of the `ssl` start flag
- `theme` ***<empty string>*** - the default theme to be installed from wordpress.org and activated with new sites
- `type` **site** - the type of the Kana site you're starting. Current options are "site" "plugin" and "theme"
- `updateInterval` **1** - the number of days Kana will wait between checking for updated Docker images and other updates. Set this to `0` to disable the check for newer images altogether (Kana will only download missing images)
- `wpdebug` **false** - the default usage of the `wpdebug` start flag
- `xdebug` **false** - the default usage of the `xdebug` start flag

You can get or set any of the above options using a similar syntax to GIT's config. For example:

`kana config adminEmail` will print the value of the admin.email setting
`kana config adminEmail myemail@somedomain.com` will change the value of the admin.email setting to "myemail@somedomain.com".

The above syntax will allow you to change the defaults for any of the options listed

## Site Config

In addition to the global config, certain items above can be overridden for any given site. For a site without a `name` flag (as seen in the start command), simply create a _.kana.json_ file in the current directory. You can populate it with the following options:

- `activate` **true** - if the project site is set to `theme` or `plugin` this will activate the project on first load
- `adminEmail` __admin@kanasite.localhost__ - the admin email address for the default admin account
- `adminPassword` **password** - the default password used to login to WordPress
- `adminUser` **admin** - the default username used to login to WordPress
- `automaticLogin` **true** - will automatically login the "admin" user when accessing the WordPress dashboard
- `database` **mariadb** - Specify the database server for WordPress, currently either `mariadb`, `mysql` or `sqlite`
- `databaseClient` **phpmyadmin** - the default database client for accessing the database directly (currently `phpmyadmin` and `tableplus` are supported)
- `databaseVersion` **11** - the default database version used for sites. 11 is chosen for the default MariaDB database. You will need to update this if you switch to MySQL.
- `environment` **local** - the default usage of the `environment` start flag
- `mailpit` **false** - the default usage of the `mailpit` start flag
- `multisite` **none** - set to either `subdirectory` or `subdomain` to create the site as the appropriate type of Multisite installation.
- `php` **8.2** - the default PHP version used for new sites (see [https://hub.docker.com/_/wordpress] for all supported versions)
- `plugins` **[]** - an array of plugins to install and activate when starting the new site. These are slugs from the Plugins section of WordPress.org.
- `removeDefaultPlugins` **false** - removes the default "Hello Dolly" and Akismet plugins when starting a new site. Note this will not restore them if they've already been removed.
- `scriptDebug` **false** - the default usage of the `scriptDebug` start flag
- `ssl` **false** - the default usage of the `ssl` start flag
- `theme` ***<empty string>*** - the default theme to be installed from wordpress.org and activated with the site
- `type` **site** - the type of the Kana site you're starting. Current options are "site" "plugin" and "theme"
- `wpdebug` **false** - the default usage of the `wpdebug` start flag
- `xdebug` **false** - the default usage of the `xdebug` start flag

### Export a sites Kana config automatically

`kana export` will create a _.kana.json_ configuration file in your current folder exporting the configuration of the current site including PHP version, active plugins and associated options as shown above

# Accessing the database directly

Currently there are two methods to access the database directly. First you can access the database via phpMyAdmin or TablePlus by running `kana open --database` for the site in question.

Note that by default Kana will open the database in [phpMyAdmin](https://www.phpmyadmin.net). You can also tell Kana to open the database in [TablePlus](https://tableplus.com) instead by setting the `databaseClient` configuration setting to `tableplus`.

You can also access the database directly by viewing the database port with `docker ps` and using the database port and the following configuration in the app of your choice:

- **Database host**: _kana\_`your site name`\_database_
- **Database name**: _wordpress_
- **Database user**: _wordpress_
- **Database password**: _wordpress_

# Verify you're using Kana

You can verify you're in a Kana environment by verifying the `IS_KANA_ENVIRONMENT` env variable. Here's an example:

```php
if ( getenv( 'IS_KANA_ENVIRONMENT' ) === true ) {
    // We're on a Kana environment
}
```

# Using Xdebug

You can setup [Xdebug](https://xdebug.org) when starting a site in Kana using the `--xdebug` flag with the start command or by setting the `xdebug` setting globally or at the site level.

Once a site is running you can see if Xdebug is running by using the `kana xdebug` command which will return _on_ if Xdebug is running or _off_ if it is not.

To start or stop Xdebug on a running site use `xdebug on` or `xdebug off` as appropriate. The output of this command will be either _on_ or _off_ to indicate the status of Xdebug when the command is complete.

Currently Kana only supports step debugging in xdebug. To use this with VSCode create a _.vscode/launch.json_ file with the following:

```{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Listen for XDebug",
            "type": "php",
            "request": "launch",
            "port": 9003,
            "log": true,
            "pathMappings": {
                "/var/www/html/wp-content/plugins/<MY KANA FOLDER NAME>/": "${workspaceFolder}",
                "/var/www/html/": "${workspaceFolder}/wordpress",
            }
        }
    ]
}
```

To trigger step debugging you'll also need the appropriate extension for your browser:

- [Xdebug Helper for Firefox](https://addons.mozilla.org/en-GB/firefox/addon/xdebug-helper-for-firefox/) ([source](https://github.com/BrianGilbert/xdebug-helper-for-firefox)).
- [Xdebug Helper for Chrome](https://chrome.google.com/extensions/detail/eadndfjplgieldjbigjakmdgkmoaaaoc) ([source](https://github.com/mac-cain13/xdebug-helper-for-chrome)).
- [XDebugToggle for Safari](https://apps.apple.com/app/safari-xdebug-toggle/id1437227804?mt=12) ([source](https://github.com/kampfq/SafariXDebugToggle)).

# Flushing cache and transients

Two wp-cli commands I find myself using regularly when working on WordPress are `wp transient delete --all` and `wp cache flush`. I use them so often that it seemed like a good idea to make them easier to access with Kana. As a result I've added the `kana flush` command which will call both on the specified site.

# Viewing the Kana changelog

It's always good to know what's changed before updating. You can use `kana changelog` to take to Kana's releases on GitHub where you can view the current changes and look for anything you might want to wait on.

# This project is under active development

Note that I am using this project for my own work and it is under active development. Some of the things I'm currently working on include:

- Better site management commands
- Much more clear prompts and messages on the commands themselves
- Writing a lot more tests (it's a personal project, I start where I can)
- A proper website for all this documentation (I already bought a domain, after all)
- Bugfixes and other tweaks as I find them necessary for my own use
- Better integration with tools like VSCode and others

# Completely Uninstalling Kana

I hate apps that leave leftovers on your machine. When stopping a site all Docker resources except the images will be removed. To remove the app completely beyond that you'll want to delete the following:

1. Delete the application from your $GOBIN or system path (or run `brew uninstall kana` if installed via homebrew)
2. Delete the `~/.config/kana` folder which contains all site and app configuration
3. (Mac only) Delete the `Kana Development CA` certificate from the _System_ keychain in the _Keychain Access_ app
4. If installed via homebrew run `brew untap ChrisWiegman/kana` to remove the Homebrew tap

You can also safely remove any new images added however it is not a requirement. Many other apps might share those images leading to your system simply needing to download them again.

# Using Kana in other projects

While Kana cannot easily be used as a package itself, you can import the binary itself into your project. If you do so, consider using the `output-json` flag on all commands. This will convert all output to JSON format to make consumption easier when the Kana application is embedded elsewhere.

Why do this? This will make it easier for me to work with Kana in a small toolbar app I'm building as well as with a [Visual Studio Code](https://code.visualstudio.com/) extension I have planned which will allow me to see what is going on with Kana and control it beyond the terminal.
