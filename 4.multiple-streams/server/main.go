package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/xieyuschen/quic-example/util"
	"log"
)

const (
	listenAddress = "127.0.0.1:4244"
)

func main() {
	certFile, keyFile := util.GetCertFilesPath()
	var err error
	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: certs,
		NextProtos:   []string{"multiple-streams-quic-demo"},
	}

	fmt.Println("Quic server is running")

	listener, err := quic.ListenAddr(listenAddress, tlsConfig, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Printf("encounter error when accept: %s\n", err)
			continue
		}
		go handleQuicConnection(conn)
	}
}

func handleQuicConnection(conn quic.Connection) {
	for {
		// why AcceptStream receives a context?
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("failed to accept a quic stream,err: %s\n", err)
			continue
		}
		go handleQuicStream(stream)
	}
}

func handleQuicStream(stream quic.Stream) {

}
