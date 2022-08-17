package appSetup

import (
	"os"
	"os/exec"
	"path"

	"github.com/ChrisWiegman/kana/internal/appConfig"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

func SetupApp(staticConfig appConfig.StaticConfig) error {

	err := ensureStaticConfig(staticConfig)
	if err != nil {
		return err
	}

	return ensureCerts(staticConfig)

}

// ensureStaticConfig Ensures the application's static config files have been generated and are where they need to be
func ensureStaticConfig(staticConfig appConfig.StaticConfig) error {

	return writeFileArrayToDisk(configFiles, staticConfig.AppDirectory)

}

// ensureCerts Ensures SSL certificates have been generated and are where they need to be
func ensureCerts(staticConfig appConfig.StaticConfig) error {

	createCert := false
	rootCert := path.Join(staticConfig.AppDirectory, "certs", staticConfig.RootCert)

	_, err := os.Stat(rootCert)
	if err != nil && os.IsNotExist(err) {
		createCert = true
	}

	if createCert {

		err = os.MkdirAll(path.Join(staticConfig.AppDirectory, "certs"), 0750)
		if err != nil {
			return err
		}

		err = minica.GenCerts(staticConfig)
		if err != nil {
			return err
		}

		installCertCommand := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", rootCert)
		return installCertCommand.Run()

	}

	return nil

}