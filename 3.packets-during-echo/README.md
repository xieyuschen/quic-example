# 3.packets-during-echo

This section aims to learn about the packets during the echo example. 
If we want to use wireshark to capture and analysis packets, 
the key pair and certificate should be fixed as we need it to decrypt the
packets. As a result, the cert files are stored under `/cert` folder.

## How to decrypt the packets
As the quic standard document has been released in 2021.05, wireshark 
could analysis it automatically. 

**Note that you need to use a version that supports RFC draft-27**.  

When you capture the quic packets, you might see this error information 
in wireshark. It's caused by you don't have the master secret of the TLS.
```
Expert Info (Warning/Decryption): 
Failed to create decryption context: Secrets are not available
```

You could configure the TLS config and get the key, then load this key 
to wireshark so you can see the decrypted packets content.  
[Issue 5](https://github.com/xieyuschen/quic-example/issues/5) has tracked 
it.  

## Quic Packet format
Quic packet format is defined by RFC 9000 and you can view it to find 
more. The topic here aims to introduce briefly for the following packets
analysis.

### Type of quic headers: Long and short header    
Unlike TCP where the packet header format is fixed, QUIC
has two types of packet headers. 
- Long header:  
  QUIC packets for connection establishment need to contain several pieces of information, 
  it uses the long header format. 

    
- Short header:
  The short header format that is used after the handshake is completed. 
  Once a connection is established, only certain header fields are necessary, the subsequent 
  packets use the short header format for efficiency.
  
### Packet format
- Initial packet:
- Retry packet:  
- Handshake packet:  
- Protected payload packet:  

## Handshake procedures

### Step1:Initial packet

The packet content is consisted by the following parts:

- IPv4 header and UDP(user datagram protocol) header.  
  IP header contains 20 bytes(or five 32-bit increments with max 24 bytes). The UDP header is 8 bytes length.
  Note that in the hex displaying mode, every hex character stands for 4 bits(2^4).

- QUIC IETF.  
  The RFC says:
  > Initial packets use an AEAD function, the keys for which are derived using a value that is visible
  > on the wire. Initial packets therefore do not have effective confidentiality protection. Initial
  > protection exists to ensure that the sender of the packet is on the network path. Any entity that
  > receives an Initial packet from a client can recover the keys that will allow them to both read the
  > contents of the packet and generate Initial packets that will be successfully authenticated at
  > either endpoint. The AEAD also protects Initial packets against accidental modification.

The Initial packet format is:

```
Initial Packet {
  Header Form (1) = 1,
  Fixed Bit (1) = 1,
  Long Packet Type (2) = 0,
  Reserved Bits (2),
  Packet Number Length (2),
  Version (32),
  Destination Connection ID Length (8),
  Destination Connection ID (0..160),
  Source Connection ID Length (8),
  Source Connection ID (0..160),
  Token Length (i),
  Token (..),
  Length (i),
  Packet Number (8..32),
  Packet Payload (8..),
}
```

If using `tcpdump` to analyze the packets, remember to remove the IP header and udp header. Here are some details about
a quic initial packet:
![img.png](img/initial-overview.png)
There are some interesting variables in the initial format:

- Destination/Source connection ID:  
  When the client wants to establish a new connection, the `Destination Connection ID` is valid but the `Source Connection ID` is invalid as null. However, in
  [7.2.3-Negotiating Connection IDs](https://www.rfc-editor.org/rfc/rfc9000.html#section-7.2-3):

  > When an Initial packet is sent by a client that has not previously received an Initial or Retry packet from the 
  > server, the client populates the Destination Connection ID field with an unpredictable value. This Destination 
  > Connection ID MUST be at least 8 bytes in length. Until a packet is received from the server, the client MUST use 
  > the same Destination Connection ID value on all packets in this connection.

  As a result, the finial connection ID will be generated from the server/
  
- token:  
  A variable-length integer specifying the length of the Token field, in bytes.
  This value is 0 if no token is present. Initial packets sent by the server MUST set the Token Length field to 0;
  clients that receive an Initial packet with a non-zero Token Length field MUST either discard the packet or
  generate a connection error of type PROTOCOL_VIOLATION.

- CRYPTO:
  The first packet sent by a client always includes a CRYPTO frame that contains the start or all of the first
  cryptographic handshake message. The first CRYPTO frame sent always begins at an offset of 0.
  The CRYPTO frame encapsulates a TLSv1.3 `ClientHello` directly.
  ![img.png](img/crypto.png)

### Step2:Retry packet

For the new one, it's amazing that a retry packet is behind of an initial packet. Let's look up RFC to find out **why 
the client side gets a RETRY packet**.
In [RFC9000-7-Negotiating Connection IDs](https://www.rfc-editor.org/rfc/rfc9000.html#name-negotiating-connection-ids)
> The Destination Connection ID field from the first Initial packet sent by a client is used to determine packet 
> protection keys for Initial packets. These keys change after receiving a Retry packet.

**This step aims to verify the connection valid for avoiding traffic amplification attack**. The token will be sent to 
client with a `RETRY` packet.

The code implement details is in 
`server.go:395 func (s *baseServer) handleInitialImpl(p *receivedPacket, hdr *wire.Header) error`. We could check that 
in which case the server side will send a `RETRY` packet.
```go
	if !s.config.AcceptToken(p.remoteAddr, token) {
		go func() {
			defer p.buffer.Release()
			if token != nil && token.IsRetryToken {
				if err := s.maybeSendInvalidToken(p, hdr); err != nil {
					s.logger.Debugf("Error sending INVALID_TOKEN error: %s", err)
				}
				return
			}
			if err := s.sendRetry(p.remoteAddr, hdr, p.info); err != nil {
				s.logger.Debugf("Error sending Retry: %s", err)
			}
		}()
		return nil
	}
```
As a result, the function call `s.config.AcceptToken(p.remoteAddr, token)` decides to send `RETRY` or not. Add a 
breakpoint to see the details. The server side checks the client initial packet by the following code and this is a 
great interface usage as it abstract(todo).
```go
var defaultAcceptToken = func(clientAddr net.Addr, token *Token) bool {
	if token == nil {
		return false
	}
	validity := protocol.TokenValidity
	if token.IsRetryToken {
		validity = protocol.RetryTokenValidity
	}
	if time.Now().After(token.SentTime.Add(validity)) {
		return false
	}
	var sourceAddr string
	if udpAddr, ok := clientAddr.(*net.UDPAddr); ok {
		sourceAddr = udpAddr.IP.String()
	} else {
		sourceAddr = clientAddr.String()
	}
	return sourceAddr == token.RemoteAddr
}
```
If it needs send a RETRY packet, the server side construct a token to send back:
```go
func (s *baseServer) sendRetry(remoteAddr net.Addr, hdr *wire.Header, info *packetInfo) error {
  // ignore some lines
  token, err := s.tokenGenerator.NewRetryToken(remoteAddr, hdr.DestConnectionID, srcConnID)
  // ignore some lines
}
```

### Step3:The second Initial packet
The second Initial packet takes the token and connection ID from the RETRY packet sent by server. It sends a CRYPTO 
packet with a TLSv1.3 ClientHello.
So what's the `TLSv1.3 ClientHello`?
//todo

### Step4:The Handshake packet
RFC shows the details about handshake, note that **the handshake has 1-RTT, not the whole process which include the 
initial stages** and **the server side send two packets back, one initial and one Handshake**.
But **(the quic) handshake is different from the cryptographic handshake**, the cryptographic handshake is carried 
in Initial and Handshake packets.
```
Client                                               Server

Initial (CRYPTO)
0-RTT (*)              ---------->
                                           Initial (CRYPTO)
                                         Handshake (CRYPTO)
                       <----------                1-RTT (*)
Handshake (CRYPTO)
1-RTT (*)              ---------->
                       <----------   1-RTT (HANDSHAKE_DONE)

1-RTT                  <=========>                    1-RTT
```
The most important part in handshake packet is `CRYPTO`. It contains a TLSv1.3 handshake packet which has 
multiple handshake messages(includes Encrypted extensions, Certificate and Certificate Verify).

### Step4:Protected payload after the handshake packet of server
After the handshake sent by server, there is a `protected packet`. This `protected packet` has two IETF quic packets, 
one is `CRYPTO` and the other is `NEW_CONNECTION_ID`.
- IETF packet which contains `CRYPTO` frame:  
It uses a long header and its `CRYPTO` frame stores tls handshake protocol messages
  (includes Certificate Verify and finish).  
  
- IETF packet which contains `NEW_CONNECTION_ID`:  
This packet use a **short header** and contains 3 `NEW_CONNECTION_ID` parts.
  //todo: short header definition. 
  // why here is 3 `NEW_CONNECTION_ID`?  
  
### Step5:Initial(ack) packet from client
After the server sends the handshake packet which contains the tls handshake packets and the protected payload which 
contains the `NEW_CONNECTION_ID`, the client sends an ack back with an `Initial` packet.

### Step6:Handshake packet from client
Handshake packet contains an ack and a `CRYPTO` which contains a tls handshake `finish` data.

### Step6:Two Protected Payload from client
Two protected payload follows the handshake packet and contain an ACK and a `ENTIRE_CONNECTION_ID`. After it is done, 
all works in establishing a connection has been ended and the client is waiting to receive an ACK from the server.
Note that both of those two use the **short header**.

### Step7:Handshake packet from server
The server sends back an ACK with a **short header**. After receiving this packet, a connection has been established 
for both endpoints.

### todo:
Continue to find out how to create a stream and how connection Id, congestion flow information, etc are exchanged.