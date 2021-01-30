package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ServerboundPlayerPositionAndLook{}))
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

	dec := decoder{rd}

	s.X = dec.readDouble("x")
	s.Y = dec.readDouble("y")
	s.Z = dec.readDouble("z")
	s.Yaw = dec.readFloat("yaw")
	s.Pitch = dec.readFloat("pitch")
	s.Flags = dec.readByte("flags")
	s.TeleportID = dec.readVarInt("teleport id")

	return
}
