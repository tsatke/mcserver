package entity

import (
	"encoding/binary"

	"github.com/google/uuid"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/id"
)

type nbtDecodeFunc func(nbt.Mapper) (Entity, error)

var (
	nbtDecoders = map[id.ID]nbtDecodeFunc{
		// mobs
		id.ParseID("minecraft:bat"):      decodeBat,
		id.ParseID("minecraft:creeper"):  decodeCreeper,
		id.ParseID("minecraft:cow"):      decodeCow,
		id.ParseID("minecraft:pig"):      decodePig,
		id.ParseID("minecraft:sheep"):    decodeSheep,
		id.ParseID("minecraft:skeleton"): decodeSkeleton,
		id.ParseID("minecraft:wolf"):     decodeWolf,
		// vehicles
		id.ParseID("minecraft:chest_minecart"): decodeChestMinecart,
		// item and xp orbs
		id.ParseID("minecraft:item"): decodeItem,
	}
)

func recoverAndSetErr(err *error) {
	if rec := recover(); rec != nil {
		if recErr, ok := rec.(error); ok {
			*err = recErr
		} else {
			panic(rec)
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func byteToBool(target *bool) func(nbt.Tag) error {
	return func(tag nbt.Tag) error {
		*target = tag.(*nbt.Byte).Value != 0
		return nil
	}
}

func intsToUUID(target *uuid.UUID) func(nbt.Tag) error {
	return func(tag nbt.Tag) error {
		buf := make([]byte, 16)
		for i, val := range tag.(*nbt.IntArray).Value {
			binary.BigEndian.PutUint32(buf[i*4:], uint32(val))
		}
		var err error
		*target, err = uuid.FromBytes(buf)
		return err
	}
}

func stringToID(target *id.ID) func(nbt.Tag) error {
	return func(tag nbt.Tag) error {
		*target = id.ParseID(tag.(*nbt.String).Value)
		return nil
	}
}
