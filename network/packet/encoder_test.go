package packet

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/tsatke/mcserver/game/chat"
)

func TestEncoderSuite(t *testing.T) {
	suite.Run(t, new(EncoderSuite))
}

type EncoderSuite struct {
	suite.Suite
}

func (suite *EncoderSuite) TestEncoder_WriteBoolean() {
	tests := []struct {
		name string
		val  bool
		want []byte
	}{
		{
			"true",
			true,
			[]byte{0x01},
		},
		{
			"false",
			false,
			[]byte{0x00},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteBoolean("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteByte() {
	for i := -130; i < 130; i++ { // few over and underflows
		suite.Run("byte="+strconv.Itoa(i), func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteByte("field", int8(i))
				suite.EqualValues([]byte{byte(i)}, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteByteArray() {
	tests := []struct {
		name string
		val  []byte
	}{
		{
			"empty",
			[]byte{},
		},
		{
			"single",
			[]byte{0x53},
		},
		{
			"few",
			[]byte{0, 01, 0x02, 0x03, 0x04, 0x05},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteByteArray("field", tt.val)
				if len(tt.val) == 0 {
					suite.Len(buf.Bytes(), 0)
				} else {
					suite.EqualValues(tt.val, buf.Bytes())
				}
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteChat() {
	tests := []struct {
		name string
		val  chat.Chat
		want []byte
	}{
		{
			"simple",
			chat.Chat{
				ChatFragment: chat.ChatFragment{
					Text: "text",
				},
			},
			append(
				[]byte{0x0f}, // string length prefix
				[]byte(`{"text":"text"}`)...,
			),
		},
		{
			"hello world",
			chat.Chat{
				ChatFragment: chat.ChatFragment{
					Text: "Hello, World!",
				},
			},
			append(
				[]byte{0x18}, // string length prefix
				[]byte(`{"text":"Hello, World!"}`)...,
			),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteChat("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteDouble() {
	tests := []struct {
		name string
		val  float64
	}{
		{
			"zero",
			0,
		},
		{
			"NaN",
			math.NaN(),
		},
		{
			"PosInf",
			math.Inf(1),
		},
		{
			"NegInf",
			math.Inf(-1),
		},
		{
			".1",
			.1,
		},
		{
			"big",
			876584235903485761239042873956412908253467259187039825362395293.987216546812352519247612,
		},
		{
			"big negative",
			-876584235903485761239042873956412908253467259187039825362395293.987216546812352519247612,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				var wantBuf [8]byte
				ByteOrder.PutUint64(wantBuf[:], math.Float64bits(tt.val))
				enc.WriteDouble("field", tt.val)
				suite.EqualValues(wantBuf[:], buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteFloat() {
	tests := []struct {
		name string
		val  float32
	}{
		{
			"zero",
			0,
		},
		{
			"NaN",
			float32(math.NaN()),
		},
		{
			"PosInf",
			float32(math.Inf(1)),
		},
		{
			"NegInf",
			float32(math.Inf(-1)),
		},
		{
			".1",
			.1,
		},
		{
			"big",
			87658423590348576123907039825362395293.9872165462352519247612,
		},
		{
			"big negative",
			-87658423590348576123907039825362395293.9872165462352519247612,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				var wantBuf [4]byte
				ByteOrder.PutUint32(wantBuf[:], math.Float32bits(tt.val))
				enc.WriteFloat("field", tt.val)
				suite.EqualValues(wantBuf[:], buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteInt() {
	suite.Require().Equal(binary.BigEndian, ByteOrder, "test only valid for big endian")

	tests := []struct {
		name string
		val  int32
		want []byte
	}{
		{
			"zero",
			0,
			[]byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			"one",
			1,
			[]byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			"200",
			200,
			[]byte{0x00, 0x00, 0x00, 0xc8},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteInt("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteLong() {
	suite.Require().Equal(binary.BigEndian, ByteOrder, "test only valid for big endian")

	tests := []struct {
		name string
		val  int64
		want []byte
	}{
		{
			"zero",
			0,
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			"one",
			1,
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			"200",
			200,
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc8},
		},
		{
			"big",
			276834129345232,
			[]byte{0x00, 0x00, 0xfb, 0xc7, 0x77, 0xf0, 0xfa, 0xd0},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteLong("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteString() {
	tests := []struct {
		name string
		val  string
		want []byte
	}{
		{
			"empty",
			"",
			[]byte{0x00},
		},
		{
			"single ascii",
			"a",
			[]byte{1, 'a'},
		},
		{
			"triple ascii",
			"abc",
			[]byte{3, 'a', 'b', 'c'},
		},
		{
			"multibyte",
			"\u1023",
			[]byte{3, 0xe1, 0x80, 0xa3},
		},
		{
			"long",
			string(bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}, 100)),
			append(
				[]byte{0xe8, 0x07}, // string prefix length as varint
				bytes.Repeat([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}, 100)...,
			),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteString("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteUUID() {
	tests := []struct {
		name string
		val  uuid.UUID
	}{
		{
			"uuid",
			uuid.New(),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteUUID("field", tt.val)
				suite.EqualValues(tt.val[:], buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteUbyte() {
	tests := []struct {
		name string
		val  uint8
		want []byte
	}{
		{
			"zero",
			0,
			[]byte{0x00},
		},
		{
			"12",
			12,
			[]byte{0x0c},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteUbyte("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteUshort() {
	suite.Require().Equal(binary.BigEndian, ByteOrder, "test only valid for big endian")

	tests := []struct {
		name string
		val  uint16
		want []byte
	}{
		{
			"zero",
			0,
			[]byte{0x00, 0x00},
		},
		{
			"100",
			100,
			[]byte{0x00, 0x64},
		},
		{
			"1000",
			1000,
			[]byte{0x03, 0xe8},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteUshort("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}

func (suite *EncoderSuite) TestEncoder_WriteVarInt() {
	tests := []struct {
		name string
		val  int
		want []byte
	}{
		{
			"zero",
			0,
			[]byte{0x00},
		},
		{
			"one",
			1,
			[]byte{0x01},
		},
		{
			"two",
			2,
			[]byte{0x02},
		},
		{
			"127",
			127,
			[]byte{0x7f},
		},
		{
			"128",
			128,
			[]byte{0x80, 0x01},
		},
		{
			"255",
			255,
			[]byte{0xff, 0x01},
		},
		{
			"754",
			754,
			[]byte{0xf2, 0x05},
		},
		{
			"2097151",
			2097151,
			[]byte{0xff, 0xff, 0x7f},
		},
		{
			"2147483647",
			2147483647,
			[]byte{0xff, 0xff, 0xff, 0xff, 0x07},
		},
		{
			"-1",
			-1,
			[]byte{0xff, 0xff, 0xff, 0xff, 0x0f},
		},
		{
			"-2147483648",
			-2147483648,
			[]byte{0x80, 0x80, 0x80, 0x80, 0x08},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				var buf bytes.Buffer
				enc := Encoder{&buf}

				enc.WriteVarInt("field", tt.val)
				suite.EqualValues(tt.want, buf.Bytes())
			}
			suite.NotPanics(testFn)
		})
	}
}
