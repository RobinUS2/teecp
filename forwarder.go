package main

import (
	"log"
	"net"
	"strings"
)

type Forwarder struct {
	opts  *Opts
	queue chan Payload
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
			tcpAddr, err := net.ResolveTCPAddr("tcp", forwarder.opts.OutputTcp)
			if err != nil {
				log.Fatalf("could not resolve target %s", err)
			}
			log.Printf("tcpAddr %v", tcpAddr)
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				log.Fatalf("could not resolve target %s", err)
			}
			for {
				payload := <-forwarder.queue
				err := forwarder.send(conn, payload)
				if err != nil {
					log.Printf("failed to send %s (payload %x)", err, payload.data)
				}
			}
		}()
	}
}
func (forwarder *Forwarder) send(conn *net.TCPConn, payload Payload) error {
	log.Printf("sending %x", payload.data)
	_, err := conn.Write(payload.data)
	if err != nil {
		return err
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
