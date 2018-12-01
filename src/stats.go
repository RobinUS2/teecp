package main

import (
	"log"
	"sync/atomic"
	"time"
)

func (forwarder *Forwarder) printStats() {
	if !forwarder.opts.StatsPrinter {
		return
	}
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for _ = range ticker.C {
			log.Printf("forwarded %d bytes, %d packets, %d failed packets", atomic.LoadUint64(&forwarder.bytesForwarded), atomic.LoadUint64(&forwarder.packetsForwarded), atomic.LoadUint64(&forwarder.packetsFailed))
		}
	}()
}
