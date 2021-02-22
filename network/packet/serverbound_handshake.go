package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseHandshaking, reflect.TypeOf(ServerboundHandshake{}))
}

//go:generate stringer -linecomment -output=serverbound_handshake_string.go -type=NextState

// NextState is a type for the next state constants.
type NextState uint8

// Allowed NextStates.
const (
	NextStateStatus NextState = 1 // Status
	NextStateLogin  NextState = 2 // Login
)

// ServerboundHandshake is the first message that the client sends to the server.
type ServerboundHandshake struct {
	// ProtocolVersion is the client's protocol version.
	ProtocolVersion int
	// ServerAddress is the server address that the client wants to connect to.
	ServerAddress string
	// ServerPort is the server port that the client wants to connect to.
	ServerPort int
	// NextState is the state that this client wants to be put in next.
	NextState NextState
}

// ID returns the constant packet ID.
func (ServerboundHandshake) ID() ID { return IDServerboundHandshake }

// Name returns the constant packet name.
func (ServerboundHandshake) Name() string { return "Handshake" }

// DecodeFrom will fill this struct with values read from the given reader.
func (s *ServerboundHandshake) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.ProtocolVersion = dec.ReadVarInt("protocol version")
	s.ServerAddress = dec.ReadString("server address")
	s.ServerPort = int(dec.ReadUshort("server port"))
	s.NextState = NextState(dec.ReadVarInt("next state"))

	return
}

// Validate implements the Validator interface.
func (s ServerboundHandshake) Validate() error {
	return multiValidate(
		stringNotEmpty("server address", s.ServerAddress),
		stringMaxLength("server address", 255, s.ServerAddress),
		intWithinRange("next state", 1, 2, int(s.NextState)),
	)
}
