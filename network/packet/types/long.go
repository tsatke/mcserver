package types

import (
	"fmt"
	"io"
)

type Long int64

func NewLong(value int64) *Long {
	longVal := Long(value)
	return &longVal
}

func (l *Long) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, LongSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*l = Long(ByteOrder.Uint64(buf))
	return nil
}

func (l Long) EncodeInto(w io.Writer) error {
	buf := make([]byte, LongSize)
	ByteOrder.PutUint64(buf, uint64(l))
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != LongSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", LongSize, n)
	}
	return nil
}
