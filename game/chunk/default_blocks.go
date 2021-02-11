package chunk

import (
	"fmt"

	"github.com/tsatke/mcserver/game/block"
)

var (
	airBlock     = mustFromDesc(block.Air)
	bedrockBlock = mustFromDesc(block.Bedrock)
	stoneBlock   = mustFromDesc(block.Stone)
)

func mustFromDesc(desc block.BlockDescriptor) block.Block {
	b, err := block.CreateFromDescriptor(desc)
	if err != nil {
		panic(fmt.Errorf("unable to create block %s: %w", desc.ID, err))
	}
	return b
}
