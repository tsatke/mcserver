package world

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"time"

	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/block"
	"github.com/tsatke/mcserver/game/entity"
	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
)

//go:generate stringer -linecomment -type=CompressionType

type CompressionType byte

const (
	CompressionGZip         CompressionType = iota + 1 // GZip
	CompressionZlib                                    // Zlib
	CompressionUncompressed                            // uncompressed
)

var (
	Decompressor = [4]func(io.Reader) (io.Reader, error){
		func(_ io.Reader) (io.Reader, error) {
			return nil, fmt.Errorf("unknown compression 0x00")
		},
		func(in io.Reader) (io.Reader, error) {
			return gzip.NewReader(in)
		},
		func(in io.Reader) (io.Reader, error) {
			return zlib.NewReader(in)
		},
		func(in io.Reader) (io.Reader, error) {
			return in, nil
		},
	}
)

type (
	vanillaChunk struct {
		// Coord are the chunk coordinates (XZ).
		Coord voxel.V2
		// Offset is the offset in bytes from the beginning of the region file.
		// This value is already multiplied with 4096.
		Offset int64
		// SectorCount is the ceiled length of the chunk in the region file in
		// 4096 byte sectors. For the chunk length in bytes, see paddedLength.
		SectorCount int
		// PaddedLength is the length of the chunk in the region file in
		// bytes. This value is already multiplied with 4096.
		PaddedLength int
		// Length is the actual length of the chunk, i.e. the not-padded length.
		Length int
		// CompressionType is the compression type, in which the chunk data was
		// compressed in the region file.
		CompressionType CompressionType
		LastModified    time.Time

		DataVersion int
		// XPos is the x coordingate of this chunk relative to (0,0), NOT relative
		// to the region.
		XPos int
		// ZPos is the z coordingate of this chunk relative to (0,0), NOT relative
		// to the region.
		ZPos int
		// LastUpdate is the tick in which the chunk was last saved.
		LastUpdate int64
		// InhabitedTime is the cumulative number of ticks players have been in
		// this chunk. Note that this value increases faster when more players
		// are in the chunk. Used for regional difficulty: increases the chances
		// of mobs spawning with equipment, the chances of that equipment having
		// enchantments, the chances of spiders having potion effects, the chances
		// of mobs having the ability to pick up dropped items, and the chances of
		// zombies having the ability to spawn other zombies when attacked. Note
		// that at values 3600000 and above, regional difficulty is effectively
		// maxed for this chunk. At values 0 and below, the difficulty is capped to
		// a minimum (thus, if this is set to a negative number, it behaves
		// identically to being set to 0, apart from taking time to build back up
		// to the positives).
		InhabitedTime     int64
		Biomes            []int
		Heightmaps        Heightmaps
		CarvingMasks      CarvingMasks
		Sections          [16]vanillaSection
		Entities          []entity.Entity
		TileEntities      []entity.TileEntity
		TileTicks         []TileTick
		LiquidTicks       []TileTick
		Lights            [][]int16
		LiquidsToBeTicked [][]int16
		ToBeTicked        [][]int16
		PostProcessing    [][]int16
		Status            Status
		Structures        interface{}
	}

	Heightmaps struct {
		MotionBlocking         []int64
		MotionBlockingNoLeaves []int64
		OceanFloor             []int64
		OceanFloorWG           []int64
		WorldSurface           []int64
		WorldSurfaceWG         []int64
	}

	CarvingMasks struct {
		Air    []int8
		Liquid []int8
	}

	TileTick struct {
		ID                   id.ID
		TicksUntilProcessing int
		Priority             int
		Pos                  *voxel.V3
	}
)

func (h Heightmaps) ToNBT() nbt.Tag {
	return nbt.NewCompoundTag("", []nbt.Tag{
		nbt.NewLongArrayTag("MOTION_BLOCKING", h.MotionBlocking),
	})
}

func (c *vanillaChunk) Pos() voxel.V2 {
	return voxel.V2{c.XPos, c.ZPos}
}

func (c *vanillaChunk) BlockAt(v3 voxel.V3) block.Block {
	sectionRelativePos := voxel.V3{v3.X, v3.Y % 16, v3.Z}
	return c.Sections[v3.Y>>4].BlockAt(sectionRelativePos)
}

func (c *vanillaChunk) SetBlockAt(v3 voxel.V3, newBlock block.Block) {
	section := &c.Sections[v3.Y>>4]
	paletteIndex := section.paletteIndexOf(newBlock)
	if paletteIndex == -1 {
		if len(section.Palette) == 0 {
			section.Palette = []block.Block{airBlock} // air should always be first block in palette
		}
		// we need to add the newBlock to the palette
		section.Palette = append(section.Palette, newBlock)
		paletteIndex = len(section.Palette) - 1
	}

	indices := section.PaletteIndices()
	indexOffset := v3.Y*16*16 + v3.Z*16 + v3.X
	if len(indices) <= indexOffset {
		// we need to grow the indices
		newIndices := make([]uint64, indexOffset+1)
		// fill newIndices with reference to air
		airIndex := section.paletteIndexOf(airBlock)
		if airIndex == -1 {
			// we need to add air to the palette
			section.Palette = append(section.Palette, airBlock)
			airIndex = len(section.Palette) - 1
		}
		for i := range newIndices {
			newIndices[i] = uint64(airIndex)
		}
		// copy old indices to the newIndices, overwriting the air values
		copy(newIndices, indices)
		indices = newIndices
		// save back into section
		section.paletteIndices = indices
	}
	indices[indexOffset] = uint64(paletteIndex)
}
