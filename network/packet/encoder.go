package packet

import (
	"io"

	"github.com/google/uuid"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/network/packet/types"
)

type encoder struct {
	w io.Writer
}

func (e encoder) writeVarInt(fieldName string, val int) {
	_write(e, fieldName, types.NewVarInt(val))
}

func (e encoder) writeInt(fieldName string, val int32) {
	_write(e, fieldName, types.NewInt(val))
}

func (e encoder) writeString(fieldName, s string) {
	_write(e, fieldName, types.NewString(s))
}

func (e encoder) writeUshort(fieldName string, val uint16) {
	_write(e, fieldName, types.NewUnsignedShort(val))
}

func (e encoder) writeByte(fieldName string, val int8) {
	_write(e, fieldName, types.NewByte(val))
}

func (e encoder) writeUbyte(fieldName string, val uint8) {
	_write(e, fieldName, types.NewUnsignedByte(val))
}

func (e encoder) writeBoolean(fieldName string, val bool) {
	_write(e, fieldName, types.NewBoolean(val))
}

func (e encoder) writeLong(fieldName string, val int64) {
	_write(e, fieldName, types.NewLong(val))
}

func (e encoder) writeUUID(fieldName string, val uuid.UUID) {
	_write(e, fieldName, types.NewUUID(val))
}

func (e encoder) writeChat(fieldName string, val types.Chat) {
	_write(e, fieldName, &val)
}

func (e encoder) writeID(fieldName string, val id.ID) {
	_write(e, fieldName, types.NewString(val.String()))
}

func (e encoder) writeNBT(fieldName string, val nbt.Tag) {
	_write(e, fieldName, &types.NBT{
		Tag: val,
	})
}

func _write(e encoder, fieldName string, v types.Value) {
	panicIffErr(fieldName, v.EncodeInto(e.w))
}
