package setup

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

var caCert = "certs/kana.ca.pem"
var caKey = "certs/kana.ca.key"

func EnsureCerts() {

	fmt.Println("Checking for Root CA...")

	appConfigPath, err := config.GetConfigRoot()
	if err != nil {
		panic(err)
	}

	certDir := path.Join(appConfigPath, "certs")

	if err := os.MkdirAll(certDir, 0750); err != nil {
		panic(err)
	}

	minica.GenCerts(certDir)
}
