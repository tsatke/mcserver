package types

import (
	"fmt"
	"io"
)

type Boolean bool

func NewBoolean(value bool) *Boolean {
	boolean := Boolean(value)
	return &boolean
}

func (b *Boolean) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, BooleanSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*b = buf[0] == 1
	return nil
}

func (b Boolean) EncodeInto(w io.Writer) error {
	buf := make([]byte, BooleanSize)
	if b {
		buf[0] = 1
	}
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != BooleanSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", BooleanSize, n)
	}
	return nil
}
