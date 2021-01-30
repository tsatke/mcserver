package game

import (
	"bytes"

	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/network/packet"
	"github.com/tsatke/mcserver/network/packet/types"
)

func (g *Game) processPacket(source *Player, pkg packet.Serverbound) {
	source.Lock()
	defer source.Unlock()

	switch p := pkg.(type) {
	case *packet.ServerboundPluginMessage:
		g.processServerboundPluginMessage(source, p)
	case *packet.ServerboundClientSettings:
		source.client.settings.locale = p.Locale
		source.client.settings.viewDistance = p.ViewDistance
	default:
		g.log.Warn().
			Str("name", pkg.Name()).
			Msg("unhandled packet")
	}
}

func (g *Game) processServerboundPluginMessage(source *Player, p *packet.ServerboundPluginMessage) {
	switch p.Channel {
	case id.ParseID("minecraft:brand"):
		rd := bytes.NewReader(p.Data)
		// brand is prefixed with a varint, indicating its length, which we discard, since we know the
		// length of the string plus the size of the varint
		var brandLen types.VarInt
		_ = brandLen.DecodeFrom(rd)
		source.client.brand = string(p.Data[brandLen.Len():])
	default:
		g.log.Debug().
			Stringer("channel", p.Channel).
			Msg("ignoring plugin message")
	}
}
