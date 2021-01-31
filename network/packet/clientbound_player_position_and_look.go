package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(StatePlay, reflect.TypeOf(ClientboundPlayerPositionAndLook{}))
}

type ClientboundPlayerPositionAndLook struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	Flags      int8
	TeleportID int
}

func (ClientboundPlayerPositionAndLook) ID() ID       { return IDClientboundPlayerPositionAndLook }
func (ClientboundPlayerPositionAndLook) Name() string { return "Player Position and Look" }

func (c ClientboundPlayerPositionAndLook) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteDouble("x", c.X)
	enc.WriteDouble("y", c.Y)
	enc.WriteDouble("z", c.Z)
	enc.WriteFloat("yaw", c.Yaw)
	enc.WriteFloat("pitch", c.Pitch)
	enc.WriteByte("flags", c.Flags)
	enc.WriteVarInt("teleport id", c.TeleportID)

	return
}
