package packet

import (
	"io"
	"reflect"

	"github.com/tsatke/mcserver/game/chat"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ClientboundDisconnectPlay{}))
}

type ClientboundDisconnectPlay struct {
	Reason chat.Chat
}

func (ClientboundDisconnectPlay) ID() ID       { return IDClientboundDisconnectPlay }
func (ClientboundDisconnectPlay) Name() string { return "Disconnect (play)" }

func (c ClientboundDisconnectPlay) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteChat("reason", c.Reason)

	return
}
