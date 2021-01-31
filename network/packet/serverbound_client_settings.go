package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ServerboundClientSettings{}))
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

	dec := Decoder{rd}

	s.Locale = dec.ReadString("locale")
	s.ViewDistance = int(dec.ReadByte("view distance"))
	s.ChatMode = ChatMode(dec.ReadVarInt("chat mode"))
	s.ChatColors = dec.ReadBoolean("chat colors")
	s.DisplayedSkinParts = dec.ReadUbyte("displayed skin parts")
	s.MainHand = Hand(dec.ReadVarInt("main hand"))

	return
}
