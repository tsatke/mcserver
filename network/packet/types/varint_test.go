package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeVarInt(t *testing.T) {
	tests := []struct {
		name    string
		arg     []byte
		want    VarInt
		wantErr bool
	}{
		{
			"zero",
			[]byte{0x00},
			0,
			false,
		},
		{
			"one",
			[]byte{0x01},
			1,
			false,
		},
		{
			"two",
			[]byte{0x02},
			2,
			false,
		},
		{
			"127",
			[]byte{0x7f},
			127,
			false,
		},
		{
			"128",
			[]byte{0x80, 0x01},
			128,
			false,
		},
		{
			"255",
			[]byte{0xff, 0x01},
			255,
			false,
		},
		{
			"754",
			[]byte{0xf2, 0x05},
			754,
			false,
		},
		{
			"2097151",
			[]byte{0xff, 0xff, 0x7f},
			2097151,
			false,
		},
		{
			"2147483647",
			[]byte{0xff, 0xff, 0xff, 0xff, 0x07},
			2147483647,
			false,
		},
		{
			"-1",
			[]byte{0xff, 0xff, 0xff, 0xff, 0x0f},
			-1,
			false,
		},
		{
			"-2147483648",
			[]byte{0x80, 0x80, 0x80, 0x80, 0x08},
			-2147483648,
			false,
		},
	}
	for _, tt := range tests {
		t.Run("mode=encode/"+tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.want.EncodeInto(&buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.arg, buf.Bytes())
		})
		t.Run("mode=decode/"+tt.name, func(t *testing.T) {
			got := NewVarInt(0)
			err := got.DecodeFrom(bytes.NewReader(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, *got)
		})
	}
}
