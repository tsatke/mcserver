package types

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/tsatke/nbt"
)

type NBT struct {
	nbt.Tag
}

func NewNBT(tag nbt.Tag) *NBT {
	return &NBT{
		Tag: tag,
	}
}

func (n *NBT) DecodeFrom(rd io.Reader) error {
	dec := nbt.NewDecoder(rd, binary.BigEndian)
	tag, err := dec.ReadTag()
	if err != nil {
		return fmt.Errorf("read tag: %w", err)
	}
	*n = NBT{
		Tag: tag,
	}
	return nil
}

func (n NBT) EncodeInto(w io.Writer) error {
	return nbt.NewEncoder(w, binary.BigEndian).WriteTag(n.Tag)
}
