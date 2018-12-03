package main

import (
	"github.com/pkg/errors"
	"log"
	"net"
	"strings"
	"sync/atomic"
)

type Forwarder struct {
	opts  *Opts
	queue chan Payload
	stop  bool

	packetsForwarded uint64
	packetsFailed    uint64
	bytesForwarded   uint64
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
}

func (instance *ForwarderInstance) Conn() *net.TCPConn {
	// @todo cache name resolution?
	// resolve
	tcpAddr, err := net.ResolveTCPAddr("tcp", instance.forwarder.opts.OutputTcp)
	if err != nil {
		log.Fatalf("could not resolve target %s", err)
	}
	log.Printf("tcpAddr %v", tcpAddr)

	// connect
	// @todo option to enable keepalive, reuse conections
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("could not resolve target %s", err)
		return nil
	}
	return conn
}

func (forwarder *Forwarder) Run() *ForwarderInstance {
	instance := &ForwarderInstance{
		forwarder: forwarder,
	}
	for {
		payload := <-forwarder.queue
		err := forwarder.send(instance.Conn(), payload)
		if err != nil {
			atomic.AddUint64(&forwarder.packetsFailed, 1)
			log.Printf("failed to send %s (payload %x)", err, payload.data)
			continue
		}

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
	return nil
}

func NewForwarder(opts *Opts) *Forwarder {
	forwarder := &Forwarder{
		opts:  opts,
		queue: make(chan Payload, opts.QueueSize),
	}
	return forwarder
}
