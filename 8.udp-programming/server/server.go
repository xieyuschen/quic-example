package main

import (
	"fmt"
	"net"
)

const (
	udpAddr = "localhost:9000"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", udpAddr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 2048)
	for {
		n, ad, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Server receives data: %s from address %s\n", string(buf[:n]), ad)
		_, err = conn.WriteToUDP([]byte("echo back"), ad)
		if err != nil {
			fmt.Println(err)
		}
	}
}
