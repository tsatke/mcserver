package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseStatus, reflect.TypeOf(ServerboundPing{}))
}

// ServerboundPing will be sent by the client in the status phase, indicating that he
// expects a ClientboundPong packet.
//
// Note: it seems that the vanilla client does not use the payload to measure the ping,
// but instead the time between the ServerboundPing and the ClientboundPong packet.
type ServerboundPing struct {
	// Payload is the payload that is sent by the client, which is a timestamp if sent by the
	// vanilla client. Don't rely on this though. In any case, this is to be sent back to the
	// client in a pong message unmodified.
	Payload int64
}

// ID returns the constant packet ID.
func (ServerboundPing) ID() ID { return IDServerboundPing }

// Name returns the constant packet name.
func (ServerboundPing) Name() string { return "Ping" }

// DecodeFrom will fill this struct with values read from the given reader.
func (s *ServerboundPing) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.Payload = dec.ReadLong("payload")

	return
}
