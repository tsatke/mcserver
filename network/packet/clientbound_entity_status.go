package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundEntityStatus{}))
}

type ClientboundEntityStatus struct {
	EntityID int
	Status   int8
}

func (ClientboundEntityStatus) ID() ID       { return IDClientboundEntityStatus }
func (ClientboundEntityStatus) Name() string { return "Entity Status" }

func (c ClientboundEntityStatus) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteInt("entity id", int32(c.EntityID))
	enc.WriteByte("entity status", c.Status)

	return
}
