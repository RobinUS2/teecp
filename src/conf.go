package main

import (
	"encoding/json"
	"github.com/google/gopacket/pcap"
	"log"
	"runtime"
	"strconv"
	"strings"
)

type Opts struct {
	Device        string
	BpfFilter     string
	OutputTcp     string
	MaxPacketSize int32
	QueueSize     int
	NumForwarders int
	StatsPrinter  bool

	layers map[int]bool
}

func NewOpts() *Opts {
	return &Opts{
		MaxPacketSize: 65536,
		QueueSize:     1000,
		NumForwarders: runtime.NumCPU(),
	}
}

func (opts *Opts) Print() {
	j, _ := json.Marshal(opts)
	log.Printf("config %s", string(j))
}

func (opts *Opts) Validate() {
	// @todo check that it's valid configuration
}

func (opts *Opts) ParseLayers(layerStr string) {
	opts.layers = make(map[int]bool)
	for _, str := range strings.Split(layerStr, ",") {
		i, _ := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
		if i > 0 {
			opts.layers[int(i)] = true
		}
	}
	if len(opts.layers) < 1 {
		// no filter
		opts.layers = nil
	}
}

func (opts *Opts) AutoDiscover() {
	if len(strings.TrimSpace(opts.Device)) > 0 {
		// already configured
		return
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print device information
	if verbose {
		log.Println("Devices found:")
		for _, device := range devices {
			log.Printf("Name: %s", device.Name)
			log.Printf("Description: %s", device.Description)
			for _, address := range device.Addresses {
				log.Printf("- IP address: %s", address.IP.String())
				log.Printf("- Subnet mask: %s", address.Netmask.String())
			}
		}
	}

	opts.Device = devices[0].Name
}
