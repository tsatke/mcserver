package comparr

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCompArrSuite(t *testing.T) {
	suite.Run(t, new(CompArrSuite))
}

type CompArrSuite struct {
	suite.Suite
}

func (suite *CompArrSuite) TestNew() {
	type args struct {
		length  int
		bitSize int
	}
	tests := []struct {
		name string
		args args
		want CompArr
	}{
		{
			"empty",
			args{0, 1},
			CompArr{
				length:      0,
				bitsz:       1,
				mask:        1,
				elemsPerInt: 64,
				data:        []uint64{},
			},
		},
		{
			"1/5",
			args{1, 5},
			CompArr{
				length:      1,
				bitsz:       5,
				mask:        0b11111,
				elemsPerInt: 12,
				data: []uint64{
					0,
				},
			},
		},
		{
			"12/5",
			args{12, 5},
			CompArr{
				length:      12,
				bitsz:       5,
				mask:        0b11111,
				elemsPerInt: 12,
				data: []uint64{
					0,
				},
			},
		},
		{
			"13/5",
			args{13, 5},
			CompArr{
				length:      13,
				bitsz:       5,
				mask:        0b11111,
				elemsPerInt: 12,
				data: []uint64{
					0,
					0,
				},
			},
		},
		{
			"1/64",
			args{1, 64},
			CompArr{
				length:      1,
				bitsz:       64,
				mask:        ^uint64(0),
				elemsPerInt: 1,
				data: []uint64{
					0,
				},
			},
		},
		{
			"64/8",
			args{64, 8},
			CompArr{
				length:      64,
				bitsz:       8,
				mask:        0b11111111,
				elemsPerInt: 8,
				data: []uint64{
					0,
					0,
					0,
					0,
					0,
					0,
					0,
					0,
				},
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := New(tt.args.length, tt.args.bitSize)
			suite.Equal(tt.want, got)
		})
	}
}

func (suite *CompArrSuite) TestSet() {
	suite.Run("datalen=1", func() {
		arr := New(8, 8)
		suite.Equal(arr.length, 8)
		suite.Len(arr.data, 1)

		arr.Set(0, 1)
		arr.Set(1, 2)
		arr.Set(2, 3)
		arr.Set(3, 4)
		arr.Set(4, 5)
		arr.Set(5, 6)
		arr.Set(6, 7)
		arr.Set(7, 8)

		suite.EqualValues(1, arr.data[0]&((1<<8)-1))
		suite.EqualValues(2, (arr.data[0]&(((1<<8)-1)<<8))>>8)
		suite.EqualValues(3, (arr.data[0]&(((1<<8)-1)<<16))>>16)
		suite.EqualValues(4, (arr.data[0]&(((1<<8)-1)<<24))>>24)
		suite.EqualValues(5, (arr.data[0]&(((1<<8)-1)<<32))>>32)
		suite.EqualValues(6, (arr.data[0]&(((1<<8)-1)<<40))>>40)
		suite.EqualValues(7, (arr.data[0]&(((1<<8)-1)<<48))>>48)
		suite.EqualValues(8, (arr.data[0]&(((1<<8)-1)<<56))>>56)
	})
	suite.Run("datalen=2", func() {
		arr := New(16, 8)
		suite.Equal(arr.length, 16)
		suite.Len(arr.data, 2)

		arr.Set(0, 1)
		arr.Set(1, 2)
		arr.Set(2, 3)
		arr.Set(3, 4)
		arr.Set(4, 5)
		arr.Set(5, 6)
		arr.Set(6, 7)
		arr.Set(7, 8)

		arr.Set(8, 11)
		arr.Set(9, 12)
		arr.Set(10, 13)
		arr.Set(11, 14)
		arr.Set(12, 15)
		arr.Set(13, 16)
		arr.Set(14, 17)
		arr.Set(15, 18)

		suite.EqualValues(1, arr.data[1]&((1<<8)-1))
		suite.EqualValues(2, (arr.data[1]&(((1<<8)-1)<<8))>>8)
		suite.EqualValues(3, (arr.data[1]&(((1<<8)-1)<<16))>>16)
		suite.EqualValues(4, (arr.data[1]&(((1<<8)-1)<<24))>>24)
		suite.EqualValues(5, (arr.data[1]&(((1<<8)-1)<<32))>>32)
		suite.EqualValues(6, (arr.data[1]&(((1<<8)-1)<<40))>>40)
		suite.EqualValues(7, (arr.data[1]&(((1<<8)-1)<<48))>>48)
		suite.EqualValues(8, (arr.data[1]&(((1<<8)-1)<<56))>>56)
		suite.EqualValues(11, arr.data[0]&((1<<8)-1))
		suite.EqualValues(12, (arr.data[0]&(((1<<8)-1)<<8))>>8)
		suite.EqualValues(13, (arr.data[0]&(((1<<8)-1)<<16))>>16)
		suite.EqualValues(14, (arr.data[0]&(((1<<8)-1)<<24))>>24)
		suite.EqualValues(15, (arr.data[0]&(((1<<8)-1)<<32))>>32)
		suite.EqualValues(16, (arr.data[0]&(((1<<8)-1)<<40))>>40)
		suite.EqualValues(17, (arr.data[0]&(((1<<8)-1)<<48))>>48)
		suite.EqualValues(18, (arr.data[0]&(((1<<8)-1)<<56))>>56)
	})
}

func (suite *CompArrSuite) TestGet() {
	arr := New(16, 8)
	suite.Equal(arr.length, 16)
	suite.Len(arr.data, 2)

	arr.Set(0, 1)
	arr.Set(1, 2)
	arr.Set(2, 3)
	arr.Set(3, 4)
	arr.Set(4, 5)
	arr.Set(5, 6)
	arr.Set(6, 7)
	arr.Set(7, 8)

	arr.Set(8, 11)
	arr.Set(9, 12)
	arr.Set(10, 13)
	arr.Set(11, 14)
	arr.Set(12, 15)
	arr.Set(13, 16)
	arr.Set(14, 17)
	arr.Set(15, 18)

	suite.EqualValues(1, arr.Get(0))
	suite.EqualValues(2, arr.Get(1))
	suite.EqualValues(3, arr.Get(2))
	suite.EqualValues(4, arr.Get(3))
	suite.EqualValues(5, arr.Get(4))
	suite.EqualValues(6, arr.Get(5))
	suite.EqualValues(7, arr.Get(6))
	suite.EqualValues(8, arr.Get(7))
	suite.EqualValues(11, arr.Get(8))
	suite.EqualValues(12, arr.Get(9))
	suite.EqualValues(13, arr.Get(10))
	suite.EqualValues(14, arr.Get(11))
	suite.EqualValues(15, arr.Get(12))
	suite.EqualValues(16, arr.Get(13))
	suite.EqualValues(17, arr.Get(14))
	suite.EqualValues(18, arr.Get(15))
}
