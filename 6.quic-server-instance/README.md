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