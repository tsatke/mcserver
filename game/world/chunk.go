package world

import (
	"github.com/tsatke/mcserver/game/block"
	"github.com/tsatke/mcserver/game/voxel"
)

type Chunk interface {
	Pos() voxel.V2

	BlockAt(voxel.V3) block.Block
	SetBlockAt(voxel.V3, block.Block)
}
