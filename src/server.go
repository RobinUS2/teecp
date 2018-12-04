package main

import (
	"flag"
	"log"
	"sync"
)

var device *string
var bpfFilter *string
var verbose *bool
var layerStr *string
var outputTcp *string
var statsPrinter *bool

const DefaultLayers = "44,45"

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
	flag.StringVar(layerStr, "layers", DefaultLayers, "layers comma separated - 44 = TCP, 45 = UDP https://github.com/google/gopacket/blob/master/layers/layertypes.go")

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
	flag.Parse()
}

func main() {
	opts := NewOpts()
	opts.Device = *device
	opts.Output = "tcp" + separator + *outputTcp
	opts.BpfFilter = *bpfFilter
	opts.Verbose = *verbose
	opts.StatsPrinter = *statsPrinter
	opts.ParseLayers(*layerStr)

	// prevent usage of non opts
	unsetLocal()

	// start
	server := NewServer(opts)
	err := server.Start()
	if err != nil {
		log.Fatalf("failed to start: %s", err)
	}
}

type Server struct {
	opts      *Opts
	forwarder *Forwarder
	mux       sync.RWMutex
}

func (server *Server) Forwarder() *Forwarder {
	server.mux.RLock()
	defer server.mux.RUnlock()
	return server.forwarder
}

func (server *Server) Start() error {
	// opts
	server.opts.AutoDiscover()
	err := server.opts.ParseOutput()
	if err != nil {
		return err
	}
	server.opts.Print()
	server.opts.Validate()

	// listener
	listener := NewRawListener(server.opts)

	//forwarder
	server.mux.Lock()
	server.forwarder = NewForwarder(server.opts)
	server.mux.Unlock()
	server.forwarder.Start()

	// inject into listener
	listener.SetForwarder(server.forwarder)

	// start
	return listener.Listen()
}

func NewServer(opts *Opts) *Server {
	return &Server{
		opts: opts,
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
