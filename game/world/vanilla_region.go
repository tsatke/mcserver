package world

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

type readWriterAt interface {
	io.ReaderAt
	io.WriterAt
}

type readWriteCloserAt interface {
	readWriterAt
	io.Closer
}

type vanillaRegion struct {
	log        zerolog.Logger
	source     readWriteCloserAt
	locations  []vanillaRegionChunkLocation
	timestamps []vanillaRegionChunkTimestamp
}

type vanillaRegionChunkLocation struct {
	// Offset is the offset in 4KiB sectors from the start
	// of the file.
	Offset int
	// SectorCount indicates the size of the chunk in 4KiB sectors.
	SectorCount int
}

type vanillaRegionChunkTimestamp struct {
	// Timestamp indicates the last modification time of the respective chunk
	// in epoch seconds.
	Timestamp int
}

func loadVanillaRegion(rd readWriteCloserAt) (*vanillaRegion, error) {
	// locations
	locations := make([]vanillaRegionChunkLocation, 0, 1024)
	offsetBuf := make([]byte, 4)
	sectorCountBuf := make([]byte, 1)
	for i := int64(0); i < 4096; i += 4 {
		if _, err := rd.ReadAt(offsetBuf[1:], i); err != nil {
			return nil, fmt.Errorf("read offset at %d: %w", i, err)
		}
		if _, err := rd.ReadAt(sectorCountBuf, i); err != nil {
			return nil, fmt.Errorf("read sector count at %d: %w", i, err)
		}
		locations = append(locations, vanillaRegionChunkLocation{
			Offset:      int(binary.BigEndian.Uint32(offsetBuf)),
			SectorCount: int(sectorCountBuf[0]),
		})
	}
	// timestamps
	timestamps := make([]vanillaRegionChunkTimestamp, 0, 1024)
	timestampBuf := make([]byte, 4)
	for i := int64(4096); i < 8192; i += 4 {
		if _, err := rd.ReadAt(timestampBuf, i); err != nil {
			return nil, fmt.Errorf("read at %d: %w", i, err)
		}
		timestamps = append(timestamps, vanillaRegionChunkTimestamp{
			Timestamp: int(binary.BigEndian.Uint32(timestampBuf)),
		})
	}

	return &vanillaRegion{
		source:     rd,
		locations:  locations,
		timestamps: timestamps,
	}, nil
}

func (r *vanillaRegion) Close() error {
	return r.source.Close()
}

func (r *vanillaRegion) loadChunk(chunkCoord voxel.V2) (*vanillaChunk, error) {
	chunkIndex := (chunkCoord.X & 31) + (chunkCoord.Z&31)*32
	chunkLocation := r.locations[chunkIndex]
	chunkTimestamp := r.timestamps[chunkIndex]
	c := &vanillaChunk{
		Coord:        chunkCoord,
		Offset:       int64(chunkLocation.Offset) * 4096,
		SectorCount:  chunkLocation.SectorCount,
		PaddedLength: chunkLocation.SectorCount * 4096,
		LastModified: time.Unix(int64(chunkTimestamp.Timestamp), 0),
	}

	chunkHeader := make([]byte, 5) // length: 0-3, compressionType: 4
	if _, err := r.source.ReadAt(chunkHeader, c.Offset); err != nil {
		return nil, fmt.Errorf("read chunk header: %w", err)
	}
	c.Length = int(binary.BigEndian.Uint32(chunkHeader[:4]))
	c.CompressionType = CompressionType(chunkHeader[4])

	compressedReader := io.NewSectionReader(r.source, c.Offset+int64(len(chunkHeader)), int64(c.Length-1))

	if c.CompressionType == 0 {
		// compression type of 0 suggests that there is no chunk at this position, since the memory area
		// is still zeroed
		return nil, ErrChunkNotGenerated
	}

	decompressorReader, err := Decompressor[c.CompressionType](compressedReader)
	if err != nil {
		return nil, fmt.Errorf("create decompressor for %s: %w", c.CompressionType, err)
	}

	tag, err := nbt.NewDecoder(decompressorReader, binary.BigEndian).ReadTag()
	if err != nil {
		return nil, fmt.Errorf("decode nbt: %w", err)
	}

	ch := &vanillaChunk{}
	if err := r.decodeChunkInto(ch, tag); err != nil {
		return nil, fmt.Errorf("decode chunk: %w", err)
	}

	return ch, nil
}

func (r *vanillaRegion) decodeChunkInto(ch *vanillaChunk, tag nbt.Tag) (err error) {
	mapper := nbt.NewSimpleMapper(tag)

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

	must(mapper.MapInt("DataVersion", &ch.DataVersion))
	must(mapper.MapInt("Level.xPos", &ch.XPos))
	must(mapper.MapInt("Level.zPos", &ch.ZPos))
	must(mapper.MapLong("Level.LastUpdate", &ch.LastUpdate))
	must(mapper.MapLong("Level.LastUpdate", &ch.LastUpdate))
	_ = mapper.MapIntArray("Level.Biomes", &ch.Biomes)
	must(mapper.MapLongArray("Level.Heightmaps.MOTION_BLOCKING", &ch.Heightmaps.MotionBlocking))
	must(mapper.MapLongArray("Level.Heightmaps.MOTION_BLOCKING_NO_LEAVES", &ch.Heightmaps.MotionBlockingNoLeaves))
	must(mapper.MapLongArray("Level.Heightmaps.OCEAN_FLOOR", &ch.Heightmaps.OceanFloor))
	// it seems like OCEAN_FLOOR_WG doesn't exist
	must(mapper.MapLongArray("Level.Heightmaps.WORLD_SURFACE", &ch.Heightmaps.WorldSurface))
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
			ch.Sections[secY].Y = secY
			must(mapper.MapList("Palette", func(size int) {
				ch.Sections[secY].Palette = make([]block.Block, size)
			}, func(paletteIndex int, mapper nbt.Mapper) error {
				var idString string
				must(mapper.MapString("Name", &idString))
				id := id.ParseID(idString)

				var properties []block.Property
				propertiesTag, _ := mapper.Query("Properties") // ignore the error, since properties seem to be optional
				if propertiesTag != nil {
					for name, value := range propertiesTag.(*nbt.Compound).Value {
						propertyDesc, ok := block.DescriptorForPropertyName(name)
						// if !ok {
						// 	return fmt.Errorf("no property descriptor for property %q", name)
						// }
						_, _ = propertyDesc, ok
						_ = value // TODO: decode this properly
					}
				}

				block, err := block.Create(id, properties...)
				if err != nil {
					return fmt.Errorf("create block: %w", err)
				}
				ch.Sections[secY].Palette[paletteIndex] = block
				return nil
			}))
			_ = mapper.MapByteArray("BlockLight", &ch.Sections[secY].BlockLight)
			must(mapper.MapLongArray("BlockStates", &ch.Sections[secY].BlockStates))
			_ = mapper.MapByteArray("SkyLight", &ch.Sections[secY].SkyLight)
			return nil
		},
	))
	must(mapper.MapList("Level.Entities", func(size int) {
		ch.Entities = make([]entity.Entity, size)
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

		ch.Entities[entityIndex] = e
		return nil
	}))
	return nil
}
