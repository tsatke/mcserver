package packet

import (
	"io"
	"reflect"

	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/voxel"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundChunkData{}))
}

type ClientboundChunkData struct {
	ChunkPos voxel.V2
	// FullChunk controls whether the client should create a NEW chunk.
	// When this is false, this packet is just a large multi-block-update,
	// which changes all blocks in the given sections at once. Note that
	// biome data is not sent when this is false, so once a chunk is created,
	// biome data in the chunk can not be modified.
	// Sections that are not specified in the PrimaryBitMask are either empty
	// (full=true) or not changed (full=false).
	FullChunk      bool
	PrimaryBitMask uint16
	Heightmaps     nbt.Tag
	Biomes         []int
	Data           []byte
	BlockEntities  []nbt.Tag
}

func (ClientboundChunkData) ID() ID       { return IDClientboundChunkData }
func (ClientboundChunkData) Name() string { return "Chunk Data" }

func (c ClientboundChunkData) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeInt("chunk x", int32(c.ChunkPos.X))
	enc.writeInt("chunk z", int32(c.ChunkPos.Z))
	enc.writeBoolean("full chunk", c.FullChunk)
	enc.writeVarInt("primary bit mask", int(c.PrimaryBitMask))
	enc.writeNBT("heightmaps", c.Heightmaps)
	if c.FullChunk {
		enc.writeVarInt("biomes length", len(c.Biomes))
		for _, biome := range c.Biomes {
			enc.writeVarInt("biomes", biome)
		}
	}
	enc.writeVarInt("size", len(c.Data))
	enc.writeByteArray("data", c.Data)
	enc.writeVarInt("number of block entities", len(c.BlockEntities))
	for _, entity := range c.BlockEntities {
		enc.writeNBT("block entity", entity)
	}

	return
}
