package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StateStatus, reflect.TypeOf(ClientboundPong{}))
}

type ClientboundPong struct {
	Payload int64
}

func (ClientboundPong) ID() ID       { return IDClientboundPong }
func (ClientboundPong) Name() string { return "Pong" }

func (c ClientboundPong) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeLong("payload", c.Payload)

	return
}
