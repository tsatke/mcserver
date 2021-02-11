package chunk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tsatke/mcserver/game/block"
	"github.com/tsatke/mcserver/game/voxel"
)

func TestChunkSuite(t *testing.T) {
	suite.Run(t, new(ChunkSuite))
}

type ChunkSuite struct {
	suite.Suite
}

func (suite *ChunkSuite) TestSetBlockAt() {
	ch := &Chunk{
		Coord:        voxel.V2{0, 0},
		LastModified: time.Now(),
		Data: &Data{
			Level: Level{},
		},
	}
	suite.Equal(airBlock, ch.BlockAt(voxel.V3{0, 0, 0}))
	ch.SetBlockAt(voxel.V3{0, 0, 0}, bedrockBlock)
	suite.Equal(bedrockBlock, ch.BlockAt(voxel.V3{0, 0, 0}))

	suite.Equal(airBlock, ch.BlockAt(voxel.V3{0, 0, 1}))
	ch.SetBlockAt(voxel.V3{0, 0, 1}, bedrockBlock)
	suite.Equal(bedrockBlock, ch.BlockAt(voxel.V3{0, 0, 1}))

	suite.Equal(airBlock, ch.BlockAt(voxel.V3{15, 15, 15}))
	ch.SetBlockAt(voxel.V3{15, 15, 15}, stoneBlock)
	for y := 0; y < 256; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				expectedBlock := airBlock
				if x == 0 && y == 0 && (z == 0 || z == 1) {
					expectedBlock = bedrockBlock
				} else if x == 15 && y == 15 && z == 15 {
					expectedBlock = stoneBlock
				}
				pos := voxel.V3{x, y, z}
				suite.Equal(expectedBlock, ch.BlockAt(pos), pos)
			}
		}
	}
	suite.Equal(stoneBlock, ch.BlockAt(voxel.V3{15, 15, 15}))

	// functionality is tested, now check internal values

	// check that all sections except sec[0] are still empty
	for i := 1; i < len(ch.Data.Level.Sections); i++ {
		suite.Zero(ch.Data.Level.Sections[i])
	}
	sec0 := ch.Data.Level.Sections[0]
	suite.EqualValues(0, sec0.Y)
	suite.ElementsMatch([]block.Block{
		airBlock,
		bedrockBlock,
		stoneBlock,
	}, sec0.Palette)
	suite.Equal(4096, len(sec0.paletteIndices))
	for i := 0; i < 4096; i++ {
		expected := airBlock
		switch i {
		case 0, 16:
			expected = bedrockBlock
		case 4095:
			expected = stoneBlock
		}
		suite.Equalf(expected, sec0.Palette[sec0.paletteIndices[i]], "at offset %d", i)
	}
}
