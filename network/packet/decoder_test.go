package packet

import (
	"bytes"
	"io"
	"runtime"
	"testing"
	"testing/iotest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func TestDecoderSuite(t *testing.T) {
	suite.Run(t, new(DecoderSuite))
}

type DecoderSuite struct {
	suite.Suite
}

func (suite *DecoderSuite) TestDecoder_ReadVarInt() {
	tests := []struct {
		name    string
		source  io.Reader
		want    int
		wantErr bool
	}{
		{
			"zero",
			bytes.NewReader([]byte{0x00}),
			0,
			false,
		},
		{
			"one",
			bytes.NewReader([]byte{0x01}),
			1,
			false,
		},
		{
			"two",
			bytes.NewReader([]byte{0x02}),
			2,
			false,
		},
		{
			"127",
			bytes.NewReader([]byte{0x7f}),
			127,
			false,
		},
		{
			"128",
			bytes.NewReader([]byte{0x80, 0x01}),
			128,
			false,
		},
		{
			"128 onebyte",
			iotest.OneByteReader(bytes.NewReader([]byte{0x80, 0x01})),
			128,
			false,
		},
		{
			"128 half",
			iotest.HalfReader(bytes.NewReader([]byte{0x80, 0x01})),
			128,
			false,
		},
		{
			"128 data err",
			iotest.DataErrReader(bytes.NewReader([]byte{0x80, 0x01})),
			128,
			false,
		},
		{
			"255",
			bytes.NewReader([]byte{0xff, 0x01}),
			255,
			false,
		},
		{
			"754",
			bytes.NewReader([]byte{0xf2, 0x05}),
			754,
			false,
		},
		{
			"2097151",
			bytes.NewReader([]byte{0xff, 0xff, 0x7f}),
			2097151,
			false,
		},
		{
			"2147483647",
			bytes.NewReader([]byte{0xff, 0xff, 0xff, 0xff, 0x07}),
			2147483647,
			false,
		},
		{
			"-1",
			bytes.NewReader([]byte{0xff, 0xff, 0xff, 0xff, 0x0f}),
			-1,
			false,
		},
		{
			"-2147483648",
			bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80, 0x08}),
			-2147483648,
			false,
		},
		{
			"-2147483648 onebyte",
			iotest.OneByteReader(bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80, 0x08})),
			-2147483648,
			false,
		},
		{
			"-2147483648 half",
			iotest.HalfReader(bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80, 0x08})),
			-2147483648,
			false,
		},
		{
			"-2147483648 data err",
			iotest.DataErrReader(bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80, 0x08})),
			-2147483648,
			false,
		},
		{
			"-2147483648 timeout",
			iotest.TimeoutReader(bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80, 0x08})),
			-2147483648,
			true,
		},
		{
			"too big",
			bytes.NewReader([]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x08}),
			0,
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				dec := Decoder{tt.source}
				got := dec.ReadVarInt("field")
				suite.EqualValues(tt.want, got)
			}
			if tt.wantErr {
				suite.Panics(testFn) // the decoder fails by panicking an error
			} else {
				suite.NotPanics(testFn)
			}
		})
	}
}

// Test_readStringFaultyLengthPrefix tests that ReadString doesn't prematurely
// and blindly allocates memory by listening on the prefix. If this test timeouts,
// crashes with not enough memory or just hangs, chances are that the implementation
// is broken and does exactly that.
func (suite *DecoderSuite) TestDecoder_ReadStringFaultyLengthPrefix() {
	source := bytes.NewReader([]byte{
		0xff, 0xff, 0xff, 0xff, 0x07, // 2GiB length prefix
		0x40, 0x40, 0x40, 0x40, 0x40, // but only 5 bytes
	})
	dec := Decoder{source}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	oldHeap := m.HeapInuse

	suite.Panics(func() {
		_ = dec.ReadString("field")
	})

	runtime.ReadMemStats(&m)
	newHeap := m.HeapInuse
	suite.Less(newHeap-oldHeap, uint64(1<<30)) // ReadString must have allocated less than 1GiB on the heap during reading
}

func (suite *DecoderSuite) TestDecoder_ReadString() {
	tests := []struct {
		name    string
		source  io.Reader
		want    string
		wantErr bool
	}{
		{
			"empty input",
			bytes.NewReader([]byte{}),
			"",
			true,
		},
		{
			"empty",
			bytes.NewReader([]byte{0x00}),
			"",
			false,
		},
		{
			"small",
			bytes.NewReader([]byte{0x02, 0x63, 0x64}),
			"cd",
			false,
		},
		{
			"with remaining",
			bytes.NewReader([]byte{0x02, 0x63, 0x64, 0x65, 0x66}),
			"cd",
			false,
		},
		{
			"small onebyte",
			iotest.OneByteReader(bytes.NewReader([]byte{0x02, 0x63, 0x64})),
			"cd",
			false,
		},
		{
			"small data err",
			iotest.DataErrReader(bytes.NewReader([]byte{0x02, 0x63, 0x64})),
			"cd",
			false,
		},
		{
			"small timeout",
			iotest.OneByteReader(iotest.TimeoutReader(bytes.NewReader([]byte{0x04, 0x63, 0x64, 0x65, 0x66}))),// make sure there are more than two reads for the TimeoutReader

			"cd",
			true,
		},
		{
			"fewer payload than length",
			bytes.NewReader([]byte{0x4f, 0x39, 0x40}),
			"",
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				dec := Decoder{tt.source}
				got := dec.ReadString("field")
				suite.EqualValues(tt.want, got)
			}
			if tt.wantErr {
				suite.Panics(testFn) // the decoder fails by panicking an error
			} else {
				suite.NotPanics(testFn)
			}
		})
	}
}

func (suite *DecoderSuite) TestDecoder_ReadUbyte() {
	tests := []struct {
		name    string
		source  io.Reader
		want    byte
		wantErr bool
	}{
		{
			"zero",
			bytes.NewReader([]byte{0x00}),
			0,
			false,
		},
		{
			"one",
			bytes.NewReader([]byte{0x01}),
			1,
			false,
		},
		{
			"hi msb",
			bytes.NewReader([]byte{0xf0}),
			240,
			false,
		},
		{
			"more remaining",
			bytes.NewReader([]byte{0x00, 0x01, 0x02}),
			0,
			false,
		},
		{
			"dataerr",
			iotest.DataErrReader(bytes.NewReader([]byte{0x04})),
			4,
			false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				dec := Decoder{tt.source}
				got := dec.ReadUbyte("field")
				suite.EqualValues(tt.want, got)
			}
			if tt.wantErr {
				suite.Panics(testFn) // the decoder fails by panicking an error
			} else {
				suite.NotPanics(testFn)
			}
		})
	}
}

func (suite *DecoderSuite) TestDecoder_ReadByte() {
	tests := []struct {
		name    string
		source  io.Reader
		want    int8
		wantErr bool
	}{
		{
			"zero",
			bytes.NewReader([]byte{0x00}),
			0,
			false,
		},
		{
			"one",
			bytes.NewReader([]byte{0x01}),
			1,
			false,
		},
		{
			"hi msb",
			bytes.NewReader([]byte{0xf0}),
			-16,
			false,
		},
		{
			"more remaining",
			bytes.NewReader([]byte{0x00, 0x01, 0x02}),
			0,
			false,
		},
		{
			"dataerr",
			iotest.DataErrReader(bytes.NewReader([]byte{0x04})),
			4,
			false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				dec := Decoder{tt.source}
				got := dec.ReadByte("field")
				suite.EqualValues(tt.want, got)
			}
			if tt.wantErr {
				suite.Panics(testFn) // the decoder fails by panicking an error
			} else {
				suite.NotPanics(testFn)
			}
		})
	}
}

func (suite *DecoderSuite) TestDecoder_ReadBoolean() {
	tests := []struct {
		name    string
		source  io.Reader
		want    bool
		wantErr bool
	}{
		{
			"false",
			bytes.NewReader([]byte{0x00}),
			false,
			false,
		},
		{
			"true",
			bytes.NewReader([]byte{0x01}),
			true,
			false,
		},
		{
			"two",
			bytes.NewReader([]byte{0x02}),
			false,
			false,
		},
		{
			"hi msb",
			bytes.NewReader([]byte{0xf0}),
			false,
			false,
		},
		{
			"more remaining",
			bytes.NewReader([]byte{0x00, 0x01, 0x02}),
			false,
			false,
		},
		{
			"dataerr",
			iotest.DataErrReader(bytes.NewReader([]byte{0x01})),
			true,
			false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				dec := Decoder{tt.source}
				got := dec.ReadBoolean("field")
				suite.EqualValues(tt.want, got)
			}
			if tt.wantErr {
				suite.Panics(testFn) // the decoder fails by panicking an error
			} else {
				suite.NotPanics(testFn)
			}
		})
	}
}

func (suite *DecoderSuite) TestDecoder_ReadUshort() {
	tests := []struct {
		name    string
		source  io.Reader
		want    uint16
		wantErr bool
	}{
		{
			"zero",
			bytes.NewReader([]byte{0x00, 0x00}),
			0,
			false,
		},
		{
			"one",
			bytes.NewReader([]byte{0x00, 0x01}),
			1,
			false,
		},
		{
			"hi msb",
			bytes.NewReader([]byte{0xf0, 0x00}),
			240 << 8,
			false,
		},
		{
			"more remaining",
			bytes.NewReader([]byte{0x00, 0x01, 0x02}),
			1,
			false,
		},
		{
			"dataerr",
			iotest.DataErrReader(bytes.NewReader([]byte{0x00, 0x04})),
			4,
			false,
		},
		{
			"onebyte",
			iotest.OneByteReader(bytes.NewReader([]byte{0x00, 0x04})),
			4,
			false,
		},
		{
			"too few bytes",
			bytes.NewReader([]byte{0x07}),
			0,
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				dec := Decoder{tt.source}
				got := dec.ReadUshort("field")
				suite.EqualValues(tt.want, got)
			}
			if tt.wantErr {
				suite.Panics(testFn) // the decoder fails by panicking an error
			} else {
				suite.NotPanics(testFn)
			}
		})
	}
}

func (suite *DecoderSuite) TestDecoder_ReadUUID() {
	uid, err := uuid.Parse("6bb0b268-63c8-11eb-bed6-acde48001122")
	suite.NoError(err)

	tests := []struct {
		name    string
		source  io.Reader
		want    string
		wantErr bool
	}{
		{
			"uuid",
			bytes.NewReader(uid[:]),
			"6bb0b268-63c8-11eb-bed6-acde48001122",
			false,
		},
		{
			"remaining bytes",
			bytes.NewReader(append(uid[:], 0x00, 0x12, 0x42)),
			"6bb0b268-63c8-11eb-bed6-acde48001122",
			false,
		},
		{
			"onebyte",
			iotest.OneByteReader(bytes.NewReader(uid[:])),
			"6bb0b268-63c8-11eb-bed6-acde48001122",
			false,
		},
		{
			"dataerr",
			iotest.DataErrReader(bytes.NewReader(uid[:])),
			"6bb0b268-63c8-11eb-bed6-acde48001122",
			false,
		},
		{
			"half",
			iotest.HalfReader(bytes.NewReader(uid[:])),
			"6bb0b268-63c8-11eb-bed6-acde48001122",
			false,
		},
		{
			"timeout",
			iotest.OneByteReader(iotest.TimeoutReader(bytes.NewReader(uid[:]))),
			"",
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			testFn := func() {
				dec := Decoder{tt.source}
				got := dec.ReadUUID("field")

				// if the test panics, as expected when tt.wantErr==true, then
				// this will not be reached
				want, err := uuid.Parse(tt.want)
				suite.NoError(err)
				suite.Equal(want, got)
			}
			if tt.wantErr {
				suite.Panics(testFn) // the decoder fails by panicking an error
			} else {
				suite.NotPanics(testFn)
			}
		})
	}
}
