package types

import (
	"fmt"
	"io"
)

type String string

func NewString(value string) *String {
	str := String(value)
	return &str
}

func (s *String) DecodeFrom(rd io.Reader) error {
	var strLen VarInt
	if err := strLen.DecodeFrom(rd); err != nil {
		return fmt.Errorf("decode string length: %w", err)
	}

	buf := make([]byte, strLen)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	*s = String(buf)
	return nil
}

func (s String) EncodeInto(w io.Writer) error {
	if err := VarInt(len(s)).EncodeInto(w); err != nil {
		return fmt.Errorf("encode string length: %w", err)
	}

	n, err := io.WriteString(w, string(s))
	if err != nil {
		return fmt.Errorf("write string: %w", err)
	}
	if n != len(s) {
		return fmt.Errorf("need to write %d, only wrote %d", len(s), n)
	}
	return nil
}
