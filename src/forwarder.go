package main

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Forwarder struct {
	opts  *Opts
	queue chan Payload
	stop  bool

	packetsForwarded     uint64
	packetsFailed        uint64
	packetAttemptsFailed uint64
	bytesForwarded       uint64
}

func (forwarder *Forwarder) Queue(payload Payload) {
	if payload.data == nil || len(payload.data) < 1 {
		return
	}
	if verbose {
		log.Printf("queued %x %s", payload.data, strings.TrimSpace(string(payload.data)))
	}
	forwarder.queue <- payload
}

func (forwarder *Forwarder) Start() {
	for i := 0; i < forwarder.opts.NumForwarders; i++ {
		go func() {
			forwarder.Run()
		}()
	}

	forwarder.printStats()
}

type ForwarderInstance struct {
	forwarder *Forwarder

	conn    *net.TCPConn
	connMux sync.RWMutex
}

func (instance *ForwarderInstance) Conn() *net.TCPConn {
	if instance.forwarder.opts.outputKeepAlive {
		// keep alive
		instance.connMux.RLock()
		v := instance.conn
		instance.connMux.RUnlock()
		if v != nil {
			// already existing
			return v
		}
	}

	// @todo cache name resolution?
	// @todo UDP support
	// resolve
	tcpAddr, err := net.ResolveTCPAddr("tcp", instance.forwarder.opts.outputAddress)
	if err != nil {
		log.Printf("could not resolve target %s", err)
		return nil
	}
	if verbose {
		log.Printf("tcpAddr %v", tcpAddr)
	}

	// connect
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("could not resolve target %s", err)
		return nil
	}

	// keep?
	if instance.forwarder.opts.outputKeepAlive {
		instance.connMux.Lock()
		instance.conn = conn
		instance.connMux.Unlock()
	}

	return conn
}

func (forwarder *Forwarder) Run() *ForwarderInstance {
	instance := &ForwarderInstance{
		forwarder: forwarder,
	}
	for {
		payload := <-forwarder.queue
		instance.handlePayload(payload)

		// stop
		if forwarder.stop {
			break
		}
	}
	return instance
}

func (forwarder *Forwarder) send(conn *net.TCPConn, payload Payload) error {
	if conn == nil {
		return errors.New("no connection")
	}
	log.Printf("sending %x", payload.data)
	// @todo option to prefix with length of data => byte bs := make([]byte, 4) binary.LittleEndian.PutUint32(bs, 31415926)
	n, err := conn.Write(payload.data)
	if err != nil {
		return err
	}
	// stats
	atomic.AddUint64(&forwarder.bytesForwarded, uint64(n))
	atomic.AddUint64(&forwarder.packetsForwarded, 1)

	if verbose {
		log.Printf("sent %d bytes to %s", n, conn.RemoteAddr().String())
	}

	// close?
	if forwarder.opts.outputKeepAlive == false {
		// close without keepalive
		if err := conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (instance *ForwarderInstance) handlePayload(payload Payload) {
	// retry support
	var err error
	for i := 0; i < instance.forwarder.opts.MaxRetries; i++ {
		// attempt
		err = instance.forwarder.send(instance.Conn(), payload)
		if err != nil {
			if verbose {
				log.Printf("failed attempt to send %s (payload %x)", err, payload.data)
			}
			atomic.AddUint64(&instance.forwarder.packetAttemptsFailed, 1)

			// reset connection on error
			instance.resetConnection()

			// sleep a bit
			time.Sleep(time.Duration(100+(i*i*1000)) * time.Millisecond)

			// try again
			continue
		}
		break
	}
	// fatal, after retries still not sent
	if err != nil {
		log.Printf("failed to send %s (payload %x)", err, payload.data)
		atomic.AddUint64(&instance.forwarder.packetsFailed, 1)
	}
}

func (instance *ForwarderInstance) resetConnection() {
	if instance.forwarder.opts.outputKeepAlive {
		instance.connMux.Lock()
		instance.conn = nil
		instance.connMux.Unlock()
	}
}

func NewForwarder(opts *Opts) *Forwarder {
	forwarder := &Forwarder{
		opts:  opts,
		queue: make(chan Payload, opts.QueueSize),
	}
	return forwarder
}
