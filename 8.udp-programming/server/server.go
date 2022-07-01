package main

import (
	"fmt"
	"net"
)

const (
	udpAddr="localhost:9000"
)

func main(){
	addr,_:=net.ResolveUDPAddr("udp",udpAddr)
	conn,err:=net.ListenUDP("udp",addr)
	if err!=nil{
		panic(err)
	}
	buf:=make([]byte,2048)
	for{
		n,err:=conn.Read(buf)
		if err!=nil{
			panic(err)
		}
		fmt.Println("Server receives data: ",string(buf[:n]))


		//panic: write udp 127.0.0.1:9000: write: destination address required
		//_,err=conn.Write(buf)
		//if err!=nil{
		//	panic(err)
		//}
	}
}
