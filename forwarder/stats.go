package forwarder

import (
	"log"
	"os"
	"sync/atomic"
	"time"
)

const DefaultStatsIntervalMilliseconds = 10 * 1000

func (forwarder *Forwarder) printStats() {
	interval := DefaultStatsIntervalMilliseconds
	if forwarder.opts.StatsIntervalMilliseconds > 0 {
		interval = forwarder.opts.StatsIntervalMilliseconds
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		timeSinceLastForward := time.Now()
		var previousStats Stats
		for _ = range ticker.C {
			// obtain stats
			stats := forwarder.Stats()

			// print
			if forwarder.opts.StatsPrinter {
				log.Printf("stats %+v", stats)
			}

			// health
			if forwarder.opts.StatsHealthCheck {
				numForwarded := stats.PacketsForwarded - previousStats.PacketsForwarded
				if numForwarded > 0 {
					// increase means we forwarded
					timeSinceLastForward = time.Now()
				}

				// long time no forward?
				secondsSinceForward := int(time.Now().Sub(timeSinceLastForward).Seconds())
				if secondsSinceForward > 30*60 {
					log.Printf("time since last forward is %d seconds ago, exiting", secondsSinceForward)
					os.Exit(1)
				}
			}

			// assign previous
			previousStats = stats
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
