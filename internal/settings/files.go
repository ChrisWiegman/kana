package settings

import (
	_ "embed"
	"os"
	"path"
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
func (s *Settings) EnsureKanaPlugin(appDir string) error {
	pluginVars := KanaPluginVars{
		Version:  "1.0.0",
		SiteName: "my-site",
	}

	tmpl := template.Must(template.New("kanaPlugin").Parse(KanaWordPressPlugin))

	pluginPath := path.Join(appDir, "wp-content", "mu-plugins")

	_, err := os.Stat(pluginPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(pluginPath, os.FileMode(defaultDirPermissions))
		if err != nil {
			return err
		}
	}

	myFile, err := os.Create(path.Join(pluginPath, "kana-local-development.php"))
	if err != nil {
		return err
	}

	return tmpl.Execute(myFile, pluginVars)
}

// EnsureStaticConfigFiles Ensures the application's static config files have been generated and are where they need to be.
func (s *Settings) EnsureStaticConfigFiles() error {
	for _, file := range configFiles {
		filePath := path.Join(s.AppDirectory, file.LocalPath)
		destFile := path.Join(s.AppDirectory, file.LocalPath, file.Name)

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

func (s *Settings) GetHtaccess() string {
	if s.Multisite == "subdomain" {
		return SubDomainMultisiteHtaccess
	}

	return SubDirectoryMultisiteHtaccess
}
