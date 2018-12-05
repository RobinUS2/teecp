package forwarder

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
			log.Printf("stats %+v", forwarder.Stats())
		}
	}()
}

type Stats struct {
	BytesForwarded       uint64
	PacketsForwarded     uint64
	PacketAttemptsFailed uint64
	PacketsFailed        uint64
}

func (forwarder *Forwarder) Stats() Stats {
	return Stats{
		BytesForwarded:       atomic.LoadUint64(&forwarder.bytesForwarded),
		PacketsForwarded:     atomic.LoadUint64(&forwarder.packetsForwarded),
		PacketAttemptsFailed: atomic.LoadUint64(&forwarder.packetAttemptsFailed),
		PacketsFailed:        atomic.LoadUint64(&forwarder.packetsFailed),
	}
}
