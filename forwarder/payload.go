package forwarder

import (
	"bytes"
	"github.com/google/gopacket"
	"net"
)

type Payload struct {
	data   []byte
	header PayloadHeader
}

func (payload Payload) Bytes() []byte {
	if payload.header.PayloadSize < 1 {
		// no header
		return payload.data
	}

	// join header
	buf := bytes.Buffer{}
	buf.Write(payload.header.Bytes())
	// payload
	buf.Write(payload.data)
	return buf.Bytes()
}

func NewPayload(payload []byte) Payload {
	return Payload{
		data: payload,
	}
}

func (payload *Payload) SetHeader(header PayloadHeader) {
	payload.header = header
}

func PayloadHeaderFromPacket(payloadLen int, packet gopacket.Packet) PayloadHeader {
	return NewPayloadHeader(payloadLen, net.ParseIP(packet.NetworkLayer().NetworkFlow().Src().String()), strToInt(packet.TransportLayer().TransportFlow().Src().String()), net.ParseIP(packet.NetworkLayer().NetworkFlow().Dst().String()), strToInt(packet.TransportLayer().TransportFlow().Dst().String()))
}
