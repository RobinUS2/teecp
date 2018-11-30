package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

type Forwarder struct {
	opts  *Opts
	queue chan Payload
	stop  bool
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
}

type ForwarderInstance struct {
	forwarder *Forwarder
	conn      *net.TCPConn
	connMux   sync.RWMutex
}

func (instance *ForwarderInstance) Conn() *net.TCPConn {
	instance.connMux.Lock()
	defer instance.connMux.Unlock()
	if instance.conn == nil {
		// resolve
		tcpAddr, err := net.ResolveTCPAddr("tcp", instance.forwarder.opts.OutputTcp)
		if err != nil {
			log.Fatalf("could not resolve target %s", err)
		}
		log.Printf("tcpAddr %v", tcpAddr)

		// connect
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Printf("could not resolve target %s", err)
			return nil
		}
		instance.conn = conn
	}
	return instance.conn
}

func (forwarder *Forwarder) Run() *ForwarderInstance {
	instance := &ForwarderInstance{
		forwarder: forwarder,
	}
	for {
		payload := <-forwarder.queue
		err := forwarder.send(instance.Conn(), payload)
		if err != nil {
			log.Printf("failed to send %s (payload %x)", err, payload.data)
		}

		// stop
		if forwarder.stop {
			break
		}
	}
	return instance
}

func (forwarder *Forwarder) send(conn *net.TCPConn, payload Payload) error {
	log.Printf("sending %x", payload.data)
	n, err := conn.Write(payload.data)
	if err != nil {
		return err
	}
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
