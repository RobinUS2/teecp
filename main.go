package main

import (
	"flag"
	"github.com/RobinUS2/teecp/forwarder"
	"log"
)

var device *string
var bpfFilter *string
var verbose *bool
var layerStr *string
var outputTcp *string
var statsPrinter *bool
var statsHealthCheck *bool
var prefixHeader *bool

func init() {
	{
		str := ""
		device = &str
	}
	flag.StringVar(device, "device", "", "device")

	{
		str := ""
		bpfFilter = &str
	}
	flag.StringVar(bpfFilter, "bpf", "icmp", "BPF filter (e.g. tcp port 1234)")

	{
		str := ""
		layerStr = &str
	}
	flag.StringVar(layerStr, "layers", forwarder.DefaultLayers, "layers comma separated - 44 = TCP, 45 = UDP https://github.com/google/gopacket/blob/master/layers/layertypes.go")

	{
		str := ""
		outputTcp = &str
	}
	flag.StringVar(outputTcp, "output-tcp", "", "TCP output configuration (e.g. your-host:1234")

	{
		b := false
		verbose = &b
	}
	flag.BoolVar(verbose, "verbose", false, "verbose logging")

	{
		b := false
		statsPrinter = &b
	}
	flag.BoolVar(statsPrinter, "stats", true, "print stats")

	{
		b := false
		statsHealthCheck = &b
	}
	flag.BoolVar(statsHealthCheck, "health-check", true, "health check")

	{
		b := false
		prefixHeader = &b
	}
	flag.BoolVar(prefixHeader, "prefix-header", false, "prefix header (include payload size, src, src port, dst, dst port)")

	// parse
	flag.Parse()
}

func main() {
	opts := forwarder.NewOpts()
	opts.Device = *device
	opts.Output = "tcp" + forwarder.ConfSeparator + *outputTcp
	opts.BpfFilter = *bpfFilter
	opts.Verbose = *verbose
	opts.StatsPrinter = *statsPrinter
	opts.StatsHealthCheck = *statsHealthCheck
	opts.PrefixHeader = *prefixHeader
	opts.ParseLayers(*layerStr)

	// prevent usage of non opts
	unsetLocal()

	// start
	server := forwarder.NewServer(opts)
	err := server.Start()
	if err != nil {
		log.Fatalf("failed to start: %s", err)
	}
}

func unsetLocal() {
	device = nil
	bpfFilter = nil
	layerStr = nil
	outputTcp = nil
	verbose = nil
	statsPrinter = nil
}
