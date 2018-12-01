package main

type Payload struct {
	data []byte
}

func NewPayload(payload []byte) Payload {
	return Payload{
		data: payload,
	}
}
