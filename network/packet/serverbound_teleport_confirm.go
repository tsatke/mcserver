package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ServerboundTeleportConfirm{}))
}

type ServerboundTeleportConfirm struct {
	TeleportID int
}

func (ServerboundTeleportConfirm) ID() ID       { return IDServerboundTeleportConfirm }
func (ServerboundTeleportConfirm) Name() string { return "Teleport Confirm" }

func (s *ServerboundTeleportConfirm) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := decoder{rd}

	s.TeleportID = dec.readVarInt("teleport id")

	return
}
