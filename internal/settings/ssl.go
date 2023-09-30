package settings

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/ChrisWiegman/kana-cli/internal/console"
	"github.com/ChrisWiegman/kana-cli/pkg/minica"
)

var execCommand = exec.Command

const certOS = "darwin"

// EnsureSSLCerts Ensures SSL certificates have been generated and are where they need to be.
func (s *Settings) EnsureSSLCerts(consoleOutput *console.Console) error {
	createCert := false

	certPath := path.Join(s.AppDirectory, "certs")
	rootCert = path.Join(certPath, s.RootCert)

	_, err := os.Stat(rootCert)
	if err != nil && os.IsNotExist(err) {
		createCert = true
	}

	if createCert {
		err = os.MkdirAll(certPath, os.FileMode(defaultDirPermissions))
		if err != nil {
			return err
		}

		certInfo := minica.CertInfo{
			CertDir:    certPath,
			CertDomain: s.AppDomain,
			RootKey:    s.RootKey,
			RootCert:   s.RootCert,
			SiteCert:   s.SiteCert,
			SiteKey:    s.SiteKey,
		}

		err = minica.GenCerts(&certInfo)
		if err != nil {
			return err
		}
	}

	// If we're on Mac try to add the cert to the system trust.
	if s.SSL && runtime.GOOS == certOS {
		return TrustSSL(consoleOutput)
	}

	return nil
}

func TrustSSL(consoleOutput *console.Console) error {
	if runtime.GOOS != certOS {
		return fmt.Errorf("the trust command is only available for MacOS")
	}
	err := VerifySSLTrust()
	if err != nil {
		consoleOutput.Println("Adding Kana's SSL certificate to your system keychain. You will be promoted for your password.")

		installCertCommand := execCommand(
			"sudo",
			"security",
			"add-trusted-cert",
			"-d",
			"-r",
			"trustRoot",
			"-k",
			"/Library/Keychains/System.keychain",
			rootCert)

		return installCertCommand.Run()
	}

	return nil
}

func VerifySSLTrust() error {
	if runtime.GOOS == certOS {
		verifyCertCommand := execCommand(
			"security",
			"find-certificate",
			"-c",
			"Kana Development CA",
			"/Library/Keychains/System.keychain")

		return verifyCertCommand.Run()
	}

	return fmt.Errorf("the trust command is only available for MacOS")
}
