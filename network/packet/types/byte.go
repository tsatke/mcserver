package types

import (
	"fmt"
	"io"
)

type Byte int8

func NewByte(v int8) *Byte {
	b := Byte(v)
	return &b
}

func (b *Byte) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, ByteSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*b = Byte(buf[0])
	return nil
}

func (b Byte) EncodeInto(w io.Writer) error {
	buf := make([]byte, ByteSize)
	buf[0] = uint8(b)
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != ByteSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", ByteSize, n)
	}
	return nil
}
