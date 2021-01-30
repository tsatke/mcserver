package packet

import (
	"io"
	"reflect"
)

func init() {
	registerPacket(StatePlay, reflect.TypeOf(ServerboundClientSettings{}))
}

type ChatMode int

const (
	ChatModeEnabled ChatMode = iota
	ChatModeCommandsOnly
	ChatModeHidden
)

type Hand int

const (
	HandLeft Hand = iota
	HandRight
)

type ServerboundClientSettings struct {
	Locale             string
	ViewDistance       int
	ChatMode           ChatMode
	ChatColors         bool
	DisplayedSkinParts byte
	MainHand           Hand
}

func (ServerboundClientSettings) ID() ID       { return IDServerboundClientSettings }
func (ServerboundClientSettings) Name() string { return "Client Settings" }

func (s *ServerboundClientSettings) DecodeFrom(rd io.Reader) (err error) {
	defer recoverAndSetErr(&err)

	dec := decoder{rd}

	s.Locale = dec.readString("locale")
	s.ViewDistance = int(dec.readByte("view distance"))
	s.ChatMode = ChatMode(dec.readVarInt("chat mode"))
	s.ChatColors = dec.readBoolean("chat colors")
	s.DisplayedSkinParts = dec.readUbyte("displayed skin parts")
	s.MainHand = Hand(dec.readVarInt("main hand"))

	return
}
