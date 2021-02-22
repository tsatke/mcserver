package packet

import (
	"io"
	"reflect"

	"github.com/tsatke/mcserver/game/voxel"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundUpdateViewPosition{}))
}

// ClientboundUpdateViewPosition is sent by the server to tell the client which chunks he should keep
// loaded. This must be sent whenever the player crosses chunk borders. The vanilla server also sends this
// whenever the vertical position of the player changes, even if he doesn't cross chunk borders.
type ClientboundUpdateViewPosition struct {
	// Chunk is the chunk coordinate of the player.
	Chunk voxel.V2
}

// ID returns the constant packet ID.
func (ClientboundUpdateViewPosition) ID() ID { return IDClientboundUpdateViewPosition }

// Name returns the constant packet name.
func (ClientboundUpdateViewPosition) Name() string { return "Update View Position" }

// EncodeInto writes this packet into the given writer.
func (c ClientboundUpdateViewPosition) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteVarInt("chunk x", c.Chunk.X)
	enc.WriteVarInt("chunk z", c.Chunk.Z)

	return
}
