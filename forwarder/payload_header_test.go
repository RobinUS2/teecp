package forwarder_test

import (
	"github.com/RobinUS2/teecp/forwarder"
	"log"
	"net"
	"testing"
)

func TestNewPayloadHeader(t *testing.T) {
	header := forwarder.NewPayloadHeader(len([]byte("test")), net.ParseIP("127.0.0.1"), 1234, net.ParseIP("10.0.0.4"), 9999)
	log.Printf("%+v", header)
	b := header.Bytes()
	log.Printf("%d %x", len(b), b)

	{
		readHeader := forwarder.PayloadHeaderFromBytes(b)
		log.Printf("%+v", readHeader)
		src := readHeader.SrcIP()
		if src.String() != "127.0.0.1" {
			t.Error(src.String())
		}

		if readHeader.SrcPort != 1234 {
			t.Error(readHeader.SrcPort)
		}

		dst := readHeader.DstIP()
		if dst.String() != "10.0.0.4" {
			t.Error(dst.String())
		}

		if readHeader.DstPort != 9999 {
			t.Error(readHeader.DstPort)
		}
	}
}
