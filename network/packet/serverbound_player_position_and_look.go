package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ServerboundPlayerPositionAndLook{}))
}

type ServerboundPlayerPositionAndLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      int8
	TeleportID int
}

func (ServerboundPlayerPositionAndLook) ID() ID       { return IDServerboundPlayerPositionAndLook }
func (ServerboundPlayerPositionAndLook) Name() string { return "Player Position and Look" }

func (s *ServerboundPlayerPositionAndLook) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.X = dec.ReadDouble("x")
	s.Y = dec.ReadDouble("y")
	s.Z = dec.ReadDouble("z")
	s.Yaw = dec.ReadFloat("yaw")
	s.Pitch = dec.ReadFloat("pitch")
	s.Flags = dec.ReadByte("flags")
	s.TeleportID = dec.ReadVarInt("teleport id")

	return
}
