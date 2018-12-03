package main

import (
	"flag"
	"log"
)

var device string
var bpfFilter string
var verbose bool
var layerStr string
var outputTcp string
var statsPrinter bool

func init() {
	flag.StringVar(&device, "device", "", "device")
	flag.StringVar(&bpfFilter, "bpf", "icmp", "BPF filter (e.g. tcp port 1234)")
	flag.StringVar(&layerStr, "layers", "44,45", "layers comma separated - 44 = TCP, 45 = UDP https://github.com/google/gopacket/blob/master/layers/layertypes.go")
	flag.StringVar(&outputTcp, "output-tcp", "", "TCP output configuration (e.g. your-host:1234")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")
	flag.BoolVar(&statsPrinter, "stats", true, "print stats")
	flag.Parse()
}

func main() {
	opts := NewOpts()
	opts.Device = device
	opts.OutputTcp = outputTcp
	opts.BpfFilter = bpfFilter
	opts.StatsPrinter = statsPrinter
	opts.ParseLayers(layerStr)
	opts.AutoDiscover()
	opts.Print()
	opts.Validate()

	// prevent usage of non opts
	unsetLocal()

	// listener
	listener := NewRawListener(opts)

	//forwarder
	forwarder := NewForwarder(opts)
	forwarder.Start()

	// inject into listener
	listener.SetForwarder(forwarder)

	// start
	err := listener.Listen()
	if err != nil {
		log.Fatalf("failed to listen %s", err)
	}
}

func unsetLocal() {
	device = ""
	bpfFilter = ""
	layerStr = ""
	outputTcp = ""
}
