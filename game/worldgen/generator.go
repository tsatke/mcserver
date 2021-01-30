package worldgen

import (
	"github.com/tsatke/mcserver/game/chunk"
	"github.com/tsatke/mcserver/game/voxel"
)

type Generator interface {
	GenerateChunk(voxel.V2) *chunk.Chunk
}
