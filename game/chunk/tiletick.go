package chunk

import (
	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/game/voxel"
)

type TileTick struct {
	ID                   id.ID
	TicksUntilProcessing int
	Priority             int
	Pos                  *voxel.V3
}
