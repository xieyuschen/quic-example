# quic-example
This reporsitory aims to help user get familiar with the quic and use quic in the real app. It develops some demo based on the 
[quic-go](https://github.com/lucas-clemente/quic-go) library and expands the example showed in the reporsitory.

## Content
Keep updating...
|Example|Description|
|:--|:--|
|[Echo demo](1.echo/README.md)|Echo demo comes from the quic-go/example and it aims to help users to learn it better and more quickly|
|[Echo-cli](2.echo-cli/README.md)|The echo-cli provides the seperated server and client to make a demo for echo case. It mainly focuses on some details when using quic and make some discussions about quic.|
|[Capture packets during echo](3.packets-during-echo/README.md)|This part uses wireshark to capture the echo packets and analyzes the whole processes.|
|[Handling multiple servers](4.multiple-streams/README.md)|Implement quic server which can handle multiple streams.|
|[Set option on listener](5.setopt-on-listener/README.md)|How to set options on listener as the net ListenConfig.|
|[Framework implementation details](etc.md)|This part focuses on some mechanism during building quic-go|

# Contribute
Feel free to create PR to enhance this project.
