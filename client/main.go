package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func newCertPool(caFile string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	pemByte, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	for {
		var block *pem.Block
		block, pemByte = pem.Decode(pemByte)
		if block == nil {
			return certPool, nil
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certPool.AddCert(cert)
	}

	//if !certPool.AppendCertsFromPEM(pemByte) {
	//	return nil, fmt.Errorf("can't add CA cert")
	//}
	//return certPool, nil
}

func NewTlsConfig(keyFile, certFile, caFile string) *tls.Config {
	var cfg tls.Config

	if keyFile == "" || certFile == "" {
		return &cfg
	}

	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalln("read pair file error:", err)
		return nil
	}

	cfg.Certificates = []tls.Certificate{tlsCert}

	if caFile != "" {
		cp, err := newCertPool(caFile)
		if err != nil {
			log.Fatalln("read cert file error:", err)
			return nil
		}

		cfg.RootCAs = cp
	}

	return &cfg
}

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: NewTlsConfig("", "", "./certs/ca.crt"),
		},
	}

	// baseUri := "https://localhost:3000"
	// baseUri := "https://127.0.0.1:3000"
	// baseUri := "https://192.168.1.6:3000"
	baseUri := "https://host.docker.internal:3000"

	resp, err := client.Get(baseUri + "/hello/world")
	if err != nil {
		log.Fatalln("get error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
