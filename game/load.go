package game

import "github.com/tsatke/mcserver/game/voxel"

func (g *Game) loadChunksInSquare(center voxel.V2, radius int) {
	for x := center.X - radius; x <= center.X+radius; x++ {
		for z := center.X - radius; z <= center.X+radius; z++ {
			coord := voxel.V2{x, z}
			if _, err := g.world.Chunk(coord); err != nil {
				g.log.Error().
					Err(err).
					Stringer("chunk", coord).
					Msg("load chunk")
			}
		}
	}
}
