package game

import (
	"fmt"
)

func (g *Game) initialize() error {
	if err := g.loadWorld(); err != nil {
		return fmt.Errorf("load world: %w", err)
	}

	g.chunkService = newChunkService(g.world, g.world.generator)

	close(g.ready)
	return nil
}
