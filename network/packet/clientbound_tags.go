package packet

import (
	"io"
	"reflect"
	"strconv"

	"github.com/tsatke/mcserver/game/id"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundTags{}))
}

// Tag is a tag, which has an ID and a list of numeric
// IDs that are considered to have the property defined
// by the tag ID. An example is `minecraft:enderman_holdable`,
// which groups all blocks that can be held by an enderman.
type Tag struct {
	// Name is the ID of this tag.
	Name id.ID
	// Entries is a list of numeric IDs that
	// should be considered to be of this tag-type.
	Entries []int
}

// ClientboundTags sends a collection of IDs by property
// to the client. An example is a list of IDs that are considered
// blocks, where the ID is the tag and the list is a list
// of numeric IDs.
type ClientboundTags struct {
	// BlockTags are the block tags that the game uses.
	BlockTags []Tag
	// ItemTags are the item tags that the game uses.
	ItemTags []Tag
	// FluidTags are the fluid tags that the game uses.
	FluidTags []Tag
	// EntityTags are the entity tags that the game uses.
	EntityTags []Tag
}

// ID returns the constant packet ID.
func (ClientboundTags) ID() ID { return IDClientboundTags }

// Name returns the constant packet name.
func (ClientboundTags) Name() string { return "Tags" }

// EncodeInto writes this packet into the given writer.
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
