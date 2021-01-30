package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StateHandshaking, reflect.TypeOf(ServerboundHandshake{}))
}

//go:generate stringer -linecomment -output=serverbound_handshake_string.go -type=NextState

type NextState uint8

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

func (ServerboundHandshake) ID() ID       { return IDServerboundHandshake }
func (ServerboundHandshake) Name() string { return "Handshake" }

func (s *ServerboundHandshake) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := decoder{rd}

	s.ProtocolVersion = dec.readVarInt("protocol version")
	s.ServerAddress = dec.readString("server address")
	s.ServerPort = int(dec.readUshort("server port"))
	s.NextState = NextState(dec.readVarInt("next state"))

	return
}
