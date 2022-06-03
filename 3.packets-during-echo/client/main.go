package main

import (
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"log"
)

const (
	peerAddress = "127.0.0.1:4243"
	message     = "Hello world"
)

func main() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"echo-quic-demo"},
	}
	conn, err := quic.DialAddr(peerAddress, tlsConfig, nil)
	if err != nil {
		panic(err)
	}
	stream, err := conn.OpenStream()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
	buf := make([]byte, len(message))
	_, err =stream.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Client: Got '%s'\n", buf)
}
