# 3.packets-during-echo
This section aims to learn about the packets during the echo example. First, the tools used here will be discussed and 
then will show the packet details.

## Tools discussion
There are many popular tools to monitor the network such as `wireshark`, `tcpdump`, `netstat` and so on. As quic 
underlies the udp protocol, what we should do is to find a useful tool on UDP.
### Linux networking tools
- ss:  
  ss command is a tool that is used for displaying network socket related information on a Linux system.

- tcpdump:
  Tcpdump is a command line utility that allows you to capture and analyze network traffic going through your system.
  It is often used to help troubleshoot network issues, as well as a security tool.

- nmap:
  Nmap is short for Network Mapper. It is an open-source Linux cmd-line tool that is used to scan IPs 
  & ports in a nw & to detect installed apps. Nmap allows nw admins to find which devices r running 
  on their nw, discover open ports & services, and detect vulnerabilities.

- dig:
  Dig (Domain Information Groper) is a powerful cmd-line tool for querying DNS name servers.
  It allows you to query info abt various DNS records, including host addresses, mail exchanges, 
  & name servers. A most common tool among sysadmins for troubleshooting DNS problems.
  
Linux provides many network tools, but if we want to take an eye on the packets the `tcpdump`.
## Analysis
If we want to use tcpdump to see the packets, the key pair and certificate should be fixed as we need it to decrypt the 
packets. As a result, I choose the files under `/cert` files:
- See the certificate details:
```shell
quic-example/cert$ openssl x509 -in ca.pem -noout -text
```

- Tcpdump to see the udp packets:
```shell
sudo tcpdump -nnvXSs 0 -i lo port 4242 -xx -tt
```
Ignoring `-X` option can avoid converting hex to text and it helps decrypting.