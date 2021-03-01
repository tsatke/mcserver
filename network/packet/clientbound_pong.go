package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseStatus, reflect.TypeOf(ClientboundPong{}))
}

// ClientboundPong is the response to the ServerboundPing packet.
type ClientboundPong struct {
	// Payload must be the same value that the client sent in
	// the ServerboundPing packet. However, the client does
	// NOT use this for latency computation.
	Payload int64
}

// ID returns the constant packet ID.
func (ClientboundPong) ID() ID { return IDClientboundPong }

// Name returns the constant packet name.
func (ClientboundPong) Name() string { return "Pong" }

// EncodeInto writes this packet into the given writer.
func (c ClientboundPong) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteLong("payload", c.Payload)

	return
}
