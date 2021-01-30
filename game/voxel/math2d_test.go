package voxel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CircleAround(t *testing.T) {
	type args struct {
		center V2
		radius float64
	}
	tests := []struct {
		name    string
		args    args
		wantRes []V2
	}{
		{
			"zero radius",
			args{
				V2{0, 0},
				0,
			},
			[]V2{
				{0, 0},
			},
		},
		{
			"one radius",
			args{
				V2{0, 0},
				1,
			},
			[]V2{
				{0, -1},
				{-1, 0}, {0, 0}, {1, 0},
				{0, 1},
			},
		},
		{
			"1.4 radius",
			args{
				V2{0, 0},
				1.4,
			},
			[]V2{
				{0, -1},
				{-1, 0}, {0, 0}, {1, 0},
				{0, 1},
			},
		},
		{
			"1.5 radius",
			args{
				V2{0, 0},
				1.5,
			},
			[]V2{
				{-1, -1}, {0, -1}, {1, -1},
				{-1, 0}, {0, 0}, {1, 0},
				{-1, 1}, {0, 1}, {1, 1},
			},
		},
		{
			"offset",
			args{
				V2{10, 10},
				1.5,
			},
			[]V2{
				{9, 9}, {10, 9}, {11, 9},
				{9, 10}, {10, 10}, {11, 10},
				{9, 11}, {10, 11}, {11, 11},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.wantRes, CircleAround(tt.args.center, tt.args.radius))
		})
	}
}

var Result interface{}

func Benchmark_circleAround(b *testing.B) {
	for _, radius := range []int{1000} {
		b.Run(fmt.Sprintf("radius=%d", radius), _benchCircle(CircleAround, float64(radius)))
	}
}

func _benchCircle(fn func(V2, float64) []V2, radius float64) func(*testing.B) {
	return func(b *testing.B) {
		var r []V2

		for i := 0; i < b.N; i++ {
			r = fn(V2{0, 0}, radius)
		}

		Result = r
	}
}
