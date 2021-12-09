package certificate

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"time"
)

// Verification function will verify the client and server certificates using the root certificate that was provided in the function arguments.
// If generateCertificate is set to true, it will also generate new client and server certificates that is signed by root certificate.
func Verification(serverName string, clientPubKey, clientPrivKey, serverPubKey, serverPrivKey, rootPubKey, rootPrivKey string, generateCertificate bool) {
	ca, errCA := ReadCertificateAuthority(rootPubKey, rootPrivKey)
	if errCA != nil {
		log.Fatalf("failed to read ca certificate: %v", errCA)
	}

	if generateCertificate {
		// generate and sign client certificate using root certificate.
		if err := GenerateAndSignCertificate(ca, clientPubKey, clientPrivKey); err != nil {
			log.Fatalf("failed to generate client certificate: %v", err)
		}

		// generate and sign server certificate using root certificate.
		if err := GenerateAndSignCertificate(ca, serverPubKey, serverPrivKey); err != nil {
			log.Fatalf("failed to generate server certificate: %v", err)
		}
	}

	clientCert, errCert := ReadCertificate(clientPubKey, clientPrivKey)
	if errCert != nil {
		log.Fatalf("failed to read certificate: %v", errCert)
	}

	serverCert, errCert := ReadCertificate(serverPubKey, serverPrivKey)
	if errCert != nil {
		log.Fatalf("failed to read certificate: %v", errCert)
	}

	roots := x509.NewCertPool()
	roots.AddCert(ca.PublicKey)

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
		DNSName:       serverName,
	}

	// verify client certificate; return err on failure.
	if _, err := clientCert.PublicKey.Verify(opts); err != nil {
		log.Fatalf("failed to verify client certificate: %v", err)
	}

	// verify server certificate; return err on failure.
	if _, err := serverCert.PublicKey.Verify(opts); err != nil {
		log.Fatalf("failed to verify server certificate: %v", err)
	}

	log.Print("client and server cert verification succeeded")
}

// ReadCertificate reads and parses the certificates from files provided as argument to this function.
func ReadCertificate(publicKeyFile, privateKeyFile string) (*KeyPair, error) {
	cert := new(KeyPair)

	privKey, errRead := ioutil.ReadFile(privateKeyFile)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read private key: %w", errRead)
	}

	privPemBlock, _ := pem.Decode(privKey)

	// Note that we use PKCS1 to parse the private key here.
	parsedPrivKey, errParse := x509.ParsePKCS1PrivateKey(privPemBlock.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", errParse)
	}

	cert.PrivateKey = parsedPrivKey

	pubKey, errRead := ioutil.ReadFile(publicKeyFile)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read public key: %w", errRead)
	}

	publicPemBlock, _ := pem.Decode(pubKey)

	parsedPubKey, errParse := x509.ParseCertificate(publicPemBlock.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", errParse)
	}

	cert.PublicKey = parsedPubKey

	return cert, nil
}

// ReadCertificateAuthority reads and parses the root certificate from files provided as argument to this function.
func ReadCertificateAuthority(publicKeyFile, privateKeyFile string) (*KeyPair, error) {
	root := new(KeyPair)

	rootKey, errRead := ioutil.ReadFile(privateKeyFile)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read private key: %w", errRead)
	}

	privPemBlock, _ := pem.Decode(rootKey)

	// Note that we use PKCS8 to parse the private key here.
	rootPrivKey, errParse := x509.ParsePKCS8PrivateKey(privPemBlock.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", errParse)
	}

	root.PrivateKey = rootPrivKey.(*rsa.PrivateKey)

	rootCert, errRead := ioutil.ReadFile(publicKeyFile)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read public key: %w", errRead)
	}

	publicPemBlock, _ := pem.Decode(rootCert)

	rootPubCrt, errParse := x509.ParseCertificate(publicPemBlock.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", errParse)
	}

	root.PublicKey = rootPubCrt

	return root, nil
}

// GenerateAndSignCertificate method will use the root certificate's public and private key to generate a certificate and sign it.
// The certificate's public and private keys will be stored in the files provided as argument to this function.
func GenerateAndSignCertificate(root *KeyPair, publicKeyFile, privateKeyFile string) error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"ayada"},
			Country:      []string{"US"},
			Province:     []string{"California"},
			Locality:     []string{"San Francisco"},
			CommonName:   "localhost",
		},
		DNSNames:     []string{"localhost", "ayada.dev"},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, root.PublicKey, &certPrivKey.PublicKey, root.PrivateKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	if err := ioutil.WriteFile(publicKeyFile, certPEM.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write to public key file: %w", err)
	}

	if err := ioutil.WriteFile(privateKeyFile, certPrivKeyPEM.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write to private key file: %w", err)
	}

	return nil
}

type KeyPair struct {
	PublicKey  *x509.Certificate
	PrivateKey *rsa.PrivateKey
}
