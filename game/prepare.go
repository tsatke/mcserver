package game

import (
	"errors"
	"time"

	"github.com/tsatke/mcserver/game/voxel"
)

func (g *Game) prepare() {
	if err := g.loadWorld(); err != nil {
		g.log.Error().
			Err(err).
			Msg("unable to load world")
		return
	}

	g.loadSpawnArea()

	close(g.ready)
}

func (g *Game) loadSpawnArea() {
	start := time.Now()
	spawn := g.world.Spawn
	spawnChunk := voxel.V2{
		X: spawn.X >> 4,
		Z: spawn.Z >> 4,
	}

	radius := 1
	relevantChunks := voxel.CircleAround(spawnChunk, float64(radius))
	for _, chunkPos := range relevantChunks {
		if _, err := g.loadChunkAtCoord(chunkPos); err != nil {
			if errors.Is(err, ErrChunkNotGenerated) {
				_, _ = g.generateAndLoadChunk(chunkPos)
			} else {
				g.log.Error().
					Err(err).
					Stringer("chunk", chunkPos).
					Msg("unable to load chunk")
			}
		}
	}

	g.log.Debug().
		Stringer("took", time.Since(start)).
		Int("chunks", len(relevantChunks)).
		Int("radius", radius).
		Msg("spawn area loaded")
}
