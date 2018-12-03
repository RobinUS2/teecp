package main

import (
	"log"
	"sync/atomic"
	"time"
)

const DefaultStatsIntervalMilliseconds = 10 * 1000

func (forwarder *Forwarder) printStats() {
	if !forwarder.opts.StatsPrinter {
		return
	}
	interval := DefaultStatsIntervalMilliseconds
	if forwarder.opts.StatsIntervalMilliseconds > 0 {
		interval = forwarder.opts.StatsIntervalMilliseconds
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			log.Printf("forwarded %d bytes, %d packets, %d failed attempts, %d failed packets", atomic.LoadUint64(&forwarder.bytesForwarded), atomic.LoadUint64(&forwarder.packetsForwarded), atomic.LoadUint64(&forwarder.packetAttemptsFailed), atomic.LoadUint64(&forwarder.packetsFailed))
		}
	}()
}
