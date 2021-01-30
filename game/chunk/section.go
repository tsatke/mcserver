package chunk

import (
	"math"
	"reflect"

	"github.com/tsatke/mcserver/game/block"
	"github.com/tsatke/mcserver/game/voxel"
)

type Section struct {
	paletteIndices []uint64

	Y           int8
	Palette     []block.Block
	BlockLight  []int8
	BlockStates []int64
	SkyLight    []int8
}

func (s *Section) BlockAt(position voxel.V3) block.Block {
	blockPos := (position.Y%16)*16*16 + position.Z*16 + position.X
	paletteIndicesCache := s.PaletteIndices()
	if blockPos >= len(paletteIndicesCache) {
		return block.Air
	}
	paletteIndex := paletteIndicesCache[blockPos]
	if int(paletteIndex) >= len(s.Palette) {
		return block.Air
	}
	return s.Palette[paletteIndex]
}

func (s *Section) PaletteIndices() []uint64 {
	if len(s.paletteIndices) == 0 && len(s.Palette) > 0 {
		segmentLength := int(math.Floor(math.Log2(float64(len(s.Palette))))) + 1
		s.paletteIndices = splitArrayIntoBitSegments(s.BlockStates, segmentLength)
	}
	return s.paletteIndices
}

func (s *Section) paletteIndexOf(block block.Block) int {
	for i, pb := range s.Palette {
		if reflect.DeepEqual(pb, block) {
			return i
		}
	}
	return -1
}

func splitArrayIntoBitSegments(array []int64, segmentLength int) []uint64 {
	var result []uint64

	mask := createHiMask(segmentLength)
	for _, elem := range array {
		for i := 0; i < 64/segmentLength; i++ {
			result = append(result, uint64(elem)&mask)
			elem = elem >> segmentLength
		}
	}

	return result
}

func createHiMask(size int) uint64 {
	return 1<<size - 1
}