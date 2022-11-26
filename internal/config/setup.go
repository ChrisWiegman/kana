package config

import (
	"os"
	"os/exec"
	"path"

	"github.com/ChrisWiegman/kana-cli/pkg/minica"
)

// EnsureStaticConfigFiles Ensures the application's static config files have been generated and are where they need to be
func (c *Config) EnsureStaticConfigFiles() error {
	return writeFileArrayToDisk(configFiles, c.Directories.App)
}

// EnsureCerts Ensures SSL certificates have been generated and are where they need to be
func (c *Config) EnsureCerts() error {

	createCert := false
	rootCert := path.Join(c.Directories.App, "certs", c.App.RootCert)

	_, err := os.Stat(rootCert)
	if err != nil && os.IsNotExist(err) {
		createCert = true
	}

	if createCert {

		err = os.MkdirAll(path.Join(c.Directories.App, "certs"), 0750)
		if err != nil {
			return err
		}

		certInfo := minica.CertInfo{
			CertDir:    c.Directories.App,
			CertDomain: c.App.Domain,
			RootKey:    c.App.RootKey,
			RootCert:   c.App.RootCert,
			SiteCert:   c.App.SiteCert,
			SiteKey:    c.App.SiteKey,
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
