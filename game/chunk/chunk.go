package chunk

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/block"
	"github.com/tsatke/mcserver/game/entity"
	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
)

type Chunk struct {
	log zerolog.Logger

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
	Data            *Data
}

func LoadChunk(log zerolog.Logger, source io.Reader) (*Chunk, error) {
	chunk := &Chunk{
		log: log,
	}

	tag, err := nbt.NewDecoder(source, binary.BigEndian).ReadTag()
	if err != nil {
		panic(err)
	}

	if err := chunk.loadData(tag); err != nil {
		return nil, fmt.Errorf("load data: %w", err)
	}

	return chunk, nil
}

func (ch *Chunk) loadData(tag nbt.Tag) (err error) {
	mapper := nbt.NewSimpleMapper(tag)
	data := &Data{}

	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(error); ok {
				err = e
			} else {
				panic(rec)
			}
		}
	}()

	// for pretty API
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	must(mapper.MapInt("DataVersion", &data.DataVersion))
	must(mapper.MapInt("Level.xPos", &data.Level.XPos))
	must(mapper.MapInt("Level.zPos", &data.Level.ZPos))
	must(mapper.MapLong("Level.LastUpdate", &data.Level.LastUpdate))
	must(mapper.MapLong("Level.LastUpdate", &data.Level.LastUpdate))
	_ = mapper.MapIntArray("Level.Biomes", &data.Level.Biomes)
	must(mapper.MapLongArray("Level.Heightmaps.MOTION_BLOCKING", &data.Level.Heightmaps.MotionBlocking))
	must(mapper.MapLongArray("Level.Heightmaps.MOTION_BLOCKING_NO_LEAVES", &data.Level.Heightmaps.MotionBlockingNoLeaves))
	must(mapper.MapLongArray("Level.Heightmaps.OCEAN_FLOOR", &data.Level.Heightmaps.OceanFloor))
	// it seems like OCEAN_FLOOR_WG doesn't exist
	must(mapper.MapLongArray("Level.Heightmaps.WORLD_SURFACE", &data.Level.Heightmaps.WorldSurface))
	// it seems like WORLD_SURFACE_WG doesn't exist
	// it seems like CarvingMasks don't exist
	must(mapper.MapList(
		"Level.Sections",
		func(size int) {},
		func(_ int, mapper nbt.Mapper) error {
			var secY int8
			must(mapper.MapByte("Y", &secY))
			if secY == -1 {
				// apparently this value exists, but it isn't associated with any
				// other value, so we skip this section
				return nil
			}
			data.Level.Sections[secY].Y = secY
			must(mapper.MapList("Palette", func(size int) {
				data.Level.Sections[secY].Palette = make([]block.Block, size)
			}, func(paletteIndex int, mapper nbt.Mapper) error {
				target := &data.Level.Sections[secY].Palette[paletteIndex]
				var idString string
				must(mapper.MapString("Name", &idString))
				id := id.ParseID(idString)
				target.Name = id
				propertiesTag, err := mapper.Query("Properties")
				if err != nil {
					// properties seems to be optional, so skip it if there's an error
					return nil
				}
				properties := propertiesTag.(*nbt.Compound).Value
				target.Properties = make(map[string]interface{})
				for name, value := range properties {
					target.Properties[name] = value
				}
				return nil
			}))
			_ = mapper.MapByteArray("BlockLight", &data.Level.Sections[secY].BlockLight)
			must(mapper.MapLongArray("BlockStates", &data.Level.Sections[secY].BlockStates))
			_ = mapper.MapByteArray("SkyLight", &data.Level.Sections[secY].SkyLight)
			return nil
		},
	))
	must(mapper.MapList("Level.Entities", func(size int) {
		data.Level.Entities = make([]entity.Entity, size)
	}, func(entityIndex int, mapper nbt.Mapper) error {
		var entityID id.ID
		must(mapper.MapCustom("id", func(tag nbt.Tag) error {
			entityID = id.ParseID(tag.(*nbt.String).Value)
			return nil
		}))

		entityTag, queryErr := mapper.Query("")
		if queryErr != nil {
			return queryErr
		}

		e, decodeErr := entity.FromNBT(entityID, entityTag)
		if decodeErr != nil {
			return decodeErr
		}

		data.Level.Entities[entityIndex] = e
		return nil
	}))
	ch.Data = data
	return nil
}

// BlockAt returns the block at the given position. The position must be relative to
// this chunk. This implies, that X and Z can not exceed 15.
func (ch *Chunk) BlockAt(pos voxel.V3) block.Block {
	sectionRelativePos := voxel.V3{pos.X, pos.Y % 16, pos.Z}
	return ch.Data.Level.Sections[pos.Y>>4].BlockAt(sectionRelativePos)
}

func (ch *Chunk) SetBlockAt(pos voxel.V3, newBlock block.Block) {
	section := &ch.Data.Level.Sections[pos.Y>>4]
	paletteIndex := section.paletteIndexOf(newBlock)
	if paletteIndex == -1 {
		// we need to add the newBlock to the palette
		section.Palette = append(section.Palette, newBlock)
		paletteIndex = len(section.Palette) - 1
	}

	indices := section.PaletteIndices()
	indexOffset := pos.Y*16*16 + pos.Z*16 + pos.X
	if len(indices) <= indexOffset {
		// we need to grow the indices
		newIndices := make([]uint64, indexOffset+1)
		// fill newIndices with reference to air
		airIndex := section.paletteIndexOf(block.Air)
		if airIndex == -1 {
			// we need to add air to the palette
			section.Palette = append(section.Palette, block.Air)
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
