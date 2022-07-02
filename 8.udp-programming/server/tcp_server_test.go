package main

import (
	"io"
	"net"
	"testing"
)

const (
	tcpAddr = "localhost:9001"
)

func TestTcpServer(t *testing.T) {
	li, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	for {
		c, _ := li.Accept()
		io.Copy(c, c)
	}

}
