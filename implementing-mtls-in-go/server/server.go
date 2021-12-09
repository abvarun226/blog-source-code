package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

// HTTPServer will start a mTLS enabled httptest server and return the test server.
// It requires server certificate's public and private key files and root certificate's public key file as arguments.
func HTTPServer(serverPublicKey, serverPrivateKey, rootPublicKey string) *httptest.Server {
	// server certificate.
	serverCert, err := tls.LoadX509KeyPair(serverPublicKey, serverPrivateKey)
	if err != nil {
		return nil
	}

	// root certificate.
	rootCert, err := ioutil.ReadFile(rootPublicKey)
	if err != nil {
		log.Fatalf("failed to read root public key: %v", err)
	}
	rootCertPool := x509.NewCertPool()
	rootCertPool.AppendCertsFromPEM(rootCert)

	// httptest server with TLS config.
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "success!")
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    rootCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	server.StartTLS()
	return server
}

// SendRequest function will send a GET request to the server URL provided in the argument.
// It also requires client certificate's public and private key files and root certificate's public key file.
func SendRequest(serverURL, clientPublicKey, clientPrivateKey, rootPublicKey string) {
	// root certificate public key
	rootCert, errRead := ioutil.ReadFile(rootPublicKey)
	if errRead != nil {
		log.Fatalf("failed to read public key: %v", errRead)
	}
	publicPemBlock, _ := pem.Decode(rootCert)
	rootPubCrt, errParse := x509.ParseCertificate(publicPemBlock.Bytes)
	if errParse != nil {
		log.Fatalf("failed to parse public key: %v", errParse)
	}

	rootCertpool := x509.NewCertPool()
	rootCertpool.AddCert(rootPubCrt)

	// client certificates.
	cert, err := tls.LoadX509KeyPair(clientPublicKey, clientPrivateKey)
	if err != nil {
		log.Fatalf("failed to load client certificate: %v", err)
	}

	// http client with root and client certificates.
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      rootCertpool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get(serverURL)
	if err != nil {
		log.Printf("failed to GET: %v", err)
		return
	}

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		log.Printf("failed to read body: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("successful GET: %s", string(body))
}
