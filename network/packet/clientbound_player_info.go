package packet

import (
	"io"
	"reflect"

	"github.com/google/uuid"

	"github.com/tsatke/mcserver/network/packet/types"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundPlayerInfo{}))
}

type PlayerInfoPlayer struct {
	UUID uuid.UUID

	Name           string
	Properties     []PlayerInfoProperties
	Gamemode       Gamemode
	Ping           int
	HasDisplayName bool
	DisplayName    types.Chat
}

type PlayerInfoProperties struct {
	Name      string
	Value     string
	Signed    bool
	Signature string
}

type PlayerInfoAction int

const (
	PlayerInfoActionAddPlayer PlayerInfoAction = iota
	PlayerInfoActionUpdateGamemode
	PlayerInfoUpdateLatency
	PlayerInfoUpdateDisplayName
	PlayerInfoRemovePlayer
)

type ClientboundPlayerInfo struct {
	Action  PlayerInfoAction
	Players []PlayerInfoPlayer
}

func (ClientboundPlayerInfo) ID() ID       { return IDClientboundPlayerInfo }
func (ClientboundPlayerInfo) Name() string { return "PlayerInfo" }

func (c ClientboundPlayerInfo) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeVarInt("action", int(c.Action))
	enc.writeVarInt("number of players", len(c.Players))
	for _, player := range c.Players {
		enc.writeUUID("uuid", player.UUID)
		switch c.Action {
		case 0:
			enc.writeString("name", player.Name)
			enc.writeVarInt("number of properties", len(player.Properties))
			for _, property := range player.Properties {
				enc.writeString("property name", property.Name)
				enc.writeString("property value", property.Value)
				enc.writeBoolean("is signed", property.Signed)
				if property.Signed {
					enc.writeString("signature", property.Signature)
				}
			}
			enc.writeVarInt("gamemode", int(player.Gamemode))
			enc.writeVarInt("ping", player.Ping)
			enc.writeBoolean("has display name", player.HasDisplayName)
			if player.HasDisplayName {
				enc.writeChat("display name", player.DisplayName)
			}
		case 1:
			enc.writeVarInt("gamemode", int(player.Gamemode))
		case 2:
			enc.writeVarInt("ping", player.Ping)
		case 3:
			enc.writeBoolean("has display name", player.HasDisplayName)
			if player.HasDisplayName {
				enc.writeChat("display name", player.DisplayName)
			}
		default:
			// write no fields for 4
		}
	}

	return
}
