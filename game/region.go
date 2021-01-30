package game

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"

	"github.com/tsatke/mcserver/game/chunk"
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

type Region struct {
	log        zerolog.Logger
	source     readWriteCloserAt
	locations  []RegionChunkLocation
	timestamps []RegionChunkTimestamp
}

type RegionChunkLocation struct {
	// Offset is the offset in 4KiB sectors from the start
	// of the file.
	Offset int
	// SectorCount indicates the size of the chunk in 4KiB sectors.
	SectorCount int
}

type RegionChunkTimestamp struct {
	// Timestamp indicates the last modification time of the respective chunk
	// in epoch seconds.
	Timestamp int
}

func loadRegion(log zerolog.Logger, rd readWriteCloserAt) (*Region, error) {
	// locations
	locations := make([]RegionChunkLocation, 0, 1024)
	offsetBuf := make([]byte, 4)
	sectorCountBuf := make([]byte, 1)
	for i := int64(0); i < 4096; i += 4 {
		if _, err := rd.ReadAt(offsetBuf[1:], i); err != nil {
			return nil, fmt.Errorf("read offset at %d: %w", i, err)
		}
		if _, err := rd.ReadAt(sectorCountBuf, i); err != nil {
			return nil, fmt.Errorf("read sector count at %d: %w", i, err)
		}
		locations = append(locations, RegionChunkLocation{
			Offset:      int(binary.BigEndian.Uint32(offsetBuf)),
			SectorCount: int(sectorCountBuf[0]),
		})
	}
	// timestamps
	timestamps := make([]RegionChunkTimestamp, 0, 1024)
	timestampBuf := make([]byte, 4)
	for i := int64(4096); i < 8192; i += 4 {
		if _, err := rd.ReadAt(timestampBuf, i); err != nil {
			return nil, fmt.Errorf("read at %d: %w", i, err)
		}
		timestamps = append(timestamps, RegionChunkTimestamp{
			Timestamp: int(binary.BigEndian.Uint32(timestampBuf)),
		})
	}

	return &Region{
		log:        log,
		source:     rd,
		locations:  locations,
		timestamps: timestamps,
	}, nil
}

func (r *Region) Close() error {
	return r.source.Close()
}

// loadChunk will return nil if the requested chunk is not yet generated.
func (r *Region) loadChunk(chunkCoord voxel.V2) (*chunk.Chunk, error) {
	chunkIndex := (chunkCoord.X & 31) + (chunkCoord.Z&31)*32
	chunkLocation := r.locations[chunkIndex]
	chunkTimestamp := r.timestamps[chunkIndex]
	c := &chunk.Chunk{
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
	c.CompressionType = chunk.CompressionType(chunkHeader[4])

	compressedReader := io.NewSectionReader(r.source, c.Offset+int64(len(chunkHeader)), int64(c.Length-1))

	if c.CompressionType == 0 {
		// compression type of 0 suggests that there is no chunk at this position, since the memory area
		// is still zeroed
		return nil, ErrChunkNotGenerated
	}

	decompressorReader, err := chunk.Decompressor[c.CompressionType](compressedReader)
	if err != nil {
		return nil, fmt.Errorf("create decompressor for %s: %w", c.CompressionType, err)
	}

	return chunk.LoadChunk(r.log, decompressorReader)
}
