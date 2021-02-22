package world

import (
	"github.com/tsatke/mcserver/game/voxel"
)

type World interface {
	// Chunk returns and if necessary also loads the chunk on the
	// given voxel.
	Chunk(voxel.V2) (Chunk, error)
	// IsChunkLoaded determines whether the chunk on the given
	// voxel is already loaded.
	IsChunkLoaded(voxel.V2) bool
	// Unload unloads the chunk on the given voxel. If there is
	// no such chunk loaded, this is a no-op.
	Unload(voxel.V2)

	Seed() int64
}
