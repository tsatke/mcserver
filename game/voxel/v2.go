package voxel

import "fmt"

type V2 struct {
	X int
	Z int
}

func (v V2) String() string { return fmt.Sprintf("(%d,%d)", v.X, v.Z) }

func (v V2) Add(other V2) V2 {
	return V2{v.X + other.X, v.Z + other.Z}
}

func (v V2) Sub(other V2) V2 {
	return V2{v.X - other.X, v.Z - other.Z}
}
