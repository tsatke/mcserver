package entity

import (
	"github.com/google/uuid"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/id"
)

type InventoryItem struct {
	ID    id.ID
	Count int8
	Slot  int8
	Tag   map[string]interface{}
}

type Item struct {
	Data
	Age         int16
	Health      int16
	PickupDelay int16
	Owner       uuid.UUID
	Thrower     uuid.UUID
	Item        InventoryItem
}

type ExperienceOrb struct {
	Data
	Age int16
	// Health is the health of this orb. This must be read as int16, but stored as byte.
	// Range is from 0-255.
	Health byte
}

func decodeItem(mapper nbt.Mapper) (Entity, error) {
	data, err := decodeData(mapper)
	if err != nil {
		return nil, err
	}

	item := Item{
		Data: data,
	}

	must(mapper.MapShort("Age", &item.Age))
	must(mapper.MapShort("Health", &item.Health))
	must(mapper.MapShort("PickupDelay", &item.PickupDelay))
	_ = mapper.MapCustom("Owner", intsToUUID(&item.Owner))
	_ = mapper.MapCustom("Thrower", intsToUUID(&item.Thrower))
	must(mapper.MapByte("Item.Count", &item.Item.Count))
	_ = mapper.MapCustom("Item.id", func(tag nbt.Tag) error {
		item.Item.ID = id.ParseID(tag.(*nbt.String).Value)
		return nil
	})

	return &item, nil
}

func decodeInventoryItem(mapper nbt.Mapper) (InventoryItem, error) {
	item := InventoryItem{}

	must(mapper.MapByte("Count", &item.Count))
	_ = mapper.MapByte("Slot", &item.Slot)
	must(mapper.MapCustom("id", stringToID(&item.ID)))
	// TODO: tag

	return item, nil
}
