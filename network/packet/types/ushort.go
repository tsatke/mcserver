package types

import (
	"fmt"
	"io"
)

type UnsignedShort uint16

func NewUnsignedShort(value uint16) *UnsignedShort {
	ushortVal := UnsignedShort(value)
	return &ushortVal
}

func (s *UnsignedShort) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, UnsignedShortSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*s = UnsignedShort(ByteOrder.Uint16(buf))
	return nil
}

func (s UnsignedShort) EncodeInto(w io.Writer) error {
	buf := make([]byte, UnsignedShortSize)
	ByteOrder.PutUint16(buf, uint16(s))
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != UnsignedShortSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", UnsignedShortSize, n)
	}
	return nil
}
