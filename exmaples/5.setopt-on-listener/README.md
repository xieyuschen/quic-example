# Set opt on listener
The [net.ListenConfig](https://pkg.go.dev/net#ListenConfig) provides a way to set options for the listener. 
In the quic-go, there is no such api. Instead, it provides an API to listen with an established `net.PacketConn` 
with function [Listen](https://pkg.go.dev/github.com/lucas-clemente/quic-go#Listen).
