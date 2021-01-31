package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseStatus, reflect.TypeOf(ClientboundPong{}))
}

type ClientboundPong struct {
	Payload int64
}

func (ClientboundPong) ID() ID       { return IDClientboundPong }
func (ClientboundPong) Name() string { return "Pong" }

func (c ClientboundPong) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteLong("payload", c.Payload)

	return
}
