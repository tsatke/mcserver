package packet

import (
	"fmt"
	"io"
	"math/bits"
	"reflect"

	"github.com/tsatke/mcserver/game/voxel"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundUpdateLight{}))
}

type ClientboundUpdateLight struct {
	ChunkPos            voxel.V2
	TrustEdges          bool
	SkyLightMask        int
	BlockLightMask      int
	EmptySkyLightMask   int
	EmptyBlockLightMask int
	SkyLightArrays      [][2048]byte
	BlockLightArrays    [][2048]byte
}

func (ClientboundUpdateLight) ID() ID       { return IDClientboundUpdateLight }
func (ClientboundUpdateLight) Name() string { return "Update Light" }

func (c ClientboundUpdateLight) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeInt("chunk x", int32(c.ChunkPos.X))
	enc.writeInt("chunk z", int32(c.ChunkPos.Z))
	enc.writeBoolean("trust edges", c.TrustEdges)
	enc.writeVarInt("sky light mask", c.SkyLightMask)
	enc.writeVarInt("block light mask", c.BlockLightMask)
	enc.writeVarInt("empty sky light mask", c.EmptySkyLightMask)
	enc.writeVarInt("empty block light mask", c.EmptyBlockLightMask)
	skyLightMaskHiCount := bits.OnesCount(uint(c.SkyLightMask))
	if len(c.SkyLightArrays) != skyLightMaskHiCount {
		return fmt.Errorf("skyLightMaskHiCount does not match the amount of arrays")
	}
	blockLightMaskHiCount := bits.OnesCount(uint(c.BlockLightMask))
	if len(c.BlockLightArrays) != blockLightMaskHiCount {
		return fmt.Errorf("blockLightMaskHiCount does not match the amount of arrays")
	}
	for _, array := range c.SkyLightArrays {
		enc.writeVarInt("length", len(array))
		enc.writeByteArray("sky light array", array[:])
	}
	for _, array := range c.BlockLightArrays {
		enc.writeVarInt("length", len(array))
		enc.writeByteArray("block light array", array[:])
	}

	return
}
