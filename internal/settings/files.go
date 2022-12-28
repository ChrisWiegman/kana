package settings

import (
	_ "embed"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/ChrisWiegman/kana-cli/pkg/minica"
)

type File struct {
	Name, Template, LocalPath string
	Permissions               os.FileMode
}

//go:embed templates/dynamic.toml
var DynamicToml string

//go:embed templates/traefik.toml
var TraefikToml string

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
			installCertCommand := exec.Command(
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
