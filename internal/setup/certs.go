package setup

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

func EnsureCerts(kanaConfig config.KanaConfig) {

	fmt.Println("Checking for Root CA...")

	certDir := path.Join(kanaConfig.ConfigRoot, "certs")

	if err := os.MkdirAll(certDir, 0750); err != nil {
		panic(err)
	}

	minica.GenCerts(certDir, kanaConfig)
}
