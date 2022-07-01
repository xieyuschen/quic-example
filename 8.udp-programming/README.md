# Udp programming
When it comes to udp, there is a question: `why there is no accept function for udp connection`?

Look at the document of net library about listener, it says:
```text
A Listener is a generic network listener for stream-oriented protocols.

Multiple goroutines may invoke methods on a Listener simultaneously.
```
The most important one is that it works **stream-oriented** protocol. 
*I assume that listener deals with the connection establishment, after finishing it creates a new connection*.

Then, searching on stackoverflow, and it provides a visible way about accept():  
It has everything to do with the fact that TCP has connections. The whole point of accept() is to 
get a stream that contains data from only that connection, in sequence, without loss or duplication. 

//TODO: so where does the client send its data really? listener or the socket endpoint?

As UDP is a connectionless protocol, you don't need it. You get remote address information with every incoming UDP 
datagram, so you know who it's from, so you don't need a per-connection socket to tell you. 

//TODO: why connection based protocol tcp need to get the remote address information from the connection?
