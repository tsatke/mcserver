package types

import (
	"fmt"
	"io"
)

type ByteArray []byte

func NewByteArray(v []byte) *ByteArray {
	b := ByteArray(v)
	return &b
}

func (b *ByteArray) DecodeFrom(rd io.Reader) error {
	_, err := io.ReadFull(rd, *b)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func (b ByteArray) EncodeInto(w io.Writer) error {
	n, err := w.Write(b)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != ByteSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", ByteSize, n)
	}
	return nil
}
