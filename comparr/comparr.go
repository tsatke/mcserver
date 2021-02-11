package comparr

import (
	"math"
)

// CompArr is a compressed array, with elements stored at a constant bit size.
// E.g. a CompArr could be used to store multiple 5bit values inside a few 64bit
// values.
type CompArr struct {
	length      int
	bitsz       int
	mask        uint64
	elemsPerInt int
	data        []uint64
}

func New(length, bitSize int) CompArr {
	if bitSize > 64 {
		panic("bit size can't exceed 64")
	}
	if bitSize < 1 {
		panic("bit size can't be < 1")
	}

	elemsPerInt := int(64 / bitSize)
	mask := (uint64(1) << bitSize) - 1
	return CompArr{
		length:      length,
		bitsz:       bitSize,
		mask:        mask,
		elemsPerInt: elemsPerInt,
		data:        make([]uint64, int(math.Ceil(float64(length)/float64(elemsPerInt)))),
	}
}

func FromSlice(slice []int, bitSize int) CompArr {
	arr := New(len(slice), bitSize)
	for i, v := range slice {
		arr.Set(i, v)
	}
	return arr
}

// Set will set the value at the given index to the given value.
// If the given value exceeds the value range allowed by this
// compact array, the MSBs will be ignored.
func (c CompArr) Set(index, value int) {
	arrIndex := len(c.data) - 1 - (index * c.bitsz / 64)
	shiftAmount := (index * c.bitsz) % 64

	c.data[arrIndex] &= ^uint64(c.mask << shiftAmount)
	c.data[arrIndex] |= (uint64(value) & c.mask) << shiftAmount
}

// Get returns the value at the given index in this compact array.
func (c CompArr) Get(index int) int {
	arrIndex := len(c.data) - index*c.bitsz/64 - 1
	shiftAmount := (index * c.bitsz) % 64

	return int((c.data[arrIndex] & uint64(c.mask<<shiftAmount)) >> shiftAmount)
}

// Len returns the amount of values that can be stored in this compact
// array.
func (c CompArr) Len() int {
	return c.length
}

// Data returns the underlying uint64 slice, in which this compact
// array stores its data. The returned slice is not a copy, so handle
// with care.
func (c CompArr) Data() []uint64 {
	return c.data
}
