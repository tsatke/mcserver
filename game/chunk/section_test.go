package chunk

import (
	"reflect"
	"testing"
)

func Test_splitArrayIntoBitSegments(t *testing.T) {
	type args struct {
		array         []int64
		segmentLength int
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{
			"simple",
			args{
				[]int64{
					1,
				},
				4,
			},
			[]uint64{
				1, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
		{
			"double",
			args{
				[]int64{
					1,
					2,
				},
				4,
			},
			[]uint64{
				1, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,

				2, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitArrayIntoBitSegments(tt.args.array, tt.args.segmentLength); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitArrayIntoBitSegments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createHiMask(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			"zero",
			args{0},
			0,
		},
		{
			"one",
			args{1},
			0b1,
		},
		{
			"two",
			args{2},
			0b11,
		},
		{
			"12",
			args{12},
			0b111111111111,
		},
		{
			"64",
			args{64},
			0xffffffffffffffff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createHiMask(tt.args.size); got != tt.want {
				t.Errorf("createHiMask() = %v, want %v", got, tt.want)
			}
		})
	}
}
