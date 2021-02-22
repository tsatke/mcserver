package packet

import (
	"io"
	"reflect"
)

func init() {
	RegisterPacket(PhasePlay, reflect.TypeOf(ServerboundClientSettings{}))
}

// ChatMode is a type for allowed constants.
type ChatMode int

// Allowed constants.
const (
	ChatModeEnabled ChatMode = iota
	ChatModeCommandsOnly
	ChatModeHidden
)

// Hand is a type for allowed constants.
type Hand int

// Allowed constants.
const (
	HandLeft Hand = iota
	HandRight
)

// ServerboundClientSettings are sent by the client if his settings change and when
// he connects.
type ServerboundClientSettings struct {
	// Locale is the locale that is used by the client. E.g. en_US or de_CH.
	Locale string
	// ViewDistance is the view distance that the client uses. Chunks
	// should be sent accordingly.
	ViewDistance int
	// ChatMode is the chat settings that the client has enabled.
	// Chat messages should be sent accordingly.
	ChatMode ChatMode
	// ChatColors indicates whether the client has chat colors enabled.
	ChatColors bool
	// DisplayedSkinParts is a bit mask that shows which skin parts are
	// shown in the client.
	//
	//	0x01: Cape
	//	0x02: Jacket
	//	0x04: Left Sleeve
	//	0x08: Right Sleeve
	//	0x10: Left Pants Leg
	//	0x20: Right Pants Leg
	//	0x40: Hat
	//
	// If a bit is set, that means that the respective part is shown.
	DisplayedSkinParts byte
	// MainHand is the main hand that the client shows, HandLeft or HandRight.
	MainHand Hand
}

// ID returns the constant packet ID.
func (ServerboundClientSettings) ID() ID { return IDServerboundClientSettings }

// Name returns the constant packet name.
func (ServerboundClientSettings) Name() string { return "Client Settings" }

// DecodeFrom will fill this struct with values read from the given reader.
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

// Validate implements the Validator interface.
func (s *ServerboundClientSettings) Validate() error {
	return multiValidate(
		stringNotEmpty("locale", s.Locale),
		intWithinRange("view distance", 1, 48, s.ViewDistance),
		intWithinRange("chat mode", 0, 2, int(s.ChatMode)),
		intWithinRange("main hand", 0, 1, int(s.MainHand)),
	)
}
