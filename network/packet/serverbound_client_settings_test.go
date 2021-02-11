package packet

import (
	"bytes"
)

func (suite *PacketSuite) TestServerboundClientSettings_DecodeFrom() {
	var buf bytes.Buffer
	enc := Encoder{&buf}
	enc.WriteString("locale", "en_US")
	enc.WriteByte("view distance", 32)
	enc.WriteVarInt("chat mode", 0)
	enc.WriteBoolean("chat colors", true)
	enc.WriteUbyte("displayed skin parts", ^byte(0))
	enc.WriteVarInt("main hand", 0)

	var p ServerboundClientSettings
	suite.NoError(p.DecodeFrom(&buf))
	suite.Equal("en_US", p.Locale)
	suite.EqualValues(32, p.ViewDistance)
	suite.Equal(ChatModeEnabled, p.ChatMode)
	suite.Equal(true, p.ChatColors)
	suite.EqualValues(^byte(0), p.DisplayedSkinParts)
	suite.Equal(HandLeft, p.MainHand)
}

func (suite *PacketSuite) TestServerboundClientSettings_Validate() {
	type fields struct {
		Locale             string
		ViewDistance       int
		ChatMode           ChatMode
		ChatColors         bool
		DisplayedSkinParts byte
		MainHand           Hand
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"empty",
			fields{},
			true,
		},
		{
			"valid",
			fields{
				"en_US",
				32,
				ChatModeEnabled,
				true,
				0,
				HandLeft,
			},
			false,
		},
		{
			"empty locale",
			fields{
				"",
				32,
				ChatModeEnabled,
				true,
				0,
				HandLeft,
			},
			true,
		},
		{
			"invalid hand",
			fields{
				"en_US",
				32,
				ChatModeEnabled,
				true,
				0,
				7,
			},
			true,
		},
		{
			"view distance too big",
			fields{
				"en_US",
				999,
				ChatModeEnabled,
				true,
				0,
				HandLeft,
			},
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			s := ServerboundClientSettings{
				Locale:             tt.fields.Locale,
				ViewDistance:       tt.fields.ViewDistance,
				ChatMode:           tt.fields.ChatMode,
				ChatColors:         tt.fields.ChatColors,
				DisplayedSkinParts: tt.fields.DisplayedSkinParts,
				MainHand:           tt.fields.MainHand,
			}
			if tt.wantErr {
				suite.Error(s.Validate())
			} else {
				suite.NoError(s.Validate())
			}
		})
	}
}
