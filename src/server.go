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

const DefaultLayers = "44,45"

func init() {
	flag.StringVar(&device, "device", "", "device")
	flag.StringVar(&bpfFilter, "bpf", "icmp", "BPF filter (e.g. tcp port 1234)")
	flag.StringVar(&layerStr, "layers", DefaultLayers, "layers comma separated - 44 = TCP, 45 = UDP https://github.com/google/gopacket/blob/master/layers/layertypes.go")
	flag.StringVar(&outputTcp, "output-tcp", "", "TCP output configuration (e.g. your-host:1234")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")
	flag.BoolVar(&statsPrinter, "stats", true, "print stats")
	flag.Parse()
}

func main() {
	opts := NewOpts()
	opts.Device = device
	opts.Output = "tcp" + separator + outputTcp
	opts.BpfFilter = bpfFilter
	opts.Verbose = verbose
	opts.StatsPrinter = statsPrinter
	opts.ParseLayers(layerStr)

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
	opts *Opts
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
	if server.opts.Verbose {
		verbose = true
	}

	// listener
	listener := NewRawListener(server.opts)

	//forwarder
	forwarder := NewForwarder(server.opts)
	forwarder.Start()

	// inject into listener
	listener.SetForwarder(forwarder)

	// start
	return listener.Listen()
}

func NewServer(opts *Opts) *Server {
	return &Server{
		opts: opts,
	}
}

func unsetLocal() {
	device = ""
	bpfFilter = ""
	layerStr = ""
	outputTcp = ""
}
