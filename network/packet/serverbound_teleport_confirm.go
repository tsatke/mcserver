package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(StatePlay, reflect.TypeOf(ServerboundTeleportConfirm{}))
}

type ServerboundTeleportConfirm struct {
	TeleportID int
}

func (ServerboundTeleportConfirm) ID() ID       { return IDServerboundTeleportConfirm }
func (ServerboundTeleportConfirm) Name() string { return "Teleport Confirm" }

func (s *ServerboundTeleportConfirm) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.TeleportID = dec.ReadVarInt("teleport id")

	return
}
