// Code generated by "blockgen -in=minecraft_blocks.yaml -out=minecraft_blocks.go -pkg=block"; DO NOT EDIT.

package block

import "github.com/tsatke/mcserver/game/id"

var (
	Air = BlockDescriptor{
		ID: id.ID{"minecraft", "air"},
	}
	Andesite = BlockDescriptor{
		ID: id.ID{"minecraft", "andesite"},
	}
	Bedrock = BlockDescriptor{
		ID: id.ID{"minecraft", "bedrock"},
	}
	CaveAir = BlockDescriptor{
		ID: id.ID{"minecraft", "cave_air"},
	}
	Clay = BlockDescriptor{
		ID: id.ID{"minecraft", "clay"},
	}
	CoalOre = BlockDescriptor{
		ID: id.ID{"minecraft", "coal_ore"},
	}
	Cobblestone = BlockDescriptor{
		ID: id.ID{"minecraft", "cobblestone"},
	}
	Cobweb = BlockDescriptor{
		ID: id.ID{"minecraft", "cobweb"},
	}
	Diorite = BlockDescriptor{
		ID: id.ID{"minecraft", "diorite"},
	}
	Dirt = BlockDescriptor{
		ID: id.ID{"minecraft", "dirt"},
	}
	GoldOre = BlockDescriptor{
		ID: id.ID{"minecraft", "gold_ore"},
	}
	Granite = BlockDescriptor{
		ID: id.ID{"minecraft", "granite"},
	}
	Grass = BlockDescriptor{
		ID: id.ID{"minecraft", "grass"},
	}
	GrassBlock = BlockDescriptor{
		ID: id.ID{"minecraft", "grass_block"},
	}
	Gravel = BlockDescriptor{
		ID: id.ID{"minecraft", "gravel"},
	}
	IronOre = BlockDescriptor{
		ID: id.ID{"minecraft", "iron_ore"},
	}
	LapisOre = BlockDescriptor{
		ID: id.ID{"minecraft", "lapis_ore"},
	}
	LilyPad = BlockDescriptor{
		ID: id.ID{"minecraft", "lily_pad"},
	}
	OakFence = BlockDescriptor{
		ID: id.ID{"minecraft", "oak_fence"},
	}
	OakLeaves = BlockDescriptor{
		ID: id.ID{"minecraft", "oak_leaves"},
	}
	OakLog = BlockDescriptor{
		ID: id.ID{"minecraft", "oak_log"},
	}
	OakPlanks = BlockDescriptor{
		ID: id.ID{"minecraft", "oak_planks"},
	}
	Rail = BlockDescriptor{
		ID: id.ID{"minecraft", "rail"},
	}
	RedstoneOre = BlockDescriptor{
		ID: id.ID{"minecraft", "redstone_ore"},
	}
	Seagrass = BlockDescriptor{
		ID: id.ID{"minecraft", "seagrass"},
	}
	Stone = BlockDescriptor{
		ID: id.ID{"minecraft", "stone"},
	}
	TallSeagrass = BlockDescriptor{
		ID: id.ID{"minecraft", "tall_seagrass"},
	}
	Vine = BlockDescriptor{
		ID: id.ID{"minecraft", "vine"},
	}
	VoidAir = BlockDescriptor{
		ID: id.ID{"minecraft", "void_air"},
	}
	Water = BlockDescriptor{
		ID: id.ID{"minecraft", "water"},
	}
)

func init() {
	Must(RegisterBlock(Air))
	Must(RegisterBlock(Andesite))
	Must(RegisterBlock(Bedrock))
	Must(RegisterBlock(CaveAir))
	Must(RegisterBlock(Clay))
	Must(RegisterBlock(CoalOre))
	Must(RegisterBlock(Cobblestone))
	Must(RegisterBlock(Cobweb))
	Must(RegisterBlock(Diorite))
	Must(RegisterBlock(Dirt))
	Must(RegisterBlock(GoldOre))
	Must(RegisterBlock(Granite))
	Must(RegisterBlock(Grass))
	Must(RegisterBlock(GrassBlock))
	Must(RegisterBlock(Gravel))
	Must(RegisterBlock(IronOre))
	Must(RegisterBlock(LapisOre))
	Must(RegisterBlock(LilyPad))
	Must(RegisterBlock(OakFence))
	Must(RegisterBlock(OakLeaves))
	Must(RegisterBlock(OakLog))
	Must(RegisterBlock(OakPlanks))
	Must(RegisterBlock(Rail))
	Must(RegisterBlock(RedstoneOre))
	Must(RegisterBlock(Seagrass))
	Must(RegisterBlock(Stone))
	Must(RegisterBlock(TallSeagrass))
	Must(RegisterBlock(Vine))
	Must(RegisterBlock(VoidAir))
	Must(RegisterBlock(Water))
}