package settings

import (
	_ "embed"
	"os"
	"os/exec"
	"path"
	"runtime"
	"text/template"

	"github.com/ChrisWiegman/kana-cli/pkg/minica"
)

type File struct {
	Name, Template, LocalPath string
	Permissions               os.FileMode
}

type KanaPluginVars struct {
	SiteName, Version string
}

//go:embed templates/dynamic.toml
var DynamicToml string

//go:embed templates/traefik.toml
var TraefikToml string

//go:embed templates/kana-local-development.php
var KanaWordPressPlugin string

var execCommand = exec.Command

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

// EnsureSSLCerts Ensures SSL certificates have been generated and are where they need to be
func (s *Settings) EnsureSSLCerts() error {
	createCert := false
	certPath := path.Join(s.AppDirectory, "certs")
	rootCert := path.Join(certPath, s.RootCert)

	_, err := os.Stat(rootCert)
	if err != nil && os.IsNotExist(err) {
		createCert = true
	}

	if createCert {
		err = os.MkdirAll(certPath, os.FileMode(defaultDirPermissions))
		if err != nil {
			return err
		}

		certInfo := minica.CertInfo{
			CertDir:    certPath,
			CertDomain: s.AppDomain,
			RootKey:    s.RootKey,
			RootCert:   s.RootCert,
			SiteCert:   s.SiteCert,
			SiteKey:    s.SiteKey,
		}

		err = minica.GenCerts(&certInfo)
		if err != nil {
			return err
		}

		// If we're on Mac try to add the cert to the system trust
		if runtime.GOOS == "darwin" {
			installCertCommand := execCommand(
				"sudo",
				"security",
				"add-trusted-cert",
				"-d",
				"-r",
				"trustRoot",
				"-k",
				"/Library/Keychains/System.keychain",
				rootCert)
			return installCertCommand.Run()
		}
	}

	return nil
}

// EnsureStaticConfigFiles Ensures the application's static config files have been generated and are where they need to be
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
