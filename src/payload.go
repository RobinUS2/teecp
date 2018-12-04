package main

import "bytes"

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
