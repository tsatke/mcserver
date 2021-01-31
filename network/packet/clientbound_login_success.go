package packet

import (
	"io"
	"reflect"

	"github.com/google/uuid"
)

func init() {
	RegisterPacket(PhaseLogin, reflect.TypeOf(ClientboundLoginSuccess{}))
}

type ClientboundLoginSuccess struct {
	UUID     uuid.UUID
	Username string
}

func (ClientboundLoginSuccess) ID() ID       { return IDClientboundLoginSuccess }
func (ClientboundLoginSuccess) Name() string { return "Login Success" }

func (c ClientboundLoginSuccess) EncodeInto(w io.Writer) (err error) {
	defer recoverAndSetErr(&err)

	enc := Encoder{w}

	enc.WriteUUID("uuid", c.UUID)
	enc.WriteString("username", c.Username)

	return
}
