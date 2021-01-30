package packet

import (
	"io"
	"reflect"

	"github.com/tsatke/mcserver/network/packet/types"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ClientboundDisconnectPlay{}))
}

type ClientboundDisconnectPlay struct {
	Reason types.Chat
}

func (ClientboundDisconnectPlay) ID() ID       { return IDClientboundDisconnectPlay }
func (ClientboundDisconnectPlay) Name() string { return "Disconnect (play)" }

func (c ClientboundDisconnectPlay) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeChat("reason", c.Reason)

	return
}
