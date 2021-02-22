package world

import "github.com/tsatke/mcserver/game/block"

// Section is a 16x16x16 cube, which is part of a chunk.
// A chunk consists of multiple sections, where the actual
// block data is stored.
type Section interface {
	Palette() []block.Block
	Blocks() [16 * 16 * 16]int16
}

type Sectioner interface {
	Chunk
	Sections() []Section
}

// ChunkSection returns the section with the given section index from the given chunk.
// The section index is zero based, where 0 stands for the section from y=0 to y=15.
// A section index of -1 is not supported.
// If the given chunk implements the Sectioner interface, it will be used to extract
// the respective section.
// If not, the computation will be most expensive and will not be cached.
func ChunkSection(ch Chunk, sec int) Section {
	if sec < 0 {
		panic("section index must be >= 0")
	}

	if secs, ok := ch.(Sectioner); ok {
		return secs.Sections()[sec]
	}

	panic("implement me")
}
