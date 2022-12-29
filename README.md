# Kana

Kana is a simple CLI tool for developing WordPress plugins and themes efficiently.

# Why Kana?

I've gone through many different tools to run WordPress sites over the years. All of them are either extremely complex or don't support basic features such as ensuring plugin and theme development can be a first-class experience. I rarely build sites with WordPress and I wanted a tool that will allow me to build the plugins and themes I do work on as efficiently as possible.

# System requirements

## MacOS

- [Docker Desktop](https://www.docker.com)

Note: Docker Desktop should be either version 4.14.0 or later or below 4.13.0 due to a change in Docker 4.13.0 that broke many tools by trying to move the docker.sock file. See [the 4.14.0 changelog](https://docs.docker.com/desktop/release-notes/#4140) for more details.

## Linux

- [Docker Engine](https://docs.docker.com/engine/install/)

# Installing Kana

There are a few options for installing Kana. You can use [Homebrew](https://brew.sh) (recommended), you can install it from the "releases" page here or you can build it manually.

## Install from Homebrew

Installing from [Homebrew](https://brew.sh) is the recommended approach on both Max and Linux as it allows for automatic updates when needed. To install from Homebrew run the following command:

```
brew install ChrisWiegman/kana/kana
```

Note that, as there are numerous ways to install Docker, I have chosen, at least for now, to not list it as a dependency when installing via Homebrew. You'll want to make sure Docker is already installed or install it with `brew install --cask docker` if you're on Mac (see [this documentation](https://docs.docker.com/engine/install/) if you're in Linux).

## Download from GitHub releases

Simply download the latest release from our [release page](https://github.com/ChrisWiegman/kana-cli/releases) and extract the CLI to a location accessible by your system PATH

**Note for Mac users** I have not signed the download copy so you'll need to manually allow it in your Mac settings if you download it from the releases page. Install it via Homebrew to avoid this step.

## Build manually

You will need [Go](https://go.dev) installed locally to build the application for now. I hope to fix this in the new future.

1. Clone this repo `git clone git@github.com:ChrisWiegman/kana-cli.git`
2. CD into the repo and run `make install`

Assuming you have Go properly setup with GOBIN in your system path, you should now be able to use Kana. Run `kana version` to test.

# Using Kana

At it's most basic you can start a zero-config Kana site by running `kana start` in your terminal. This will create a new Kana site based on your current directory name and open it in your default browser. If it is the first time you've run Kana it will also install it's root CA in your Mac's system store.

Kana relies on [Traefik](https://traefik.io) to map real domains to local sites. You can run as many sites as you need and each will be mapped to a subdomain of _sites.kana.li_.

## Start

`kana start` will start a kana site based on your current directory and open it in your browser.

To login to the new site use the following:

- _User Name_: **admin**
- _Password_: **password**

Note: these can be changed in the config. Please see below.

### Start options

`--plugin` will map the current directory as a plugin within the created site. Use this if you are developing a plugin.

`--theme` will map the current directory as a theme within the created site. Use this if you are developing a theme.

`--local` will create a directory called "wordpress" in the current directory and map it to the main WordPress site. This will allow you easy access, if you need it, to all the WordPress files (including any other installed plugins and themes) in your IDE.

If you do not specify the `local` flag you can find Kana's site files in `~/.config/kana/sites/<SITE NAME>/app`

`--xdebug` will start Xdebug on the site (see below for usage).

`--phpmyadmin` will start an instance of [phpMyAdmin](https://www.phpmyadmin.net) to allow for easier access to the database without needing external tools.

`--name` The name flag allows you to run an arbitrary site from anywhere. For example, if you already started and stopped a site from a directory called _test_ you can run `kana start --name=test` to start that site from anywhere. If you use the `name` flag on a new site it will create that site without a link to any local folder. This can be handy for testing a plugin or other configuration but not that none of the other start flags will apply.

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

## Stop

`kana stop` will stop the current site and, if no other sites are running, will shut down shared containers as well.

## Destroy

`kana destroy` will stop and destroy the current site. This is different than `stop` in that `stop` will leave the database and files it creates alone so you can start it again later. Once destroyed a site is irrecoverable.

## Open

`kana open` will open the site in your default browser

## wp-cli

`kana wp <WP-CLI COMMAND>` will execute a [wp-cli](https://wp-cli.org) command on your site. For example `kana wp plugin list` will list all the plugins on the site and their associated statuses

# Configuring Kana

The above commands will get an individual site up and running but there are a few more options to consider that can be changed for a given site or globally

## Global Config

Kana has a handful of options that apply to all new sites created with the app. You can adjust these with the `config` command as noted below:

`kana config` will list all changeable defaults for a new site. Currently these include the following:

- `admin.email` __admin@kanasite.localhost__ - the admin email address for the default admin account
- `admin.password` **password** - the default password used to login to WordPress
- `admin.username` **admin** - the default username used to login to WordPress
- `local` **false** - the default usage of the `local` start flag
- `php` **7.4** - the default PHP version used for new sites (currently 8.0, 8.1 and 8.2 are also supported)
- `type` **site** - the type of the Kana site you're starting. Current options are "site" "plugin" and "theme"
- `xdebug` **false** - the default usage of the `xdebug` start flag
- `phpmyadmin` **false** - the default usage of the `phpmyadmin` start flag

You can get or set any of the above options using a similar syntax to GIT's config. For example:

`kana config admin.email` will print the value of the admin.email setting
`kana config admin.email myemail@somedomain.com` will change the value of the admin.email setting to "myemail@somedomain.com".

The above syntax will allow you to change the defaults for any of the options listed

## Site Config

In addition to the global config, certain items above can be overridden for any given site. For a site without a `name` flag (as seen in the start command), simply create a _.kana.json_ file in the current directory. You can populate it with the following options:

- `local` **false** - the default usage of the `local` start flag
- `php` **7.4** - the default PHP version used for new sites (currently 8.0, 8.1 and 8.2 are also supported)
- `type` **site** - the type of the Kana site you're starting. Current options are "site" "plugin" and "theme"
- `xdebug` **false** - the default usage of the `xdebug` start flag
- `phpmyadmin` **false** - the default usage of the `phpmyadmin` start flag
- `plugins` **[]** - an array of plugins to install and activate when starting the new site. These are slugs from the Plugins section of WordPress.org.

### Export

`kana export` will create a _.kana.json_ configuration file in your current folder exporting the configuration of the current site including PHP version, active plugins and associated options as shown above

# Accessing the database directly

Currently there are two methods to access the database directly. First, use the `phpmyadmin` flag or setting (set to true) to add an instance of [phpMyAdmin](https://www.phpmyadmin.net) to your site. You can access this by appending **\*phpmyadmin-** to the beginning of your site domain. For example, if your site can get found at https://mysupersite.sites.kana.li you can access phpMyAdmin at https://phpmyadmin-mysupersite.sites.kana.li if you have enabled phpMyAdmin at site start.

You can also access the database directly by viewing the database port with `docker ps` and using the database port and the following configuration in the app of your choice:

- **Database host**: _kana\_`your site name`\_database_
- **Database name**: _wordpress_
- **Database user**: _wordpress_
- **Database password**: _wordpress_

# Using Xdebug

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

Note the above example will map the current folder as a plugin and maps the _wordpress_ folder as if the `local` flag was used. You may need to adjust these paths depending on your setup.

To trigger step debugging you'll also need the appropriate extension for your browser:

- [Xdebug Helper for Firefox](https://addons.mozilla.org/en-GB/firefox/addon/xdebug-helper-for-firefox/) ([source](https://github.com/BrianGilbert/xdebug-helper-for-firefox)).
- [Xdebug Helper for Chrome](https://chrome.google.com/extensions/detail/eadndfjplgieldjbigjakmdgkmoaaaoc) ([source](https://github.com/mac-cain13/xdebug-helper-for-chrome)).
- [XDebugToggle for Safari](https://apps.apple.com/app/safari-xdebug-toggle/id1437227804?mt=12) ([source](https://github.com/kampfq/SafariXDebugToggle)).

# This project is under active development

Note that I am using this project for my own work and it is under active development. Some of the things I'm currently working on include:

- Better site management commands
- Support for more xdebug modes (trace, etc)
- Much more clear prompts and messages on the commands themselves
- Writing a lot more tests (it's a personal project, I start where I can)
- A proper website for all this documentation (I already bought a domain, after all)
- Possible support for Docker alternatives

# Completely Uninstalling Kana

I hate apps that leave leftovers on your machine. When stopping a site all Docker resources except the images will be removed. To remove the app completely beyond that you'll want to delete the following:

1. Delete the application from your $GOBIN or system path (or run `brew uninstall kana` if installed via homebrew)
2. Delete the `~/.config/kana` folder which contains all site and app configuration
3. (Mac only) Delete the `Kana Development CA` certificate from the _System_ keychain in the _Keychain Access_ app
4. If installed via homebrew run `brew untap ChrisWiegman/kana` to remove the Homebrew tap

You can also safely remove any new images added however it is not a requirement. Many other apps might share those images leading to your system simply needing to download them again.
