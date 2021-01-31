package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(StatePlay, reflect.TypeOf(ClientboundHeldItemChange{}))
}

type ClientboundHeldItemChange struct {
	Slot int8
}

func (ClientboundHeldItemChange) ID() ID       { return IDClientboundHeldItemChange }
func (ClientboundHeldItemChange) Name() string { return "Held item change" }

func (c ClientboundHeldItemChange) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteByte("slot", c.Slot)

	return
}
