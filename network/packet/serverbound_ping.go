package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StateStatus, reflect.TypeOf(ServerboundPing{}))
}

type ServerboundPing struct {
	Payload int64
}

func (ServerboundPing) ID() ID       { return IDServerboundPing }
func (ServerboundPing) Name() string { return "Ping" }

func (s *ServerboundPing) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := decoder{rd}

	s.Payload = dec.readLong("payload")

	return
}
