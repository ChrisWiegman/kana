package setup

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

var rootKey = "kana.root.key"
var rootCert = "kana.root.pem"
var siteCert = "kana.site.pem"
var siteKey = "kana.site.key"

func EnsureCerts(kanaConfig config.KanaConfig) {

	fmt.Println("Checking for Root CA...")

	certDir := path.Join(kanaConfig.ConfigRoot, "certs")

	if err := os.MkdirAll(certDir, 0750); err != nil {
		panic(err)
	}

	minica.GenCerts(certDir, rootKey, rootCert, siteCert, siteKey, kanaConfig)
}
