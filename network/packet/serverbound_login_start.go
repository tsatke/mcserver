package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StateLogin, reflect.TypeOf(ServerboundLoginStart{}))
}

type ServerboundLoginStart struct {
	// Username is the username of the player that is trying to connect.
	Username string
}

func (ServerboundLoginStart) ID() ID       { return IDServerboundLoginStart }
func (ServerboundLoginStart) Name() string { return "Login Start" }

func (s *ServerboundLoginStart) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := decoder{rd}

	s.Username = dec.readString("name")

	return
}
