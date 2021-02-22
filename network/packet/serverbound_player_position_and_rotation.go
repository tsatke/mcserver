package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ServerboundPlayerPositionAndRotation{}))
}

// ServerboundPlayerPositionAndRotation is a combination packet of
// ServerboundPlayerRotation and ServerboundPlayerPosition.
type ServerboundPlayerPositionAndRotation struct {
	// X is the absolute X position of the player.
	X float64
	// FeetY is the absolute Y position of the player. This is usually HeadY - 1.62.
	FeetY float64
	// Z is the absolute Z position of the player.
	Z float64
	// Yaw is the absolute rotation on the X axis in degrees. Yaw is not clamped
	// to [0,360].
	Yaw float32
	// Pitch is a value in [-90,90], where -90 is to be interpreted as looking
	// straight up, 0 is looking straight, and 90 is looking straight down.
	Pitch    float32
	OnGround bool
}

// ID returns the constant packet ID.
func (ServerboundPlayerPositionAndRotation) ID() ID { return IDServerboundPlayerPositionAndRotation }

// Name returns the constant packet name.
func (ServerboundPlayerPositionAndRotation) Name() string { return "Player Position and Rotation" }

// DecodeFrom will fill this struct with values read from the given reader.
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
