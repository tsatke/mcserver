package entity

import (
	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
)

type TileEntity interface {
	BlockID() id.ID
}

type BlockEntityData struct {
	ID         id.ID
	Pos        *voxel.V3
	KeepPacked byte
}

func (d *BlockEntityData) BlockID() id.ID {
	return d.ID
}
