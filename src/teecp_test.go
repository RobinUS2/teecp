package main_test

import (
	teecp "../src"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const startPort = 9000
const protocol = "tcp"
const host = "127.0.0.1"

func TestNewForwarder(t *testing.T) {
	const primaryPort = startPort + 1
	const teePort = startPort + 2
	opts := teecp.NewOpts()
	opts.Device = "lo0"
	opts.BpfFilter = fmt.Sprintf("port %d and dst %s", primaryPort, host)
	opts.ParseLayers(teecp.DefaultLayers)
	opts.Output = fmt.Sprintf("tcp|%s:%d", host, teePort)
	opts.Verbose = true
	opts.StatsPrinter = true
	opts.StatsIntervalMilliseconds = 500
	controls := runTest(t, opts, primaryPort, teePort, nil)
	// await
	<-controls.shutdown
}

func TestNewForwarderSilent(t *testing.T) {
	const primaryPort = startPort + 3
	const teePort = startPort + 4
	opts := teecp.NewOpts()
	opts.Device = "lo0"
	opts.BpfFilter = fmt.Sprintf("port %d and dst %s", primaryPort, host)
	opts.ParseLayers(teecp.DefaultLayers)
	opts.Output = fmt.Sprintf("tcp|%s:%d", host, teePort)
	opts.Verbose = false
	opts.StatsPrinter = false
	opts.StatsIntervalMilliseconds = 500
	controls := runTest(t, opts, primaryPort, teePort, nil)
	// await
	<-controls.shutdown
}

func TestNewForwarderTeeDown(t *testing.T) {
	const primaryPort = startPort + 5
	const teePort = startPort + 6
	opts := teecp.NewOpts()
	opts.Device = "lo0"
	opts.BpfFilter = fmt.Sprintf("port %d and dst %s", primaryPort, host)
	opts.ParseLayers(teecp.DefaultLayers)
	opts.Output = fmt.Sprintf("tcp|%s:%d", host, teePort)
	opts.Verbose = false
	opts.StatsPrinter = false
	opts.StatsIntervalMilliseconds = 500
	counter := uint64(0)

	testOpts := &TestOpts{}
	controls := runTest(t, opts, primaryPort, teePort, testOpts)

	// inject on message hook
	testOpts.mux.Lock()
	testOpts.onMsg = func(port int, msg []byte) {
		if atomic.AddUint64(&counter, 1) == 5 {
			// in middle, close connection of tee
			if err := controls.conTee.Close(); err != nil {
				t.Error(err)
			}
		}
	}
	testOpts.mux.Unlock()

	// await a bit for the retries, but don't wait for final shutdown
	time.Sleep(1 * time.Second)
}

type TestControls struct {
	shutdown   chan bool
	conPrimary *net.TCPListener
	conTee     *net.TCPListener
	forwarder  *teecp.Server
}

type TestOpts struct {
	onMsg func(port int, msg []byte)
	mux   sync.RWMutex
}

func (testOpts *TestOpts) OnMsg() func(port int, msg []byte) {
	testOpts.mux.RLock()
	defer testOpts.mux.RUnlock()
	return testOpts.onMsg
}

func runTest(t *testing.T, opts *teecp.Opts, primaryPort int, teePort int, testOpts *TestOpts) *TestControls {
	shutdown := make(chan bool, 1)
	payload := fmt.Sprintf("test-msg-%d", time.Now().UnixNano())
	numTee := uint64(0)
	numPrimary := uint64(0)

	// primary receiver (e.g. your web server)
	conPrimary := getConnection(t, primaryPort)
	readConnection(t, conPrimary, func(msg []byte) {
		// test hook
		if testOpts != nil && testOpts.OnMsg() != nil {
			testOpts.OnMsg()(primaryPort, msg)
		}
		// validate
		str := string(msg)
		log.Printf("primary %s", str)
		if !strings.HasPrefix(str, payload) {
			return
		}
		atomic.AddUint64(&numPrimary, 1)
	})

	// tee receiver (e.g. your staging environment)
	conTee := getConnection(t, teePort)
	readConnection(t, conTee, func(msg []byte) {
		// test hook
		if testOpts != nil && testOpts.OnMsg() != nil {
			testOpts.OnMsg()(teePort, msg)
		}
		// validate
		str := string(msg)
		log.Printf("tee %s", str)
		if !strings.HasPrefix(str, payload) {
			return
		}
		atomic.AddUint64(&numTee, 1)
	})

	// tee forwarder
	forwarder := teecp.NewServer(opts)
	go func() {
		err := forwarder.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	// send to primary
	go func() {
		for i := 0; i < 10; i++ {
			addr, err := net.ResolveTCPAddr(protocol, fmt.Sprintf("%s:%d", host, primaryPort))
			if err != nil {
				t.Error(err)
			}
			con, err := net.DialTCP(protocol, nil, addr)
			if err != nil {
				t.Error(err)
			}
			_, err = con.Write([]byte(payload + fmt.Sprintf("-%d", i)))
			if err != nil {
				t.Error(err)
			}
		}
	}()

	// wait for at least 1 stats print
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)

			// done :)
			if atomic.LoadUint64(&numTee) == 10 && atomic.LoadUint64(&numPrimary) == 10 {
				shutdown <- true
			}
		}
	}()

	return &TestControls{
		shutdown:   shutdown,
		conPrimary: conPrimary,
		conTee:     conTee,
		forwarder:  forwarder,
	}
}

func getConnection(t *testing.T, port int) *net.TCPListener {
	addr, err := net.ResolveTCPAddr(protocol, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		t.Error(err)
	}
	con, err := net.ListenTCP(protocol, addr)
	if err != nil {
		t.Error(err)
	}
	return con
}

func readConnection(t *testing.T, listener *net.TCPListener, onMsg func(msg []byte)) {
	go func() {
		for {
			con, err := listener.Accept()
			if con == nil {
				break
			}
			if err != nil {
				t.Error(err)
			}
			buf := make([]byte, teecp.DefaultMaxPacketSize)
			n, err := con.Read(buf)
			if err != nil {
				t.Error(err)
			}
			buf = buf[0:n]
			log.Printf("con %v %v %s", con, buf, string(buf))
			onMsg(buf)
		}
	}()
}
