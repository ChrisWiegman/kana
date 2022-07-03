package setup

import (
	"os"
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

func EnsureCerts(kanaConfig config.KanaConfig) {

	if err := os.MkdirAll(path.Join(kanaConfig.AppDirectory, "certs"), 0750); err != nil {
		panic(err)
	}

	minica.GenCerts(kanaConfig)
}
