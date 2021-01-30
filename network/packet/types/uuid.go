package types

import (
	"fmt"
	"io"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func NewUUID(uuid uuid.UUID) *UUID {
	u := UUID(uuid)
	return &u
}

func (u *UUID) DecodeFrom(rd io.Reader) error {
	var uuid uuid.UUID
	_, err := io.ReadFull(rd, uuid[:])
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	*u = UUID(uuid)
	return nil
}

func (u UUID) EncodeInto(w io.Writer) error {
	n, err := w.Write(u[:])
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != UUIDSize {
		return fmt.Errorf("need to write %d bytes, but wrote %d", UUIDSize, n)
	}
	return nil
}
