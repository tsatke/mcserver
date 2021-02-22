package worldgen

import (
	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
	"github.com/tsatke/mcserver/game/world"
)

type Generator interface {
	ID() id.ID
	GenerateChunk(voxel.V2) world.Chunk
}
