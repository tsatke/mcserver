package game

import (
	"fmt"

	"github.com/tsatke/mcserver/game/chunk"
	"github.com/tsatke/mcserver/game/voxel"
	"github.com/tsatke/mcserver/game/worldgen"
)

type ChunkService interface {
	// Chunk returns and if necessary also loads the chunk on the
	// given voxel.
	Chunk(voxel.V2) (*chunk.Chunk, error)
	// IsChunkLoaded determines whether the chunk on the given
	// voxel is already loaded.
	IsChunkLoaded(voxel.V2) bool
	// Unload unloads the chunk on the given voxel. If there is
	// no such chunk loaded, this is a no-op.
	Unload(voxel.V2)
}

type regionSource interface {
	LoadRegion(coord voxel.V2) (*Region, error)
}

type chunkService struct {
	regionSource   regionSource
	chunkGenerator worldgen.Generator
	loadedChunks   map[voxel.V2]*chunk.Chunk
}

func newChunkService(source regionSource, generator worldgen.Generator) *chunkService {
	return &chunkService{
		regionSource:   source,
		chunkGenerator: generator,
		loadedChunks:   make(map[voxel.V2]*chunk.Chunk),
	}
}

func (s *chunkService) Chunk(coord voxel.V2) (*chunk.Chunk, error) {
	if loaded, ok := s.loadedChunks[coord]; ok {
		return loaded, nil
	}

	return s.loadChunk(coord)
}

func (s *chunkService) loadChunk(coord voxel.V2) (*chunk.Chunk, error) {
	regionCoord := voxel.V2{coord.X >> 5, coord.Z >> 5}
	region, err := s.regionSource.LoadRegion(regionCoord)
	if err != nil {
		return nil, err
	}
	loaded, err := region.loadChunk(coord)
	if err != nil {
		return nil, fmt.Errorf("load chunk: %w", err)
	}
	s.loadedChunks[coord] = loaded
	return loaded, nil
}

func (s *chunkService) IsChunkLoaded(coord voxel.V2) bool {
	_, ok := s.loadedChunks[coord]
	return ok
}

func (s *chunkService) Unload(coord voxel.V2) {
	delete(s.loadedChunks, coord)
}
