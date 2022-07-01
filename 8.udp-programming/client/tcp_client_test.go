package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

const (
	tcpAddr="localhost:9001"
)

func TestClient(t *testing.T) {
	c,err:=net.Dial("tcp",tcpAddr)
	if err!=nil{
		panic(err)
	}

	c.Write([]byte("hello"))
	buf:=make([]byte,2048)
	n,_:=c.Read(buf)
	fmt.Println(string(buf[:n]))
	time.Sleep(time.Minute)
}
