package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhaseStatus, reflect.TypeOf(ServerboundRequest{}))
}

// ServerboundRequest is sent by the client during the Status phase.
// After this packet, the client expects the server to send a
// ClientboundResponse.
type ServerboundRequest struct {
}

// ID returns the constant packet ID.
func (ServerboundRequest) ID() ID { return IDServerboundRequest }

// Name returns the constant packet name.
func (ServerboundRequest) Name() string { return "Request" }

// DecodeFrom will fill this struct with values read from the given reader.
func (s *ServerboundRequest) DecodeFrom(rd io.Reader) (err error) {
	return
}
