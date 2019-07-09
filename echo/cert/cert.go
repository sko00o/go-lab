package cert

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

var (
	keyPair  *tls.Certificate
	certPool *x509.CertPool
)

// GetCert returns a certificate pair and pool
func GetCert() (*tls.Certificate, *x509.CertPool) {
	serverCrt, err := ioutil.ReadFile("cert/server.crt")
	if err != nil {
		log.Fatal(err)
	}

	serverKey, err := ioutil.ReadFile("cert/server.key")
	if err != nil {
		log.Fatal(err)
	}

	pair, err := tls.X509KeyPair(serverCrt, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	keyPair = &pair
	certPool = x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(serverCrt)
	if !ok {
		log.Fatal("bad certs")
	}

	return keyPair, certPool
}
