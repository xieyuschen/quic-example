package main

import "net"

const (
	udpAddr="localhost:9000"
)

func main(){
	addr,_:=net.ResolveUDPAddr("udp",udpAddr)
	conn,err:=net.ListenUDP("udp",addr)
	if err!=nil{
		panic(err)
	}

}
