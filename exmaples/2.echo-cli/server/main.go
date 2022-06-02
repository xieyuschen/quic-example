package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"math/big"
)

const (
	listenAddress = "127.0.0.1:4242"
)
func main(){
	var writerType string
	flag.StringVar(&writerType,"writer","normal",
		"writer type is used to specify how server writes to the stream. twice will write it twice")
	flag.Parse()

	fmt.Println("Quic server is running, it will exit after a stream is done")
	listener,err:=quic.ListenAddr(listenAddress,generateTlsConfig(),nil)
	if err != nil {
		log.Fatalln(err)
	}
	conn,err:=listener.Accept(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Connection is established\n")
	stream,err:= conn.AcceptStream(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Create a new stream id: %d\n",stream.StreamID())
	// good point!
	_,err=io.Copy(loggingWriter{Writer:stream,writeType: writerType},stream)

	if err != nil {
		fmt.Printf("stream %d is closed, err:%s\n",stream.StreamID(),err)
	}
}

// generateTlsConfig generates the key pair based on the RSA and certificates the key-pair
// with x509 algorithm.
func generateTlsConfig()*tls.Config{
	key,err:=rsa.GenerateKey(rand.Reader,4096)
	if err != nil {
		panic(err)
	}
	template:=x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER,err:= x509.CreateCertificate(rand.Reader,&template,&template,&key.PublicKey,key)
	if err != nil {
		panic(err)
	}

	keyPEM:=pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY",Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM:=pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE",Bytes: certDER})

	tlsCert,err:=tls.X509KeyPair(certPEM,keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{ tlsCert},
		NextProtos: []string{"echo-quic-demo"},
	}
}

// loggingWriter is a good example that how to wrap a type
// good point!
type loggingWriter struct {
	writeType string
	io.Writer
}


func (w loggingWriter) Write(b []byte)  (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	if w.writeType=="twice"{
		l:=len(b)
		n,err:=w.Writer.Write(b[:l/2])
		nn,err:=w.Writer.Write(b[l/2:])
		return n+nn,err
	}
	return w.Writer.Write(b)
}

