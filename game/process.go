package game

import (
	"bytes"

	"github.com/tsatke/mcserver/game/id"
	"github.com/tsatke/mcserver/network/packet"
)

// processPacket checks the given packet and delegates it to the appropriate processing method.
// While the packet is processed, this method holds the lock on the given player.
func (g *Game) processPacket(source *Player, pkg packet.Serverbound) {
	source.Lock()
	defer source.Unlock()

	switch p := pkg.(type) {
	case *packet.ServerboundPluginMessage:
		g.processServerboundPluginMessage(source, p)
	case *packet.ServerboundClientSettings:
		g.processServerboundClientSettings(source, p)
	default:
		g.log.Warn().
			Str("name", pkg.Name()).
			Msg("no processor for packet")
		return
	}
}

func (g *Game) processServerboundClientSettings(source *Player, p *packet.ServerboundClientSettings) {
	source.client.settings.locale = p.Locale
	source.client.settings.viewDistance = p.ViewDistance
}

func (g *Game) processServerboundPluginMessage(source *Player, p *packet.ServerboundPluginMessage) {
	switch p.Channel {
	case id.ParseID("minecraft:brand"):
		rd := bytes.NewReader(p.Data)
		if err := func() (e error) {
			defer func() {
				if rec := recover(); rec != nil {
					if recErr, ok := rec.(error); ok {
						e = recErr
					} else {
						panic(rec)
					}
				}
			}()
			source.client.brand = packet.Decoder{rd}.ReadString("client brand")
			return
		}(); err != nil {
			g.log.Error().
				Err(err).
				Msg("client sent invalid brand, ignoring packet")
		}
	default:
		g.log.Debug().
			Stringer("channel", p.Channel).
			Msg("ignoring plugin message")
	}
}
