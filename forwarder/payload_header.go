package forwarder

import (
	"bytes"
	"encoding/binary"
	"net"
)

type PayloadHeader struct {
	PayloadSize uint32 // 4 bytes : number of bytes of the actual payload (excludes the fixed length payload fields below)
	Src         []byte // 16 bytes
	SrcPort     uint32 // 4 bytes
	Dst         []byte // 16 bytes
	DstPort     uint32 // 4 bytes
}

func NewPayloadHeader(payloadSize int, src net.IP, srcPort int, dst net.IP, dstPort int) PayloadHeader {
	return PayloadHeader{
		PayloadSize: uint32(payloadSize),
		Src:         toBytes(src),
		SrcPort:     uint32(srcPort),
		Dst:         toBytes(dst),
		DstPort:     uint32(dstPort),
	}
}

func (header PayloadHeader) SrcIP() net.IP {
	return net.IP(header.Src)
}

func (header PayloadHeader) DstIP() net.IP {
	return net.IP(header.Dst)
}

func (header PayloadHeader) Bytes() []byte {
	if header.PayloadSize < 1 {
		panic("missing payload size")
	}
	buf := bytes.Buffer{}

	// payload length
	writeUint32(header.PayloadSize, &buf)

	// src
	buf.Write(header.Src)

	// src port
	writeUint32(header.SrcPort, &buf)

	// dst
	buf.Write(header.Dst)

	// dst port
	writeUint32(header.DstPort, &buf)

	return buf.Bytes()
}

func PayloadHeaderFromBytes(b []byte) PayloadHeader {
	buf := bytes.NewBuffer(b)
	payloadSize := readUint32(buf)

	// src
	src := make([]byte, 16)
	buf.Read(src)

	// src port
	srcPort := readUint32(buf)

	// dst
	dst := make([]byte, 16)
	buf.Read(dst)

	// dst port
	dstPort := readUint32(buf)

	// dst port
	return PayloadHeader{
		PayloadSize: payloadSize,
		Src:         src,
		SrcPort:     srcPort,
		Dst:         dst,
		DstPort:     dstPort,
	}
}

func toBytes(ip net.IP) []byte {
	return ip.To16()
}

func writeUint32(v uint32, buf *bytes.Buffer) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, v)
	buf.Write(bs)
}

func readUint32(buf *bytes.Buffer) uint32 {
	bs := make([]byte, 4)
	buf.Read(bs)
	return binary.BigEndian.Uint32(bs)
}
