package types

import (
	"fmt"
	"io"
	"math"
)

type Double float64

func NewDouble(value float64) *Double {
	longVal := Double(value)
	return &longVal
}

func (l *Double) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, DoubleSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*l = Double(math.Float64frombits(ByteOrder.Uint64(buf)))
	return nil
}

func (l Double) EncodeInto(w io.Writer) error {
	buf := make([]byte, DoubleSize)
	ByteOrder.PutUint64(buf, math.Float64bits(float64(l)))
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != DoubleSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", DoubleSize, n)
	}
	return nil
}
