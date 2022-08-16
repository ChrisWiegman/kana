package setup

import (
	"os"
	"os/exec"
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

func SetupApp(appConfig config.AppConfig) error {

	err := ensureAppConfig(appConfig)
	if err != nil {
		return err
	}

	return ensureCerts(appConfig)

}

// ensureAppConfig Ensures the application's config has been generated and is where it needs to be
func ensureAppConfig(kanaConfig config.AppConfig) error {

	return writeFileArrayToDisk(configFiles, kanaConfig.AppHomeDirectory)

}

// ensureCerts Ensures SSL certificates have been generated and are where they need to be
func ensureCerts(kanaConfig config.AppConfig) error {

	createCert := false
	rootCert := path.Join(kanaConfig.AppHomeDirectory, "certs", kanaConfig.RootCert)

	_, err := os.Stat(rootCert)
	if err != nil && os.IsNotExist(err) {
		createCert = true
	}

	if createCert {

		err = os.MkdirAll(path.Join(kanaConfig.AppHomeDirectory, "certs"), 0750)
		if err != nil {
			return err
		}

		err = minica.GenCerts(kanaConfig)
		if err != nil {
			return err
		}

		installCertCommand := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", rootCert)
		return installCertCommand.Run()

	}

	return nil

}
