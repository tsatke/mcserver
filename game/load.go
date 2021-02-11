package game

import "github.com/tsatke/mcserver/game/voxel"

func (g *Game) loadChunksInSquare(center voxel.V2, radius int) {
	for x := center.X - radius; x <= center.X+radius; x++ {
		for z := center.X - radius; z <= center.X+radius; z++ {
			coord := voxel.V2{x, z}
			if _, err := g.chunkService.Chunk(coord); err != nil {
				g.log.Error().
					Err(err).
					Stringer("chunk", coord).
					Msg("load chunk")
			}
		}
	}
}

// getNotifiedAfter will execute the given function in a separate goroutine
// and close the returned channel after the given function has returned.
// The following code snippets are probably functionally equivalent.
//
//	fn := ...
//	fn()
//
//	fn := ...
//	<-getNotifiedAfter(fn)
func getNotifiedAfter(fn func()) <-chan struct{} {
	ch := make(chan struct{})
	defer close(ch)
	go fn()
	return ch
}
