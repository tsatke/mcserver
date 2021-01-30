package types

import (
	"encoding/binary"
	"io"
)

var (
	// ByteOrder is the byte order that is used by this server.
	ByteOrder = binary.BigEndian
)

type Value interface {
	DecodeFrom(io.Reader) error
	EncodeInto(io.Writer) error
}
