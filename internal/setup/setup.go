package setup

import (
	"os"
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

	return writeFileArrayToDisk(configFiles, kanaConfig.AppDirectory)

}

// ensureCerts Ensures SSL certificates have been generated and are where they need to be
func ensureCerts(kanaConfig config.AppConfig) error {

	err := os.MkdirAll(path.Join(kanaConfig.AppDirectory, "certs"), 0750)
	if err != nil {
		return err
	}

	return minica.GenCerts(kanaConfig)

}
