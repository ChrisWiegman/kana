package settings

import (
	_ "embed"
	"os"
	"os/exec"
	"path"

	"github.com/ChrisWiegman/kana-cli/pkg/minica"
)

type File struct {
	Name, Template, LocalPath string
	Permissions               os.FileMode
}

//go:embed templates/dynamic.toml
var DYNAMIC_TOML string

//go:embed templates/traefik.toml
var TRAEFIK_TOML string

var configFiles = []File{
	{
		Name:        "dynamic.toml",
		Template:    DYNAMIC_TOML,
		LocalPath:   "config/traefik",
		Permissions: 0644,
	},
	{
		Name:        "traefik.toml",
		Template:    TRAEFIK_TOML,
		LocalPath:   "config/traefik",
		Permissions: 0644,
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

		err = os.MkdirAll(certPath, 0750)
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

		err = minica.GenCerts(certInfo)
		if err != nil {
			return err
		}

		installCertCommand := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", rootCert)
		return installCertCommand.Run()
	}

	return nil
}

// EnsureStaticConfigFiles Ensures the application's static config files have been generated and are where they need to be
func (s *Settings) EnsureStaticConfigFiles() error {

	for _, file := range configFiles {

		filePath := path.Join(s.AppDirectory, file.LocalPath)
		destFile := path.Join(s.AppDirectory, file.LocalPath, file.Name)

		if err := os.MkdirAll(filePath, 0750); err != nil {
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
