package packet

import (
	"io"
	"reflect"
	"strconv"

	"github.com/tsatke/mcserver/game/id"
)

func init() {
	RegisterPacket(StatePlay, reflect.TypeOf(ClientboundTags{}))
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

	enc := Encoder{w}

	enc.WriteVarInt("block tags length", len(c.BlockTags))
	for i, blockTag := range c.BlockTags {
		enc.WriteID("blockTags["+strconv.Itoa(i)+"] type", blockTag.Name)
		enc.WriteVarInt("blockTags["+strconv.Itoa(i)+"] count", len(blockTag.Entries))
		for j, entry := range blockTag.Entries {
			enc.WriteVarInt("blockTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}
	enc.WriteVarInt("item tags length", len(c.ItemTags))
	for i, itemTag := range c.ItemTags {
		enc.WriteID("itemTags["+strconv.Itoa(i)+"] type", itemTag.Name)
		enc.WriteVarInt("itemTags["+strconv.Itoa(i)+"] count", len(itemTag.Entries))
		for j, entry := range itemTag.Entries {
			enc.WriteVarInt("itemTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}
	enc.WriteVarInt("fluid tags length", len(c.FluidTags))
	for i, fluidTag := range c.FluidTags {
		enc.WriteID("fluidTags["+strconv.Itoa(i)+"] type", fluidTag.Name)
		enc.WriteVarInt("fluidTags["+strconv.Itoa(i)+"] count", len(fluidTag.Entries))
		for j, entry := range fluidTag.Entries {
			enc.WriteVarInt("fluidTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}
	enc.WriteVarInt("entity tags length", len(c.EntityTags))
	for i, entityTag := range c.EntityTags {
		enc.WriteID("entityTags["+strconv.Itoa(i)+"] type", entityTag.Name)
		enc.WriteVarInt("entityTags["+strconv.Itoa(i)+"] count", len(entityTag.Entries))
		for j, entry := range entityTag.Entries {
			enc.WriteVarInt("entityTags["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] count", entry)
		}
	}

	return
}
