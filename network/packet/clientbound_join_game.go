package packet

import (
	"io"
	"reflect"
	"strconv"

	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/id"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundJoinGame{}))
}

type ClientboundJoinGame struct {
	EntityID            int32
	Hardcore            bool
	Gamemode            int
	PreviousGamemode    int
	WorldNames          []id.ID
	DimensionCodec      nbt.Tag
	Dimension           nbt.Tag
	WorldName           id.ID
	HashedSeed          int64
	MaxPlayers          int
	ViewDistance        int
	ReducedDebugInfo    bool
	EnableRespawnScreen bool
	Debug               bool
	Flat                bool
}

func (ClientboundJoinGame) ID() ID       { return IDClientboundJoinGame }
func (ClientboundJoinGame) Name() string { return "Join Game" }

func (c ClientboundJoinGame) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteInt("entity ID", c.EntityID)
	enc.WriteBoolean("is hardcore", c.Hardcore)
	enc.WriteUbyte("gamemode", uint8(c.Gamemode))
	enc.WriteByte("previous gamemode", int8(c.PreviousGamemode))
	enc.WriteVarInt("world count", len(c.WorldNames))
	for i, worldName := range c.WorldNames {
		enc.WriteID("world names["+strconv.Itoa(i)+"]", worldName)
	}
	enc.WriteNBT("dimension codec", c.DimensionCodec)
	enc.WriteNBT("dimension", c.Dimension)
	enc.WriteID("world name", c.WorldName)
	enc.WriteLong("hashed seed", c.HashedSeed)
	enc.WriteVarInt("max players", c.MaxPlayers)
	enc.WriteVarInt("view distance", c.ViewDistance)
	enc.WriteBoolean("reduced debug info", c.ReducedDebugInfo)
	enc.WriteBoolean("enable respawn screen", c.EnableRespawnScreen)
	enc.WriteBoolean("is debug", c.Debug)
	enc.WriteBoolean("is flat", c.Flat)

	return
}
