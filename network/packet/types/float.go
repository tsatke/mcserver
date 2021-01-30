package types

import (
	"fmt"
	"io"
	"math"
)

type Float float32

func NewFloat(value float32) *Float {
	longVal := Float(value)
	return &longVal
}

func (l *Float) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, FloatSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*l = Float(math.Float32frombits(ByteOrder.Uint32(buf)))
	return nil
}

func (l Float) EncodeInto(w io.Writer) error {
	buf := make([]byte, FloatSize)
	ByteOrder.PutUint32(buf, math.Float32bits(float32(l)))
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != FloatSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", FloatSize, n)
	}
	return nil
}
