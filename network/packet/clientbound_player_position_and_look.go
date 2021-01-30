package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundPlayerPositionAndLook{}))
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

	enc := encoder{w}

	enc.writeDouble("x", c.X)
	enc.writeDouble("y", c.Y)
	enc.writeDouble("z", c.Z)
	enc.writeFloat("yaw", c.Yaw)
	enc.writeFloat("pitch", c.Pitch)
	enc.writeByte("flags", c.Flags)
	enc.writeVarInt("teleport id", c.TeleportID)

	return
}
