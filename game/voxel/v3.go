package voxel

import "fmt"

type V3 struct {
	X int
	Y int
	Z int
}

func (v V3) String() string { return fmt.Sprintf("(%d,%d,%d)", v.X, v.Y, v.Z) }
