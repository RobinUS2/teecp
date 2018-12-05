package forwarder

import (
	"sync"
)

type Server struct {
	opts      *Opts
	forwarder *Forwarder
	mux       sync.RWMutex
}

func (server *Server) Forwarder() *Forwarder {
	server.mux.RLock()
	defer server.mux.RUnlock()
	return server.forwarder
}

func (server *Server) Stats() Stats {
	return server.Forwarder().Stats()
}

func (server *Server) Start() error {
	// opts
	server.opts.AutoDiscover()
	err := server.opts.ParseOutput()
	if err != nil {
		return err
	}
	server.opts.Print()
	server.opts.Validate()

	// listener
	listener := NewRawListener(server.opts)

	//forwarder
	server.mux.Lock()
	server.forwarder = NewForwarder(server.opts)
	server.mux.Unlock()
	server.forwarder.Start()

	// inject into listener
	listener.SetForwarder(server.forwarder)

	// start
	return listener.Listen()
}

func NewServer(opts *Opts) *Server {
	return &Server{
		opts: opts,
	}
}
