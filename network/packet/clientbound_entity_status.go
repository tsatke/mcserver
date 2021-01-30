package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundEntityStatus{}))
}

type ClientboundEntityStatus struct {
	EntityID int
	Status   int8
}

func (ClientboundEntityStatus) ID() ID       { return IDClientboundEntityStatus }
func (ClientboundEntityStatus) Name() string { return "Entity Status" }

func (c ClientboundEntityStatus) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeInt("entity id", int32(c.EntityID))
	enc.writeByte("entity status", c.Status)

	return
}
