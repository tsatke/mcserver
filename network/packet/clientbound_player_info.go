package packet

import (
	"io"
	"reflect"

	"github.com/google/uuid"

	"github.com/tsatke/mcserver/game/chat"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundPlayerInfo{}))
}

type PlayerInfoPlayer struct {
	UUID uuid.UUID

	Name           string
	Properties     []PlayerInfoProperties
	Gamemode       int
	Ping           int
	HasDisplayName bool
	DisplayName    chat.Chat
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
func (ClientboundPlayerInfo) Name() string { return "Player Info" }

func (c ClientboundPlayerInfo) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteVarInt("action", int(c.Action))
	enc.WriteVarInt("number of players", len(c.Players))
	for _, player := range c.Players {
		enc.WriteUUID("uuid", player.UUID)
		switch c.Action {
		case 0:
			enc.WriteString("name", player.Name)
			enc.WriteVarInt("number of properties", len(player.Properties))
			for _, property := range player.Properties {
				enc.WriteString("property name", property.Name)
				enc.WriteString("property value", property.Value)
				enc.WriteBoolean("is signed", property.Signed)
				if property.Signed {
					enc.WriteString("signature", property.Signature)
				}
			}
			enc.WriteVarInt("gamemode", int(player.Gamemode))
			enc.WriteVarInt("ping", player.Ping)
			enc.WriteBoolean("has display name", player.HasDisplayName)
			if player.HasDisplayName {
				enc.WriteChat("display name", player.DisplayName)
			}
		case 1:
			enc.WriteVarInt("gamemode", int(player.Gamemode))
		case 2:
			enc.WriteVarInt("ping", player.Ping)
		case 3:
			enc.WriteBoolean("has display name", player.HasDisplayName)
			if player.HasDisplayName {
				enc.WriteChat("display name", player.DisplayName)
			}
		default:
			// write no fields for 4
		}
	}

	return
}
