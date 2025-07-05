package settings

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/pkg/minica"
)

var execCommand = exec.Command

// EnsureSSLCerts Ensures SSL certificates have been generated and are where they need to be.
func EnsureSSLCerts(appDirectory string, useSSL bool, consoleOutput *console.Console) error {
	createCert := false

	certPath := filepath.Join(appDirectory, "certs")
	rootCertFile := filepath.Join(certPath, rootCert)

	_, err := os.Stat(rootCertFile)
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
			CertDomain: domain,
			RootKey:    rootKey,
			RootCert:   rootCert,
			SiteCert:   siteCert,
			SiteKey:    siteKey,
		}

		err = minica.GenCerts(&certInfo)
		if err != nil {
			return err
		}
	}

	// If we're on Mac try to add the cert to the system trust.
	if useSSL && runtime.GOOS == certOS {
		return TrustSSL(rootCert, appDirectory, consoleOutput)
	}

	return nil
}

// TrustSSL Adds the Kana certificate to the Apple Keychain.
func TrustSSL(rootCert, appDirectory string, consoleOutput *console.Console) error {
	if runtime.GOOS != certOS {
		return fmt.Errorf("the trust command is only available for MacOS")
	}

	err := VerifySSLTrust()
	if err != nil {
		consoleOutput.Println("Adding Kana's SSL certificate to your system keychain. You will be promoted for your password.")

		certPath := filepath.Join(appDirectory, "certs")
		rootCertFile := filepath.Join(certPath, rootCert)

		installCertCommand := execCommand(
			"sudo",
			"security",
			"add-trusted-cert",
			"-d",
			"-r",
			"trustRoot",
			"-k",
			"/Library/Keychains/System.keychain",
			rootCertFile)

		err = installCertCommand.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// VerifySSLTrust verifies the SSL certificate has been added to that Apple Keychain.
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
