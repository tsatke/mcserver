package chunk

import (
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/entity"
)

type (
	Data struct {
		DataVersion int
		Level       Level
	}

	Level struct {
		// XPos is the x coordingate of this chunk relative to (0,0), NOT relative
		// to the region.
		XPos int
		// ZPos is the z coordingate of this chunk relative to (0,0), NOT relative
		// to the region.
		ZPos int
		// LastUpdate is the tick in which the chunk was last saved.
		LastUpdate int64
		// InhabitedTime is the cumulative number of ticks players have been in
		// this chunk. Note that this value increases faster when more players
		// are in the chunk. Used for regional difficulty: increases the chances
		// of mobs spawning with equipment, the chances of that equipment having
		// enchantments, the chances of spiders having potion effects, the chances
		// of mobs having the ability to pick up dropped items, and the chances of
		// zombies having the ability to spawn other zombies when attacked. Note
		// that at values 3600000 and above, regional difficulty is effectively
		// maxed for this chunk. At values 0 and below, the difficulty is capped to
		// a minimum (thus, if this is set to a negative number, it behaves
		// identically to being set to 0, apart from taking time to build back up
		// to the positives).
		InhabitedTime     int64
		Biomes            []int
		Heightmaps        Heightmaps
		CarvingMasks      CarvingMasks
		Sections          [16]Section
		Entities          []entity.Entity
		TileEntities      []entity.TileEntity
		TileTicks         []TileTick
		LiquidTicks       []TileTick
		Lights            [][]int16
		LiquidsToBeTicked [][]int16
		ToBeTicked        [][]int16
		PostProcessing    [][]int16
		Status            Status
		Structures        interface{}
	}

	Heightmaps struct {
		MotionBlocking         []int64
		MotionBlockingNoLeaves []int64
		OceanFloor             []int64
		OceanFloorWG           []int64
		WorldSurface           []int64
		WorldSurfaceWG         []int64
	}

	CarvingMasks struct {
		Air    []int8
		Liquid []int8
	}
)

func (h Heightmaps) ToNBT() nbt.Tag {
	return nbt.NewCompoundTag("", []nbt.Tag{
		nbt.NewLongArrayTag("MOTION_BLOCKING", h.MotionBlocking),
	})
}
