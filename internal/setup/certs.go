package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/ChrisWiegman/kana/internal/config"
)

var caCert = "certs/kana.ca.pem"
var caKey = "certs/kana.ca.key"

func GenerateCA() {

	fmt.Println("Checking for Root CA...")

	appConfigPath, err := config.GetConfigRoot()
	if err != nil {
		panic(err)
	}

	certDir := path.Join(appConfigPath, "certs")
	caCertFile := path.Join(appConfigPath, caCert)
	caKeyFile := path.Join(appConfigPath, caKey)

	_, err = os.Stat(caKeyFile)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
	}

	if os.IsNotExist(err) {

		fmt.Println("Root CA not found. Generating Root CA...")

		os.MkdirAll(certDir, 0700)

		err = exec.Command(
			"openssl",
			"genrsa",
			"-out",
			caKeyFile,
			"4096").Run()
		if err != nil {
			fmt.Println(err)
		}

		err = exec.Command(
			"openssl",
			"req",
			"-x509",
			"-new",
			"-nodes",
			"-key",
			caKeyFile,
			"-sha256",
			"-days",
			"7300",
			"-out",
			caCertFile,
			"-subj",
			"/C=US/ST=Florida/L=Sarasota/O=Kana/OU=Development/CN=Kana Development CA").Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}
