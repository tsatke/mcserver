package types

import (
	"fmt"
	"io"
)

type Int int32

func NewInt(value int32) *Int {
	intVal := Int(value)
	return &intVal
}

func (i *Int) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, IntSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*i = Int(ByteOrder.Uint32(buf))
	return nil
}

func (i Int) EncodeInto(w io.Writer) error {
	buf := make([]byte, IntSize)
	ByteOrder.PutUint32(buf, uint32(i))
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != IntSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", IntSize, n)
	}
	return nil
}
