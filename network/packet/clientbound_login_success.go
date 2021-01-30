package packet

import (
	"io"
	"reflect"

	"github.com/google/uuid"
)

func init() {
	registerPacket(StateLogin, reflect.TypeOf(ClientboundLoginSuccess{}))
}

type ClientboundLoginSuccess struct {
	UUID     uuid.UUID
	Username string
}

func (ClientboundLoginSuccess) ID() ID       { return IDClientboundLoginSuccess }
func (ClientboundLoginSuccess) Name() string { return "Login Success" }

func (c ClientboundLoginSuccess) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := encoder{w}

	enc.writeUUID("uuid", c.UUID)
	enc.writeString("username", c.Username)

	return
}
