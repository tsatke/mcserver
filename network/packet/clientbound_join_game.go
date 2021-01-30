package packet

import (
	"io"
	"reflect"
	"strconv"

	"github.com/tsatke/nbt"

	"github.com/tsatke/mcserver/game/id"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundJoinGame{}))
}

type Gamemode int8

const (
	GamemodeUnknown Gamemode = iota - 1
	GamemodeSurvival
	GamemodeCreative
	GamemodeAdventure
	GamemodeSpectator
)

type ClientboundJoinGame struct {
	EntityID            int32
	Hardcore            bool
	Gamemode            Gamemode
	PreviousGamemode    Gamemode
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

	enc := encoder{w}

	enc.writeInt("entity ID", c.EntityID)
	enc.writeBoolean("is hardcore", c.Hardcore)
	enc.writeUbyte("gamemode", uint8(c.Gamemode))
	enc.writeByte("previous gamemode", int8(c.PreviousGamemode))
	enc.writeVarInt("world count", len(c.WorldNames))
	for i, worldName := range c.WorldNames {
		enc.writeID("world names["+strconv.Itoa(i)+"]", worldName)
	}
	enc.writeNBT("dimension codec", c.DimensionCodec)
	enc.writeNBT("dimension", c.Dimension)
	enc.writeID("world name", c.WorldName)
	enc.writeLong("hashed seed", c.HashedSeed)
	enc.writeVarInt("max players", c.MaxPlayers)
	enc.writeVarInt("view distance", c.ViewDistance)
	enc.writeBoolean("reduced debug info", c.ReducedDebugInfo)
	enc.writeBoolean("enable respawn screen", c.EnableRespawnScreen)
	enc.writeBoolean("is debug", c.Debug)
	enc.writeBoolean("is flat", c.Flat)

	return
}
