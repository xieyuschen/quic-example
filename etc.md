# ETC
This file records some useful acknowledgements about building large project during learning quic-go.

## Golang select priority
Golang select has no guarantee of the order of the cases if there are multiple cases ready, see 
[reference](https://stackoverflow.com/questions/47645808/how-does-select-work-when-multiple-channels-are-involved). If 
we want to set the priority of the cases, could refer to the following 
[link](https://stackoverflow.com/questions/46200343/force-priority-of-go-select-statement/46202533#46202533).

That's how quic-go solves the order between `error` and receiving packets:
```go
func (s *baseServer) run() {
	defer close(s.running)
	for {
		select {
		case <-s.errorChan:
			return
		default:
		}
		select {
		case <-s.errorChan:
			return
		case p := <-s.receivedPackets:
			if bufferStillInUse := s.handlePacketImpl(p); !bufferStillInUse {
				p.buffer.Release()
			}
		}
	}
}
```