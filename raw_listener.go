package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
)

type RawListener struct {
	opts *Opts
}

func (listener *RawListener) Listen() error {
	if handle, err := pcap.OpenLive(listener.opts.Device, listener.opts.MaxPacketSize, true, pcap.BlockForever); err != nil {
		return err
	} else if err := handle.SetBPFFilter(bpfFilter); err != nil { // optional
		return err
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			go listener.handlePacket(packet)
		}
	}
	return nil
}

func (listener *RawListener) handlePacket(packet gopacket.Packet) {
	if verbose {
		log.Printf("recv packet %+v", packet) // Do something with a packet here.
	}
	for _, layer := range packet.Layers() {
		if verbose {
			log.Printf("  layer type=%d payload=%x (%s) content=%x", layer.LayerType(), layer.LayerPayload(), string(layer.LayerPayload()), layer.LayerContents())
		}
		if listener.opts.layers != nil && !listener.opts.layers[int(layer.LayerType())] {
			// not in expected layers
			continue
		}
		log.Printf("%s", strings.TrimSpace(string(layer.LayerPayload())))
	}
}

func NewRawListener(opts *Opts) *RawListener {
	return &RawListener{
		opts: opts,
	}
}
