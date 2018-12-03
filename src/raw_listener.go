package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
)

type RawListener struct {
	opts      *Opts
	forwarder *Forwarder
}

func (listener *RawListener) Forwarder() *Forwarder {
	return listener.forwarder
}

func (listener *RawListener) SetForwarder(forwarder *Forwarder) {
	listener.forwarder = forwarder
}

func (listener *RawListener) Listen() error {
	if handle, err := pcap.OpenLive(listener.opts.Device, listener.opts.MaxPacketSize, true, pcap.BlockForever); err != nil {
		return err
	} else if err := handle.SetBPFFilter(listener.opts.BpfFilter); err != nil { // optional
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

	// specific layer?
	if listener.opts.layers != nil {
		for _, layer := range packet.Layers() {
			if verbose {
				log.Printf("  layer type=%d payload=%x (%s) content=%x", layer.LayerType(), layer.LayerPayload(), string(layer.LayerPayload()), layer.LayerContents())
			}
			if listener.opts.layers != nil && !listener.opts.layers[int(layer.LayerType())] {
				// not in expected layers
				continue
			}
			// just this layer's payload
			listener.Forwarder().Queue(NewPayload(layer.LayerPayload()))
		}
	} else {
		// full packet
		listener.Forwarder().Queue(NewPayload(packet.Data()))
	}
}

func NewRawListener(opts *Opts) *RawListener {
	return &RawListener{
		opts: opts,
	}
}
