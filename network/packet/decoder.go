package packet

import (
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/network/packet/types"
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

func panicIffErr(fieldName string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", fieldName, err))
	}
}

type decoder struct {
	rd io.Reader
}

func (d decoder) readVarInt(fieldName string) int {
	val := types.NewVarInt(0)
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return int(*val)
}

func (d decoder) readString(fieldName string) string {
	val := types.NewString("")
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return string(*val)
}

func (d decoder) readUbyte(fieldName string) byte {
	val := types.NewUnsignedByte(0)
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return byte(*val)
}

func (d decoder) readUshort(fieldName string) uint16 {
	val := types.NewUnsignedShort(0)
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return uint16(*val)
}

func (d decoder) readByte(fieldName string) int8 {
	val := types.NewByte(0)
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return int8(*val)
}

func (d decoder) readBoolean(fieldName string) bool {
	val := types.NewBoolean(false)
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return bool(*val)
}

func (d decoder) readLong(fieldName string) int64 {
	val := types.NewLong(0)
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return int64(*val)
}

func (d decoder) readUUID(fieldName string) uuid.UUID {
	val := types.NewUUID([16]byte{})
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return uuid.UUID(*val)
}

func (d decoder) readID(fieldName string) id.ID {
	val := types.String("")
	panicIffErr(fieldName, val.DecodeFrom(d.rd))
	return id.ParseID(string(val))
}
