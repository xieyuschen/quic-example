package util

import (
	"path"
	"runtime"
)

var (
	certFile, keyFile string
	sslLogFile        string
)

const (
	listenAddress    = "127.0.0.1:4243"
	certPemPath      = "../cert/cert.pem"
	privKeyPath      = "../cert/priv.key"
	sslOutputLogPath = "../ssl.log"
)

func GetCertFilesPath() (certPem, privKey string) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current frame")
	}

	certFolderPath := path.Dir(filename)

	return path.Join(certFolderPath, certPemPath), path.Join(certFolderPath, privKeyPath)
}
