package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
)

const (
	serverDomain     = "https://localhost:4244"
	certPemPath      = "../../cert/cert.pem"
	privKeyPath      = "../../cert/priv.key"
	sslOutputLogPath = "../ssl.log"
)

var (
	certFile, keyFile, sslLogFile string
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current frame")
	}

	certPath := path.Dir(filename)
	certFile, keyFile, sslLogFile = path.Join(certPath, certPemPath),
		path.Join(certPath, privKeyPath),
		path.Join(certPath, sslLogFile)
}
func main() {
	keyLogFile := flag.String("keylog", "", "key log file")
	insecure := flag.Bool("insecure", false, "skip certificate verification")
	flag.Parse()

	fmt.Println("Start http3 client")

	var keyLog io.Writer
	if len(*keyLogFile) > 0 {
		f, err := os.Create(*keyLogFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		keyLog = f
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	caCertRaw, err := ioutil.ReadFile(certFile)
	if err != nil {
		panic(err)
	}
	if ok := pool.AppendCertsFromPEM(caCertRaw); !ok {
		panic("Could not add root ceritificate to pool.")
	}

	var qconf quic.Config

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: *insecure,
			KeyLogWriter:       keyLog,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	doneCh := make(chan struct{}, 1)
	go func(addr string) {
		fmt.Printf("GET %s", addr)
		rsp, err := hclient.Get(addr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Got response for %s: %#v", addr, rsp)

		body := &bytes.Buffer{}
		_, err = io.Copy(body, rsp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Response Body: %d bytes", body.Len())
		fmt.Printf("Response Body:")
		fmt.Printf("%s", body.Bytes())
		doneCh <- struct{}{}
	}(serverDomain)
	<-doneCh
}
