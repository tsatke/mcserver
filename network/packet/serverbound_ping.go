package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseStatus, reflect.TypeOf(ServerboundPing{}))
}

type ServerboundPing struct {
	// Payload is the payload that is sent by the client, which is a timestamp if sent by the
	// vanilla client. Don't rely on this though. In any case, this is to be sent back to the
	// client in a pong message unmodified.
	Payload int64
}

func (ServerboundPing) ID() ID       { return IDServerboundPing }
func (ServerboundPing) Name() string { return "Ping" }

func (s *ServerboundPing) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.Payload = dec.ReadLong("payload")

	return
}
