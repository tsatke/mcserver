package packet

import (
	"io"
	"reflect"
	"strconv"

	"github.com/tsatke/mcserver/game/id"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundTags{}))
}

type Tag struct {
	Name    id.ID
	Entries []int
}

type ClientboundTags struct {
	BlockTags  []Tag
	ItemTags   []Tag
	FluidTags  []Tag
	EntityTags []Tag
}

func (ClientboundTags) ID() ID       { return IDClientboundTags }
func (ClientboundTags) Name() string { return "Tags" }

func (c ClientboundTags) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeVarInt("block tags length", len(c.BlockTags))
	for i, blockTag := range c.BlockTags {
		enc.writeID("blockTags["+strconv.Itoa(i)+"] type", blockTag.Name)
		enc.writeVarInt("blockTags["+strconv.Itoa(i)+"] count", len(blockTag.Entries))
		for j, entry := range blockTag.Entries {
			enc.writeVarInt("blockTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}
	enc.writeVarInt("item tags length", len(c.ItemTags))
	for i, itemTag := range c.ItemTags {
		enc.writeID("itemTags["+strconv.Itoa(i)+"] type", itemTag.Name)
		enc.writeVarInt("itemTags["+strconv.Itoa(i)+"] count", len(itemTag.Entries))
		for j, entry := range itemTag.Entries {
			enc.writeVarInt("itemTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}
	enc.writeVarInt("fluid tags length", len(c.FluidTags))
	for i, fluidTag := range c.FluidTags {
		enc.writeID("fluidTags["+strconv.Itoa(i)+"] type", fluidTag.Name)
		enc.writeVarInt("fluidTags["+strconv.Itoa(i)+"] count", len(fluidTag.Entries))
		for j, entry := range fluidTag.Entries {
			enc.writeVarInt("fluidTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}
	enc.writeVarInt("entity tags length", len(c.EntityTags))
	for i, entityTag := range c.EntityTags {
		enc.writeID("entityTags["+strconv.Itoa(i)+"] type", entityTag.Name)
		enc.writeVarInt("entityTags["+strconv.Itoa(i)+"] count", len(entityTag.Entries))
		for j, entry := range entityTag.Entries {
			enc.writeVarInt("entityTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}

	return
}
