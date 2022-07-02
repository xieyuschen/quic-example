package main

import (
	"fmt"
	"net"
)

const (
	udpAddr = "localhost:9000"
)

func main() {
	conn, err := net.Dial("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("hello"))
	fmt.Println("Client writes data: hello")
}
