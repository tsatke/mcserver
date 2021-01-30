package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundHeldItemChange{}))
}

type ClientboundHeldItemChange struct {
	Slot int8
}

func (ClientboundHeldItemChange) ID() ID       { return IDClientboundHeldItemChange }
func (ClientboundHeldItemChange) Name() string { return "Held item change" }

func (c ClientboundHeldItemChange) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeByte("slot", c.Slot)

	return
}
