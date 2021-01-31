package packet

import (
	"fmt"
	"io"
	"math/bits"
	"reflect"

	"github.com/tsatke/mcserver/game/voxel"
)

func init() {
	RegisterPacket(StatePlay, reflect.TypeOf(ClientboundUpdateLight{}))
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

	enc := Encoder{w}

	enc.WriteInt("chunk x", int32(c.ChunkPos.X))
	enc.WriteInt("chunk z", int32(c.ChunkPos.Z))
	enc.WriteBoolean("trust edges", c.TrustEdges)
	enc.WriteVarInt("sky light mask", c.SkyLightMask)
	enc.WriteVarInt("block light mask", c.BlockLightMask)
	enc.WriteVarInt("empty sky light mask", c.EmptySkyLightMask)
	enc.WriteVarInt("empty block light mask", c.EmptyBlockLightMask)
	skyLightMaskHiCount := bits.OnesCount(uint(c.SkyLightMask))
	if len(c.SkyLightArrays) != skyLightMaskHiCount {
		return fmt.Errorf("skyLightMaskHiCount does not match the amount of arrays")
	}
	blockLightMaskHiCount := bits.OnesCount(uint(c.BlockLightMask))
	if len(c.BlockLightArrays) != blockLightMaskHiCount {
		return fmt.Errorf("blockLightMaskHiCount does not match the amount of arrays")
	}
	for i := range c.SkyLightArrays {
		enc.WriteVarInt("length", len(c.SkyLightArrays[i]))
		enc.WriteByteArray("sky light array", c.SkyLightArrays[i][:])
	}
	for i := range c.BlockLightArrays {
		enc.WriteVarInt("length", len(c.BlockLightArrays[i]))
		enc.WriteByteArray("block light array", c.BlockLightArrays[i][:])
	}

	return
}
