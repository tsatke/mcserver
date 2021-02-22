package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseLogin, reflect.TypeOf(ServerboundLoginStart{}))
}

// ServerboundLoginStart is sent by the client to set the username that he wants
// to use.
type ServerboundLoginStart struct {
	// Username is the username of the player that is trying to connect.
	Username string
}

// ID returns the constant packet ID.
func (ServerboundLoginStart) ID() ID { return IDServerboundLoginStart }

// Name returns the constant packet name.
func (ServerboundLoginStart) Name() string { return "Login Start" }

// DecodeFrom will fill this struct with values read from the given reader.
func (s *ServerboundLoginStart) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.Username = dec.ReadString("username")

	return
}

// Validate implements the Validator interface.
func (s ServerboundLoginStart) Validate() error {
	return multiValidate(
		stringNotContains("username", s.Username, " "),
		stringNotEmpty("username", s.Username),
		stringMaxLength("username", 16, s.Username),
	)
}
