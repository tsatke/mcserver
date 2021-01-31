package packet

import (
	"io"
	"reflect"

	"github.com/tsatke/mcserver/game/voxel"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundUpdateViewPosition{}))
}

type ClientboundUpdateViewPosition struct {
	Chunk voxel.V2
}

func (ClientboundUpdateViewPosition) ID() ID       { return IDClientboundUpdateViewPosition }
func (ClientboundUpdateViewPosition) Name() string { return "Update View Position" }

func (c ClientboundUpdateViewPosition) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteVarInt("chunk x", c.Chunk.X)
	enc.WriteVarInt("chunk z", c.Chunk.Z)

	return
}
