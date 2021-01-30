package types

import (
	"fmt"
	"io"
)

type VarInt int32

func NewVarInt(value int) *VarInt {
	varInt := VarInt(value)
	return &varInt
}

func (i *VarInt) DecodeFrom(rd io.Reader) error {
	var res int32
	var readCnt int

	buf := make([]byte, 1)
	for {
		_, err := io.ReadFull(rd, buf)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		val := buf[0] & 0b01111111
		res |= int32(val) << (7 * readCnt)
		readCnt++
		if readCnt > 5 {
			return fmt.Errorf("VarInt too big")
		}

		if buf[0]&0b10000000 == 0 {
			break
		}
	}

	*i = VarInt(res)
	return nil
}

func (i VarInt) EncodeInto(w io.Writer) error {
	value := uint32(i)
	buf := make([]byte, 0, VarIntMaxSize)
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
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != len(buf) {
		return fmt.Errorf("need to write %d bytes, but wrote %d", len(buf), n)
	}
	return nil
}

func (i VarInt) Len() int {
	// TODO: optimize
	value := uint32(i)
	buf := make([]byte, 0, VarIntMaxSize)
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
	return len(buf)
}
