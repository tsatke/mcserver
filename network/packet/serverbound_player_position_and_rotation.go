package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ServerboundPlayerPositionAndRotation{}))
}

type ServerboundPlayerPositionAndRotation struct {
	X, FeetY, Z float64
	Yaw, Pitch  float32
	OnGround    bool
}

func (ServerboundPlayerPositionAndRotation) ID() ID       { return IDServerboundPlayerPositionAndRotation }
func (ServerboundPlayerPositionAndRotation) Name() string { return "Player Position and Rotation" }

func (s *ServerboundPlayerPositionAndRotation) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.X = dec.ReadDouble("x")
	s.FeetY = dec.ReadDouble("feet y")
	s.Z = dec.ReadDouble("z")
	s.Yaw = dec.ReadFloat("yaw")
	s.Pitch = dec.ReadFloat("pitch")
	s.OnGround = dec.ReadBoolean("on ground")

	return
}
