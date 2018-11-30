package main

import (
	"flag"
	"log"
)

var device string
var bpfFilter string
var verbose bool
var layerStr string
var layers map[int]bool

func init() {
	flag.StringVar(&device, "device", "", "device")
	flag.StringVar(&bpfFilter, "bpf", "icmp", "BPF filter")
	flag.StringVar(&layerStr, "layers", "44,45", "layers comma separated - https://github.com/google/gopacket/blob/master/layers/layertypes.go")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")
	flag.Parse()
}

func main() {
	opts := NewOpts()
	opts.Device = device
	opts.ParseLayers(layerStr)
	opts.AutoDiscover()
	opts.Print()

	listener := NewRawListener(opts)
	err := listener.Listen()
	if err != nil {
		log.Fatalf("failed to listen %s", err)
	}
}
