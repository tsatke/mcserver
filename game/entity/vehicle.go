package entity

import (
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/block"
)

type (
	Minecart struct {
		CustomDisplayTitle bool
		DisplayState       block.Block
		DisplayOffset      int
	}

	ChestMinecart struct {
		Data
		Minecart
		Items         []InventoryItem
		LootTable     interface{}
		LootTableSeed int64
	}
)

func decodeMinecart(mapper nbt.Mapper) (Minecart, error) {
	mc := Minecart{}
	_ = mapper.MapCustom("CustomDisplayTitle", byteToBool(&mc.CustomDisplayTitle))
	// TODO: DisplayState
	_ = mapper.MapInt("DisplayOffset", &mc.DisplayOffset)
	return mc, nil
}

func decodeChestMinecart(mapper nbt.Mapper) (Entity, error) {
	base, err := decodeMinecart(mapper)
	if err != nil {
		return nil, err
	}

	mc := ChestMinecart{
		Minecart: base,
	}

	_ = mapper.MapList(
		"Items",
		func(size int) {
			mc.Items = make([]InventoryItem, size)
		},
		func(i int, mapper nbt.Mapper) error {

			return nil
		},
	)

	return &mc, nil
}
