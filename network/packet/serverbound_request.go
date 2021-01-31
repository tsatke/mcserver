package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(StateStatus, reflect.TypeOf(ServerboundRequest{}))
}

type ServerboundRequest struct {
}

func (ServerboundRequest) ID() ID       { return IDServerboundRequest }
func (ServerboundRequest) Name() string { return "Request" }

func (s *ServerboundRequest) DecodeFrom(rd io.Reader) (err error) {
	return
}
