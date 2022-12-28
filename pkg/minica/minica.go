package minica

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"os"
	"path"
	"time"
)

type issuer struct {
	key  crypto.Signer
	cert *x509.Certificate
}

type CertInfo struct {
	CertDir    string
	CertDomain string
	RootKey    string
	RootCert   string
	SiteCert   string
	SiteKey    string
}

var fileOpenMode = 0600
var hundredYears = 100
var twoYears = 2
var thirtyDays = 30
var certBytes = 2048

func GenCerts(certInfo *CertInfo) error {
	caKey := path.Join(certInfo.CertDir, certInfo.RootKey)
	caCert := path.Join(certInfo.CertDir, certInfo.RootCert)
	domains := []string{
		fmt.Sprintf("*.%s", certInfo.CertDomain),
	}

	issuer, err := getIssuer(caKey, caCert)
	if err != nil {
		return err
	}

	_, err = sign(issuer, domains, certInfo.CertDir, certInfo.SiteCert, certInfo.SiteKey)

	return err
}

func getIssuer(keyFile, certFile string) (*issuer, error) {
	keyContents, keyErr := os.ReadFile(keyFile)
	certContents, certErr := os.ReadFile(certFile)

	if os.IsNotExist(keyErr) && os.IsNotExist(certErr) {
		err := makeIssuer(keyFile, certFile)
		if err != nil {
			return nil, err
		}

		return getIssuer(keyFile, certFile)
	} else if keyErr != nil {
		return nil, fmt.Errorf("%s (but %s exists)", keyErr, certFile)
	} else if certErr != nil {
		return nil, fmt.Errorf("%s (but %s exists)", certErr, keyFile)
	}

	key, err := readPrivateKey(keyContents)
	if err != nil {
		return nil, fmt.Errorf("reading private key from %s: %s", keyFile, err)
	}

	cert, err := readCert(certContents)
	if err != nil {
		return nil, fmt.Errorf("reading CA certificate from %s: %s", certFile, err)
	}

	equal, err := publicKeysEqual(key.Public(), cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("comparing public keys: %s", err)
	} else if !equal {
		return nil, fmt.Errorf("public key in CA certificate %s doesn't match private key in %s",
			certFile, keyFile)
	}

	return &issuer{key, cert}, nil
}

func readPrivateKey(keyContents []byte) (crypto.Signer, error) {
	block, _ := pem.Decode(keyContents)
	if block == nil {
		return nil, fmt.Errorf("no PEM found")
	} else if block.Type != "RSA PRIVATE KEY" && block.Type != "ECDSA PRIVATE KEY" {
		return nil, fmt.Errorf("incorrect PEM type %s", block.Type)
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func readCert(certContents []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certContents)
	if block == nil {
		return nil, fmt.Errorf("no PEM found")
	} else if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("incorrect PEM type %s", block.Type)
	}

	return x509.ParseCertificate(block.Bytes)
}

func makeIssuer(keyFile, certFile string) error {
	key, err := makeKey(keyFile)
	if err != nil {
		return err
	}

	_, err = makeRootCert(key, certFile)
	if err != nil {
		return err
	}

	return nil
}

func makeKey(filename string) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, certBytes)
	if err != nil {
		return nil, err
	}

	der := x509.MarshalPKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(fileOpenMode))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	err = pem.Encode(file, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	})

	if err != nil {
		return nil, err
	}

	return key, nil
}

func makeRootCert(key crypto.Signer, filename string) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

	if err != nil {
		return nil, err
	}

	skid, err := calculateSKID(key.Public())
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName:         "Kana Development CA",
			Country:            []string{"US"},
			Province:           []string{"Wyoming"},
			Locality:           []string{"Casper"},
			Organization:       []string{"Kana"},
			OrganizationalUnit: []string{"Development"},
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(hundredYears, 0, 0),

		SubjectKeyId:          skid,
		AuthorityKeyId:        skid,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(fileOpenMode))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(der)
}

func publicKeysEqual(a, b interface{}) (bool, error) {
	aBytes, err := x509.MarshalPKIXPublicKey(a)
	if err != nil {
		return false, err
	}

	bBytes, err := x509.MarshalPKIXPublicKey(b)
	if err != nil {
		return false, err
	}

	return bytes.Equal(aBytes, bBytes), nil
}

func calculateSKID(pubKey crypto.PublicKey) ([]byte, error) {
	spkiASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}

	_, err = asn1.Unmarshal(spkiASN1, &spki)
	if err != nil {
		return nil, err
	}

	skid := sha1.Sum(spki.SubjectPublicKey.Bytes) //nolint:gosec
	return skid[:], nil
}

func sign(iss *issuer, domains []string, certPath, siteCert, siteKey string) (*x509.Certificate, error) {
	cn := domains[0]

	key, err := makeKey(path.Join(certPath, siteKey))
	if err != nil {
		return nil, err
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		DNSNames: domains,
		Subject: pkix.Name{
			CommonName: cn,
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		// Set the validity period to 2 years and 30 days, to satisfy the iOS and
		// macOS requirements that all server certificates must have validity
		// shorter than 825 days:
		// https://derflounder.wordpress.com/2019/06/06/new-tls-security-requirements-for-ios-13-and-macos-catalina-10-15/
		NotAfter: time.Now().AddDate(twoYears, 0, thirtyDays),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, iss.cert, key.Public(), iss.key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path.Join(certPath, siteCert), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(fileOpenMode))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})

	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(der)
}
