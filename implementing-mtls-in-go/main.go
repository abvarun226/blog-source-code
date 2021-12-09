package main

import (
	"github.com/abvarun226/mtls-proxy/certificate"
	"github.com/abvarun226/mtls-proxy/server"
)

const (
	clientPublicKey  = "certs/client.pem"
	clientPrivateKey = "certs/client.key"

	serverPublicKey  = "certs/server.pem"
	serverPrivateKey = "certs/server.key"

	rootPublicKey  = "certs/ca.pem"
	rootPrivateKey = "certs/ca.key"
)

func main() {
	certificate.Verification("localhost", clientPublicKey, clientPrivateKey, serverPublicKey, serverPrivateKey, rootPublicKey, rootPrivateKey, false)

	srv := server.HTTPServer(serverPublicKey, serverPrivateKey, rootPublicKey)
	defer srv.Close()

	server.SendRequest(srv.URL, clientPublicKey, clientPrivateKey, rootPublicKey)
}
