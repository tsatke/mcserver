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

// Encoder is a decorating struct, which uses specialized algorithms to
// write data into an underlying io.Writer.
type Encoder struct {
	// W is the underlying writer into which data will eventually
	// be written.
	W io.Writer
}

// WriteVarInt writes the given int into the writer as a VarInt.
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
	_write(e.W, fieldName, buf)
}

// WriteInt writes the given int32 with ByteOrder into the writer.
func (e Encoder) WriteInt(fieldName string, val int32) {
	var buf [IntSize]byte
	ByteOrder.PutUint32(buf[:], uint32(val))
	_write(e.W, fieldName, buf[:])
}

// WriteString writes a VarInt into the writer, indicating the length of the given
// string. After that, the string is written as byte array.
// See Encoder.WriteVarInt.
func (e Encoder) WriteString(fieldName, s string) {
	e.WriteVarInt(fieldName+" string length", len(s))
	_write(e.W, fieldName, []byte(s))
}

// WriteUshort writes the given uint16 with ByteOrder into the writer.
func (e Encoder) WriteUshort(fieldName string, val uint16) {
	var buf [UnsignedShortSize]byte
	ByteOrder.PutUint16(buf[:], val)
	_write(e.W, fieldName, buf[:])
}

// WriteByte writes the given int8 into the writer as unsigned value.
func (e Encoder) WriteByte(fieldName string, val int8) {
	_write(e.W, fieldName, []byte{byte(val)})
}

// WriteByteArray writes the given byte array into the writer.
func (e Encoder) WriteByteArray(fieldName string, val []byte) {
	_write(e.W, fieldName, val)
}

// WriteUbyte writes the given byte into the writer.
func (e Encoder) WriteUbyte(fieldName string, val uint8) {
	_write(e.W, fieldName, []byte{val})
}

// WriteBoolean writes a single byte into the writer, 0x01 for true, 0x00 for false.
func (e Encoder) WriteBoolean(fieldName string, val bool) {
	if val {
		_write(e.W, fieldName, []byte{0x01})
	} else {
		_write(e.W, fieldName, []byte{0x00})
	}
}

// WriteLong writes the given int64 with ByteOrder into the writer.
func (e Encoder) WriteLong(fieldName string, val int64) {
	var buf [LongSize]byte
	ByteOrder.PutUint64(buf[:], uint64(val))
	_write(e.W, fieldName, buf[:])
}

// WriteDouble writes the given float64 with ByteOrder into the writer
// as IEEE754 encoded value.
func (e Encoder) WriteDouble(fieldName string, val float64) {
	var buf [DoubleSize]byte
	ByteOrder.PutUint64(buf[:], math.Float64bits(val))
	_write(e.W, fieldName, buf[:])
}

// WriteFloat writes the given float32 with ByteOrder into the writer
// as IEEE754 encoded value.
func (e Encoder) WriteFloat(fieldName string, val float32) {
	var buf [FloatSize]byte
	ByteOrder.PutUint32(buf[:], math.Float32bits(val))
	_write(e.W, fieldName, buf[:])
}

// WriteUUID writes the 16 bytes of the UUID as plain bytes into the writer.
func (e Encoder) WriteUUID(fieldName string, val uuid.UUID) {
	_write(e.W, fieldName, val[:])
}

// WriteChat writes the json marshalled value of the given chat object
// into the writer as a string.
func (e Encoder) WriteChat(fieldName string, val chat.Chat) {
	data, err := json.Marshal(val)
	panicIffErr(fieldName, err)
	e.WriteString(fieldName, string(data))
}

// WriteID writes the given ID into the writer as a string.
// See Encoder.WriteString.
func (e Encoder) WriteID(fieldName string, val id.ID) {
	e.WriteString(fieldName, val.String())
}

// WriteNBT writes an nbt tag into the writer.
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
