package packet

import (
	"bytes"
	"strings"
)

func (suite *PacketSuite) TestServerboundHandshake_DecodeFrom() {
	var buf bytes.Buffer
	enc := Encoder{&buf}
	enc.WriteVarInt("protocol version", 123)
	enc.WriteString("server address", "my.server.com")
	enc.WriteUshort("server port", 17)
	enc.WriteVarInt("next state", 1)

	var p ServerboundHandshake
	suite.NoError(p.DecodeFrom(&buf))
	suite.EqualValues(123, p.ProtocolVersion)
	suite.EqualValues("my.server.com", p.ServerAddress)
	suite.EqualValues(17, p.ServerPort)
	suite.EqualValues(NextStateStatus, p.NextState)
}

func (suite *PacketSuite) TestServerboundHandshake_Validate() {
	type fields struct {
		ProtocolVersion int
		ServerAddress   string
		ServerPort      int
		NextState       NextState
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
				123,
				"abc.de",
				65234,
				1,
			},
			false,
		},
		{
			"invalid server address",
			fields{
				123,
				strings.Repeat("0123456789abcdef", 16), // 256 characters
				65234,
				1,
			},
			true,
		},
		{
			"empty server address",
			fields{
				123,
				"",
				65234,
				1,
			},
			true,
		},
		{
			"invalid state 0",
			fields{
				123,
				"abc.de",
				65234,
				0,
			},
			true,
		},
		{
			"invalid state 3",
			fields{
				123,
				"abc.de",
				65234,
				3,
			},
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			s := ServerboundHandshake{
				ProtocolVersion: tt.fields.ProtocolVersion,
				ServerAddress:   tt.fields.ServerAddress,
				ServerPort:      tt.fields.ServerPort,
				NextState:       tt.fields.NextState,
			}
			if tt.wantErr {
				suite.Error(s.Validate())
			} else {
				suite.NoError(s.Validate())
			}
		})
	}
}
