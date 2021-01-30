package types

import (
	"fmt"
	"io"
)

type Position struct {
	X, Y, Z int
}

func (p *Position) DecodeFrom(rd io.Reader) error {
	buf := make([]byte, PositionSize)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	val := ByteOrder.Uint64(buf)

	(*p).X = int(val >> 38)
	(*p).Y = int(val & 0xfff)
	(*p).Z = int(val << 26 >> 38)
	return nil
}

func (p Position) EncodeInto(w io.Writer) error {
	val := ((uint64(p.X) & 0x3FFFFFF) << 38) | ((uint64(p.Z) & 0x3FFFFFF) << 12) | (uint64(p.Y) & 0xFFF)
	buf := make([]byte, PositionSize)
	ByteOrder.PutUint64(buf, val)
	n, err := w.Write(buf)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != len(buf) {
		return fmt.Errorf("need to write %d bytes, but wrote %d", len(buf), n)
	}
	return nil
}
