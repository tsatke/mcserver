package types

import (
	"fmt"
	"io"
)

type UnsignedByte uint8

func NewUnsignedByte(v uint8) *UnsignedByte {
	ub := UnsignedByte(v)
	return &ub
}

func (b *UnsignedByte) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, UnsignedByteSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*b = UnsignedByte(buf[0])
	return nil
}

func (b UnsignedByte) EncodeInto(w io.Writer) error {
	buf := make([]byte, UnsignedByteSize)
	buf[0] = uint8(b)
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != UnsignedByteSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", UnsignedByteSize, n)
	}
	return nil
}
