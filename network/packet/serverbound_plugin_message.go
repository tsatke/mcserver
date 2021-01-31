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

type ServerboundPluginMessage struct {
	Channel id.ID
	Data    []byte
}

func (ServerboundPluginMessage) ID() ID       { return IDServerboundPluginMessage }
func (ServerboundPluginMessage) Name() string { return "Plugin Message" }

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
