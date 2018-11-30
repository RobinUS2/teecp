# teecp
TCP tee implementation (Linux, Mac OS X, Windows) - duplicate TCP packets

## Purpose
Run outside of the regular traffic flow, listen to the TCP packets and duplicate them
to other sources with minimal impact. This means it needs no changes in the
existing applications that run their. For example you run a process on port 
1234 TCP. You can start the copying process teecp that monitors that TCP port
and copies all the individual packets to another location.

## How does it work?
It relies on the promiscuous mode ethernet sniffing mode which is also used by 
tools like [WireShark](https://www.wireshark.org/), [WinPcap](https://www.winpcap.org/), etc.

It is built around Google's [gopacket](https://github.com/google/gopacket) library and written in GoLang.

< Diagram here >

## Related projects
- GoReplay HTTP(S) https://github.com/buger/goreplay
- TCPCopy https://github.com/session-replay-tools/tcpcopy
- Duplicator https://github.com/agnoster/duplicator
- IPTables `iptables -t mangle -A POSTROUTING -p tcp --dport 1234 -j TEE --gateway IP_HOST_B`