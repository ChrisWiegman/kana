package settings

import (
	_ "embed"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/subdomain.htaccess
var SubDomainMultisiteHtaccess string

//go:embed templates/subdirectory.htaccess
var SubDirectoryMultisiteHtaccess string

//go:embed templates/dynamic.toml
var DynamicToml string

//go:embed templates/traefik.toml
var TraefikToml string

//go:embed templates/kana-local-development.php
var KanaWordPressPlugin string

var configFiles = []File{
	{
		Name:        "dynamic.toml",
		Template:    DynamicToml,
		LocalPath:   "config/traefik",
		Permissions: os.FileMode(defaultFilePermissions),
	},
	{
		Name:        "traefik.toml",
		Template:    TraefikToml,
		LocalPath:   "config/traefik",
		Permissions: os.FileMode(defaultFilePermissions),
	},
}

// EnsureKanaPlugin ensures the Kana plugin file is in place and ready to go.
func EnsureKanaPlugin(siteDirectory, version, siteName string) error {
	pluginVars := PluginVersion{
		Version:  version,
		SiteName: siteName,
	}

	tmpl := template.Must(template.New("kanaPlugin").Parse(KanaWordPressPlugin))

	pluginPath := filepath.Join(siteDirectory, "wp-content", "mu-plugins")

	_, err := os.Stat(pluginPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(pluginPath, os.FileMode(defaultDirPermissions))
		if err != nil {
			return err
		}
	}

	myFile, err := os.Create(filepath.Join(pluginPath, "kana-local-development.php"))
	if err != nil {
		return err
	}

	return tmpl.Execute(myFile, pluginVars)
}

// GetDefaultFilePermissions returns the default directory permissions and the default file permissions.
func GetDefaultFilePermissions() (dirPerms, filePerms int) {
	return defaultDirPermissions, defaultFilePermissions
}

// GetHtaccess Returns the correct .htaccess file for the multisite type.
func GetHtaccess(multisite string) string {
	if multisite == "subdomain" {
		return SubDomainMultisiteHtaccess
	}

	return SubDirectoryMultisiteHtaccess
}

// esureStaticConfigFiles Ensures the application's static config files have been generated and are where they need to be.
func ensureStaticConfigFiles(appDirectory string) error {
	for _, file := range configFiles {
		filePath := filepath.Join(appDirectory, file.LocalPath)
		destFile := filepath.Join(appDirectory, file.LocalPath, file.Name)

		if err := os.MkdirAll(filePath, os.FileMode(defaultDirPermissions)); err != nil {
			return err
		}

		finalTemplate := []byte(file.Template)

		err := os.WriteFile(destFile, finalTemplate, file.Permissions)
		if err != nil {
			return err
		}
	}

	return nil
}
