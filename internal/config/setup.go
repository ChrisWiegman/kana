package config

import (
	"fmt"
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
	certPath := path.Join(c.Directories.App, "certs")
	rootCert := path.Join(certPath, c.App.RootCert)

	_, err := os.Stat(rootCert)
	if err != nil && os.IsNotExist(err) {
		createCert = true
	}

	fmt.Println(rootCert)

	if createCert {

		err = os.MkdirAll(certPath, 0750)
		if err != nil {
			return err
		}

		certInfo := minica.CertInfo{
			CertDir:    certPath,
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
