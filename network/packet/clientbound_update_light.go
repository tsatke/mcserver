package packet

import (
	"fmt"
	"io"
	"math/bits"
	"reflect"

	"github.com/tsatke/mcserver/game/voxel"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundUpdateLight{}))
}

// ClientboundUpdateLight is sent to update light levels in a chunk
type ClientboundUpdateLight struct {
	// ChunkPos is the position of the referred chunk.
	ChunkPos voxel.V2
	// TrustEdges determines if edges should be trusted for light updates.
	// I have no idea what that means.
	TrustEdges bool
	// SkyLightMask is a bitmask containing 18 bits, with the lowest bit
	// corresponding to chunk section -1 (in the void, y=-16 to y=-1) and
	// the highest bit for chunk section 16 (above the world, y=256 to y=271).
	SkyLightMask int
	// BlockLightMask is a bitmask containing 18 bits, with the lowest bit
	// corresponding to chunk section -1 (in the void, y=-16 to y=-1) and
	// the highest bit for chunk section 16 (above the world, y=256 to y=271).
	BlockLightMask int
	// EmptySkyLightMask is a bit mask containing 18 bits, which indicates
	// sections that have 0 for all their sky light values. If a section
	// is set in both this mask and the main sky light mask, it is ignored
	// for this mask and it is included in the sky light arrays (the
	// vanilla server does not create such masks). If it is only set in
	// this mask, it is not included in the sky light arrays.
	EmptySkyLightMask int
	// EmptyBlockLightMask is a bit mask containing 18 bits which indicates
	// sections that have 0 for all their block light values. If a section
	// is set in both this mask and the main block light mask, it is ignored
	// for this mask and it is included in the block light arrays (the
	// vanilla server does not create such masks). If it is only set in
	// this mask, it is not included in the block light arrays.
	EmptyBlockLightMask int
	// SkyLightArrays contains 1 array for each bit set to true in the
	// sky light mask, starting with the lowest value. Half a
	// byte per light value.
	SkyLightArrays [][2048]byte
	// BlockLightArrays contains 1 array for each bit set to true in the
	// sky light mask, starting with the lowest value. Half a
	// byte per light value.
	BlockLightArrays [][2048]byte
}

// ID returns the constant packet ID.
func (ClientboundUpdateLight) ID() ID { return IDClientboundUpdateLight }

// Name returns the constant packet name.
func (ClientboundUpdateLight) Name() string { return "Update Light" }

// EncodeInto writes this packet into the given writer.
func (c ClientboundUpdateLight) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteVarInt("chunk x", c.ChunkPos.X)
	enc.WriteVarInt("chunk z", c.ChunkPos.Z)
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
