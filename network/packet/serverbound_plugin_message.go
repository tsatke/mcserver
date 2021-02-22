package packet

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/tsatke/mcserver/game/id"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ServerboundPluginMessage{}))
}

// ServerboundPluginMessage is used by the client to send out-of-protocol
// data, such as the client brand.
type ServerboundPluginMessage struct {
	Channel id.ID
	Data    []byte
}

// ID returns the constant packet ID.
func (ServerboundPluginMessage) ID() ID { return IDServerboundPluginMessage }

// Name returns the constant packet name.
func (ServerboundPluginMessage) Name() string { return "Plugin Message" }

// DecodeFrom will fill this struct with values read from the given reader.
func (s *ServerboundPluginMessage) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := Decoder{rd}

	s.Channel = dec.ReadID("channel")
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}
	s.Data = data

	return
}
