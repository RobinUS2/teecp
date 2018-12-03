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
tools like [WireShark](https://www.wireshark.org/), 
[WinPcap](https://www.winpcap.org/), 
[tcpdump](https://www.tcpdump.org/), etc.

It is built around Google's [gopacket](https://github.com/google/gopacket) library and written in GoLang.

< Diagram here >

By default the payload of the packet is forwarded (without the encapsulating layers). 
It is however possible to forward the entire packet payload without any filters.

## How to run?
The below will listen on interface `lo0`, filter traffic on port 1234, log 
all details (very verbose, turn off in production), and copy it's packet payloads 
(by default TCP & UDP) towards localhost port 8080.
```
./teecp --device=lo0 --bpf='port 1234' --verbose=true --output-tcp 'localhost:8080'
```

The `--bpf` flag can handle [Berkeley Packet Filter](https://en.wikipedia.org/wiki/Berkeley_Packet_Filter) syntax. 

A handful of examples:

| Example                 | Syntax                                         |
|-------------------------|------------------------------------------------|
| TCP only                | tcp                                            |
| TCP for a specific port | tcp port 1234                                  |
| + specific source       | tcp port 1234 and src 1.2.3.4                  |
| + specific destination  | tcp port 1234 and src 1.2.3.4 and dst 10.0.0.1 | 

## Build & test
The application relies upon libpcap (for Windows, [download developer pack](https://www.winpcap.org/devel.htm)) and [GoLang](https://golang.org/doc/install). 

OS X
```
brew install libpcap
```

Putting it all together
```
go vet . && go fmt . && go test -v . && go build . && ./teecp --device=lo0 --bpf='port 1234' --verbose=true --output-tcp "test.com:123"
```

## Used by
- [Route42](https://route42.nl/)
- open a PR and add YourCompany!

## Related projects
- GoReplay HTTP(S) https://github.com/buger/goreplay
- TCPCopy https://github.com/session-replay-tools/tcpcopy
- Duplicator https://github.com/agnoster/duplicator
- IPTables `iptables -t mangle -A POSTROUTING -p tcp --dport 1234 -j TEE --gateway IP_HOST_B`