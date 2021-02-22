package world

import (
	"github.com/tsatke/mcserver/game/block"
	"github.com/tsatke/mcserver/game/voxel"
)

var (
	airBlock     block.Block
	bedrockBlock block.Block
	stoneBlock   block.Block
)

func init() {
	must := func(b block.Block, err error) block.Block {
		if err != nil {
			panic(err)
		}
		return b
	}
	airBlock = must(block.CreateFromDescriptor(block.Air))
	bedrockBlock = must(block.CreateFromDescriptor(block.Bedrock))
	stoneBlock = must(block.CreateFromDescriptor(block.Stone))
}

var _ Section = (*emptySection)(nil)

type emptySection struct{}

func (e emptySection) Palette() []block.Block { return []block.Block{airBlock} }

func (e emptySection) Blocks() (res [4096]int16) { return }

var _ Chunk = (*emptyChunk)(nil)

type emptyChunk struct {
	pos voxel.V2
}

func (e emptyChunk) Pos() voxel.V2 { return e.pos }

func (e emptyChunk) BlockAt(v3 voxel.V3) block.Block {
	return airBlock
}

func (e emptyChunk) SetBlockAt(v3 voxel.V3, b block.Block) {
	panic("can't set block in immutable chunk")
}
