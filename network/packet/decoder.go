package packet

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/google/uuid"

	"github.com/tsatke/mcserver/game/id"
)

// Decoder is a struct that can decode protocol values from a reader.
// Please note, that for concise and clean API reasons, all Read methods
// fail by panicking an error. Avoid use outside of the packet package if possible.
// If you need to use this, make sure that you recover any errors.
type Decoder struct {
	Rd io.Reader
}

func (d Decoder) ReadVarInt(fieldName string) int {
	var res int32
	var readCnt int

	var buf [1]byte
	for {
		_, err := io.ReadFull(d.Rd, buf[:])
		panicIffErr(fieldName, err)

		val := buf[0] & (1<<7 - 1)
		res |= int32(val) << (7 * readCnt)
		readCnt++
		if readCnt > 5 {
			panicIffErr(fieldName, fmt.Errorf("VarInt too big"))
		}

		if buf[0]&(1<<7) == 0 {
			break
		}
	}

	return int(res)
}

func (d Decoder) ReadString(fieldName string) string {
	strLen := d.ReadVarInt(fieldName + " length")

	/*
		This way of reading the string may be slower than some of
		the other methods seen in this file. However, this doesn't
		allocate memory according to some bytes a user sends.
		If we don't read a string like this, a user could send an
		(invalid) string with a length of 2GiB, but no payload.
		The reading will error out due to EOF, however, the 2GiB
		would still have been allocated, which may cause the server
		to lag pretty hard. This way, the user kind of has to prove
		the length by actually sending the payload.
	*/
	var buf bytes.Buffer
	n, err := buf.ReadFrom(io.LimitReader(d.Rd, int64(strLen)))
	panicIffErr(fieldName, err)
	if n != int64(strLen) {
		panicIffErr(fieldName, fmt.Errorf("prefix indicated length of %d, but only got %d bytes", strLen, n))
	}

	return buf.String()
}

func (d Decoder) ReadUbyte(fieldName string) byte {
	var buf [ByteSize]byte
	_, err := io.ReadFull(d.Rd, buf[:])
	panicIffErr(fieldName, err)

	return buf[0]
}

func (d Decoder) ReadUshort(fieldName string) uint16 {
	var buf [UnsignedShortSize]byte
	_, err := io.ReadFull(d.Rd, buf[:])
	panicIffErr(fieldName, err)

	return ByteOrder.Uint16(buf[:])
}

func (d Decoder) ReadByte(fieldName string) int8 {
	return int8(d.ReadUbyte(fieldName))
}

func (d Decoder) ReadBoolean(fieldName string) bool {
	var buf [BooleanSize]byte
	_, err := io.ReadFull(d.Rd, buf[:])
	panicIffErr(fieldName, err)

	return buf[0] == 1
}

func (d Decoder) ReadUUID(fieldName string) uuid.UUID {
	var uuid uuid.UUID
	_, err := io.ReadFull(d.Rd, uuid[:])
	panicIffErr(fieldName, err)
	return uuid
}

func (d Decoder) ReadLong(fieldName string) int64 {
	var buf [LongSize]byte
	_, err := io.ReadFull(d.Rd, buf[:])
	panicIffErr(fieldName, err)
	return int64(ByteOrder.Uint64(buf[:]))
}

func (d Decoder) ReadDouble(fieldName string) float64 {
	var buf [DoubleSize]byte
	_, err := io.ReadFull(d.Rd, buf[:])
	panicIffErr(fieldName, err)

	return math.Float64frombits(ByteOrder.Uint64(buf[:]))
}

func (d Decoder) ReadFloat(fieldName string) float32 {
	var buf [FloatSize]byte
	_, err := io.ReadFull(d.Rd, buf[:])
	panicIffErr(fieldName, err)

	return math.Float32frombits(ByteOrder.Uint32(buf[:]))
}

func (d Decoder) ReadID(fieldName string) id.ID {
	return id.ParseID(d.ReadString(fieldName))
}
