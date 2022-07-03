package setup

import (
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/internal/docker"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

func SetupApp(controller *docker.Controller) {

	ensureAppConfig(controller.Config)
	ensureCerts(controller.Config)

}

// ensureAppConfig Ensures the application's config has been generated and is where it needs to be
func ensureAppConfig(kanaConfig config.KanaConfig) error {

	return writeFileArrayToDisk(configFiles, kanaConfig.AppDirectory)

}

// ensureCerts Ensures SSL certificates have been generated and are where they need to be
func ensureCerts(kanaConfig config.KanaConfig) {

	if err := os.MkdirAll(path.Join(kanaConfig.AppDirectory, "certs"), 0750); err != nil {
		panic(err)
	}

	minica.GenCerts(kanaConfig)
}
