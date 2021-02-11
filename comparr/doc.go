// Package comparr implements a data structure we call "compact array".
// As opposed to a normal array or slice of int8, int16, int32 or similar,
// a compact array is not bound to the predefined bit-size of an integral
// element. This means, that a compact array can not only hold 8, 16 or 32
// bit values, but also 2, 5 or 17 bit values. The maximum supported bit
// size is 64, the minimum required bit size is 1.
//
// All data is held in an uint64 slice. It is not efficient to use a compact
// array as an alternative to a []int64. It is also not efficient to define
// a large bit size, which 64 is not divisible by. Consider a compact
// array with a bit size of 22. This would mean, that per uint64 in memory,
// 64-(2*22)=20 bit would be unused. For a compact array of size 100, this
// would mean 1000 unused bits, or 125 unused bytes.
//
// A compact array can be created in two different ways. Creating an empty
// array, or converting an existing int slice to a compact array.
// Both approaches require you to set a bit size for the compact array.
//
//	carr := comparr.New(10, 5)
//	carr.Len() // == 10
//	len(carr.Data()) // == 1
//
// The above example shows how to initialize an empty compact array.
// It has a length of 10 and a bit size of 5, which means that in total,
// this compact array requires 10*5=50 bits of memory. This is why the
// length of the underlying uint64 slice is 1.
//
//	carr.Set(0, 7)
//
// In this example, the value at index 0 is set to 7. Since the compact
// array has a bit size of 5, specifying any value greater than (1<<5)-1
// will result in the higher bits getting lost.
//
//	carr.Set(0, 32) // 32 is out of the 5 bit range, MSBs will be ignored
//	carr.Get(0)     // == 0
//
// The other way to create a compact array is to convert an existing int slice
// to a compact array.
//
//	carr := comparr.FromSlice([]int{1, 2, 3}, 2)
//	carr.Len()       // == 3
//	len(carr.Data()) // == 1
//
// Just as with CompArr.Set, values exceeding the given bit size will
// have their MSBs ignored.
package comparr
