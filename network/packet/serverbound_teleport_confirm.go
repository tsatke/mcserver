package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ServerboundTeleportConfirm{}))
}

// ServerboundTeleportConfirm is sent by the client as confirmation of a ClientboundPlayerPositionAndLook.
type ServerboundTeleportConfirm struct {
	TeleportID int
}

// ID returns the constant packet ID.
func (ServerboundTeleportConfirm) ID() ID { return IDServerboundTeleportConfirm }

// Name returns the constant packet name.
func (ServerboundTeleportConfirm) Name() string { return "Teleport Confirm" }

// DecodeFrom will fill this struct with values read from the given reader.
func (s *ServerboundTeleportConfirm) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.TeleportID = dec.ReadVarInt("teleport id")

	return
}
