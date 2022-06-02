# 3.packets-during-echo
This section aims to learn about the packets during the echo example. First, the tools used here will be discussed and 
then will show the packet details.

## Tools discussion
There are many popular tools to monitor the TCP network such as `wireshark`, `tcpdump`, `netstat` and so on. As quic 
underlies the udp protocol, what we should do is to find a useful tool on UDP.
- wireshark.  
Wireshark is an open-source software with a good GUI.  
  
- netstat.  
Netstat integrates with the Linux administrator and can get manual by `man netstat` on ubuntu(can also visit the 
  [website](https://linux.die.net/man/8/netstat). So I prefer the 
  netstat to analyze the udp packets.
  
## Analysis
