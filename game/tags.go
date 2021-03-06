package game

import (
	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/network/packet"
)

func (g *Game) sendTags(p *Player) {
	g.WritePacket(p, packet.ClientboundTags{
		BlockTags: []packet.Tag{
			{Name: id.ParseID("minecraft:enderman_holdable")},
			{Name: id.ParseID("minecraft:banners")},
			{Name: id.ParseID("minecraft:soul_fire_base_blocks")},
			{Name: id.ParseID("minecraft:campfires")},
			{Name: id.ParseID("minecraft:infiniburn_nether")},
			{Name: id.ParseID("minecraft:flower_pots")},
			{Name: id.ParseID("minecraft:infiniburn_overworld")},
			{Name: id.ParseID("minecraft:wooden_fences")},
			{Name: id.ParseID("minecraft:piglin_repellents")},
			{Name: id.ParseID("minecraft:wall_post_override")},
			{Name: id.ParseID("minecraft:wooden_slabs")},
			{Name: id.ParseID("minecraft:portals")},
			{Name: id.ParseID("minecraft:small_flowers")},
			{Name: id.ParseID("minecraft:bamboo_plantable_on")},
			{Name: id.ParseID("minecraft:wooden_trapdoors")},
			{Name: id.ParseID("minecraft:pressure_plates")},
			{Name: id.ParseID("minecraft:jungle_logs")},
			{Name: id.ParseID("minecraft:wooden_stairs")},
			{Name: id.ParseID("minecraft:spruce_logs")},
			{Name: id.ParseID("minecraft:signs")},
			{Name: id.ParseID("minecraft:carpets")},
			{Name: id.ParseID("minecraft:base_stone_overworld")},
			{Name: id.ParseID("minecraft:wool")},
			{Name: id.ParseID("minecraft:wooden_buttons")},
			{Name: id.ParseID("minecraft:stairs")},
			{Name: id.ParseID("minecraft:wither_summon_base_blocks")},
			{Name: id.ParseID("minecraft:logs")},
			{Name: id.ParseID("minecraft:stone_bricks")},
			{Name: id.ParseID("minecraft:hoglin_repellents")},
			{Name: id.ParseID("minecraft:fire")},
			{Name: id.ParseID("minecraft:beehives")},
			{Name: id.ParseID("minecraft:ice")},
			{Name: id.ParseID("minecraft:base_stone_nether")},
			{Name: id.ParseID("minecraft:dragon_immune")},
			{Name: id.ParseID("minecraft:crops")},
			{Name: id.ParseID("minecraft:wall_signs")},
			{Name: id.ParseID("minecraft:slabs")},
			{Name: id.ParseID("minecraft:valid_spawn")},
			{Name: id.ParseID("minecraft:mushroom_grow_block")},
			{Name: id.ParseID("minecraft:guarded_by_piglins")},
			{Name: id.ParseID("minecraft:wooden_doors")},
			{Name: id.ParseID("minecraft:warped_stems")},
			{Name: id.ParseID("minecraft:standing_signs")},
			{Name: id.ParseID("minecraft:infiniburn_end")},
			{Name: id.ParseID("minecraft:trapdoors")},
			{Name: id.ParseID("minecraft:crimson_stems")},
			{Name: id.ParseID("minecraft:buttons")},
			{Name: id.ParseID("minecraft:flowers")},
			{Name: id.ParseID("minecraft:corals")},
			{Name: id.ParseID("minecraft:prevent_mob_spawning_inside")},
			{Name: id.ParseID("minecraft:wart_blocks")},
			{Name: id.ParseID("minecraft:climbable")},
			{Name: id.ParseID("minecraft:planks")},
			{Name: id.ParseID("minecraft:soul_speed_blocks")},
			{Name: id.ParseID("minecraft:dark_oak_logs")},
			{Name: id.ParseID("minecraft:rails")},
			{Name: id.ParseID("minecraft:coral_plants")},
			{Name: id.ParseID("minecraft:non_flammable_wood")},
			{Name: id.ParseID("minecraft:leaves")},
			{Name: id.ParseID("minecraft:walls")},
			{Name: id.ParseID("minecraft:coral_blocks")},
			{Name: id.ParseID("minecraft:beacon_base_blocks")},
			{Name: id.ParseID("minecraft:strider_warm_blocks")},
			{Name: id.ParseID("minecraft:fence_gates")},
			{Name: id.ParseID("minecraft:bee_growables")},
			{Name: id.ParseID("minecraft:shulker_boxes")},
			{Name: id.ParseID("minecraft:wooden_pressure_plates")},
			{Name: id.ParseID("minecraft:wither_immune")},
			{Name: id.ParseID("minecraft:acacia_logs")},
			{Name: id.ParseID("minecraft:anvil")},
			{Name: id.ParseID("minecraft:birch_logs")},
			{Name: id.ParseID("minecraft:tall_flowers")},
			{Name: id.ParseID("minecraft:wall_corals")},
			{Name: id.ParseID("minecraft:underwater_bonemeals")},
			{Name: id.ParseID("minecraft:stone_pressure_plates")},
			{Name: id.ParseID("minecraft:impermeable")},
			{Name: id.ParseID("minecraft:sand")},
			{Name: id.ParseID("minecraft:nylium")},
			{Name: id.ParseID("minecraft:gold_ores")},
			{Name: id.ParseID("minecraft:logs_that_burn")},
			{Name: id.ParseID("minecraft:fences")},
			{Name: id.ParseID("minecraft:saplings")},
			{Name: id.ParseID("minecraft:beds")},
			{Name: id.ParseID("minecraft:oak_logs")},
			{Name: id.ParseID("minecraft:unstable_bottom_center")},
			{Name: id.ParseID("minecraft:doors")},
		}, // TODO: send BlockTags
		ItemTags: []packet.Tag{
			{Name: id.ParseID("minecraft:banners")},
			{Name: id.ParseID("minecraft:soul_fire_base_blocks")},
			{Name: id.ParseID("minecraft:stone_crafting_materials")},
			{Name: id.ParseID("minecraft:wooden_fences")},
			{Name: id.ParseID("minecraft:piglin_repellents")},
			{Name: id.ParseID("minecraft:beacon_payment_items")},
			{Name: id.ParseID("minecraft:wooden_slabs")},
			{Name: id.ParseID("minecraft:small_flowers")},
			{Name: id.ParseID("minecraft:wooden_trapdoors")},
			{Name: id.ParseID("minecraft:jungle_logs")},
			{Name: id.ParseID("minecraft:lectern_books")},
			{Name: id.ParseID("minecraft:wooden_stairs")},
			{Name: id.ParseID("minecraft:spruce_logs")},
			{Name: id.ParseID("minecraft:signs")},
			{Name: id.ParseID("minecraft:carpets")},
			{Name: id.ParseID("minecraft:wool")},
			{Name: id.ParseID("minecraft:wooden_buttons")},
			{Name: id.ParseID("minecraft:stairs")},
			{Name: id.ParseID("minecraft:fishes")},
			{Name: id.ParseID("minecraft:logs")},
			{Name: id.ParseID("minecraft:stone_bricks")},
			{Name: id.ParseID("minecraft:creeper_drop_music_discs")},
			{Name: id.ParseID("minecraft:arrows")},
			{Name: id.ParseID("minecraft:slabs")},
			{Name: id.ParseID("minecraft:wooden_doors")},
			{Name: id.ParseID("minecraft:warped_stems")},
			{Name: id.ParseID("minecraft:trapdoors")},
			{Name: id.ParseID("minecraft:crimson_stems")},
			{Name: id.ParseID("minecraft:buttons")},
			{Name: id.ParseID("minecraft:flowers")},
			{Name: id.ParseID("minecraft:stone_tool_materials")},
			{Name: id.ParseID("minecraft:planks")},
			{Name: id.ParseID("minecraft:boats")},
			{Name: id.ParseID("minecraft:dark_oak_logs")},
			{Name: id.ParseID("minecraft:rails")},
			{Name: id.ParseID("minecraft:non_flammable_wood")},
			{Name: id.ParseID("minecraft:leaves")},
			{Name: id.ParseID("minecraft:walls")},
			{Name: id.ParseID("minecraft:coals")},
			{Name: id.ParseID("minecraft:wooden_pressure_plates")},
			{Name: id.ParseID("minecraft:acacia_logs")},
			{Name: id.ParseID("minecraft:anvil")},
			{Name: id.ParseID("minecraft:piglin_loved")},
			{Name: id.ParseID("minecraft:music_discs")},
			{Name: id.ParseID("minecraft:birch_logs")},
			{Name: id.ParseID("minecraft:tall_flowers")},
			{Name: id.ParseID("minecraft:sand")},
			{Name: id.ParseID("minecraft:gold_ores")},
			{Name: id.ParseID("minecraft:logs_that_burn")},
			{Name: id.ParseID("minecraft:fences")},
			{Name: id.ParseID("minecraft:saplings")},
			{Name: id.ParseID("minecraft:beds")},
			{Name: id.ParseID("minecraft:oak_logs")},
			{Name: id.ParseID("minecraft:doors")},
		}, // TODO: send ItemTags
		FluidTags: []packet.Tag{
			{Name: id.ParseID("minecraft:lava")},
			{Name: id.ParseID("minecraft:water")},
		}, // TODO: send FluidTags
		EntityTags: []packet.Tag{
			{Name: id.ParseID("minecraft:skeletons")},
			{Name: id.ParseID("minecraft:raiders")},
			{Name: id.ParseID("minecraft:arrows")},
			{Name: id.ParseID("minecraft:beehive_inhabitors")},
			{Name: id.ParseID("minecraft:impact_projectiles")},
		}, // TODO: send EntityTags
	})
}
