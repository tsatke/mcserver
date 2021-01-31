package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseLogin, reflect.TypeOf(ServerboundLoginStart{}))
}

type ServerboundLoginStart struct {
	// Username is the username of the player that is trying to connect.
	Username string
}

func (ServerboundLoginStart) ID() ID       { return IDServerboundLoginStart }
func (ServerboundLoginStart) Name() string { return "Login Start" }

func (s *ServerboundLoginStart) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.Username = dec.ReadString("username")

	return
}

func (s ServerboundLoginStart) Validate() error {
	return multiValidate(
		stringNotEmpty("username", s.Username),
		stringMaxLength("username", 16, s.Username),
	)
}
