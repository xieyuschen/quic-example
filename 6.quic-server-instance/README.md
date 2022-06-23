# quic-server-instance
This chapter focus on how the server instance works.

## Server workflow
The server is designed and encapsulated well. It uses channel to pass data from different go routines which benefits 
from the mechanism of go routine.  

## Multiplexer
One interface and prepare one interface for multiplexing. Return the interface to the caller instead of struct. The 
interface usage forbids the caller to know the implementation and visits the data member.

```go
type multiplexer interface {
AddConn(c net.PacketConn, connIDLen int, statelessResetKey []byte, tracer logging.Tracer) (packetHandlerManager, error)
RemoveConn(indexableConn) error
}

// The connMultiplexer listens on multiple net.PacketConns and dispatches
// incoming packets to the connection handler.
type connMultiplexer struct {
	mutex sync.Mutex

	conns                   map[string] /* LocalAddr().String() */ connManager
	newPacketHandlerManager func(net.PacketConn, int, []byte, logging.Tracer, utils.Logger) (packetHandlerManager, error) // so it can be replaced in the tests

	logger utils.Logger
}
```

Inside `connMultiplexer`, there is a function `newPacketHandlerManager` which is used to create a new packet handler 
manager. The reason for this is that the handler should bind with the `connMultiplexer` but shouldn't as an interface. 

The function could of course call directly instead of being stored in the `connMultiplexer` but this is not a good as 
it introduces complexity if we want to change the handler to a different implementation.

Of course, this might be the habit of the user.

## Server-ConnMultiplexer-PacketHandlerManager
The abstraction in quic-go does quite well. The server start and then enroll ConnMultiplexer and PacketHandlerManager. 
All of them connect the other through interface, instead of concreted struct.

- Server:  
The server starts a `run` method in another go routine to detect whether there is a new packet to handle, if there is,
it will call `handlePacketImpl` to handle the packet.

- PacketHandlerManager:  
The packet handler manager is responsible for handling the packet. It is a stateless handler. It is a singleton. It 
executes the `newPacketHandlerMap` to create a new packet handler map and then call `go m.listen()` to listen on the 
`packetHandlerManager`. 

If there comes new packets, it will use channel to inform the server by `s.receivedPackets <- p` and it is captured by 
the endless loop in `run` method.  
**It's the duty to call `s.handlePacketImpl` to handle the packet**. Each connection has its own packet handler/manager.

- ConnMultiplexer:  
The ConnMultiplexer is responsible for listening on the net.PacketConn and creating a new PacketHandlerManager. As there 
comes packets first, a connection later. **It is a singleton and a stateless handler**.

The `ConnMultiplexer` doesn't charge for handling connection details. It just listens on the net.PacketConn and create 
`packetHandlerManager` for each connection. After the connection being established, the connection will be created by 
the server function `func (s *baseServer) handleInitialImpl(p *receivedPacket, hdr *wire.Header) error` and start a new 
go routine to handle the connection packet with function call `go conn.run()`.