package packet

import (
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/google/uuid"
	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/chat"
	"github.com/tsatke/mcserver/game/id"
)

type Encoder struct {
	W io.Writer
}

func (e Encoder) WriteVarInt(fieldName string, val int) {
	value := uint32(val)
	buf := make([]byte, 0)
	for {
		tmp := byte(value & 0b01111111)
		value >>= 7
		if value != 0 {
			tmp |= 0b10000000
		}
		buf = append(buf, tmp)
		if value == 0 {
			break
		}
	}
	_write(e.W, fieldName, buf[:])
}

func (e Encoder) WriteInt(fieldName string, val int32) {
	var buf [IntSize]byte
	ByteOrder.PutUint32(buf[:], uint32(val))
	_write(e.W, fieldName, buf[:])
}

func (e Encoder) WriteString(fieldName, s string) {
	e.WriteVarInt(fieldName+" string length", len(s))
	_write(e.W, fieldName, []byte(s))
}

func (e Encoder) WriteUshort(fieldName string, val uint16) {
	var buf [UnsignedShortSize]byte
	ByteOrder.PutUint16(buf[:], val)
	_write(e.W, fieldName, buf[:])
}

func (e Encoder) WriteByte(fieldName string, val int8) {
	_write(e.W, fieldName, []byte{byte(val)})
}

func (e Encoder) WriteByteArray(fieldName string, val []byte) {
	_write(e.W, fieldName, val)
}

func (e Encoder) WriteUbyte(fieldName string, val uint8) {
	_write(e.W, fieldName, []byte{val})
}

func (e Encoder) WriteBoolean(fieldName string, val bool) {
	if val {
		_write(e.W, fieldName, []byte{0x01})
	} else {
		_write(e.W, fieldName, []byte{0x00})
	}
}

func (e Encoder) WriteLong(fieldName string, val int64) {
	var buf [LongSize]byte
	ByteOrder.PutUint64(buf[:], uint64(val))
	_write(e.W, fieldName, buf[:])
}

func (e Encoder) WriteDouble(fieldName string, val float64) {
	var buf [DoubleSize]byte
	ByteOrder.PutUint64(buf[:], math.Float64bits(val))
	_write(e.W, fieldName, buf[:])
}

func (e Encoder) WriteFloat(fieldName string, val float32) {
	var buf [FloatSize]byte
	ByteOrder.PutUint32(buf[:], math.Float32bits(val))
	_write(e.W, fieldName, buf[:])
}

func (e Encoder) WriteUUID(fieldName string, val uuid.UUID) {
	_write(e.W, fieldName, val[:])
}

func (e Encoder) WriteChat(fieldName string, val chat.Chat) {
	data, err := json.Marshal(val)
	panicIffErr(fieldName, err)
	e.WriteString(fieldName, string(data))
}

func (e Encoder) WriteID(fieldName string, val id.ID) {
	e.WriteString(fieldName, val.String())
}

func (e Encoder) WriteNBT(fieldName string, val nbt.Tag) {
	enc := nbt.NewEncoder(e.W, ByteOrder)
	panicIffErr(fieldName, enc.WriteTag(val))
}

func _write(w io.Writer, fieldName string, buf []byte) {
	n, err := w.Write(buf)
	if err != nil {
		panicIffErr(fieldName, fmt.Errorf("write: %w", err))
	}
	if n != len(buf) {
		panicIffErr(fieldName, fmt.Errorf("need to write %d bytes, but wrote %d", len(buf), n))
	}
}
