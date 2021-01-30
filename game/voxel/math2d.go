package voxel

import (
	"math"
)

func CircleAround(center V2, radius float64) []V2 {
	rad := int(radius)
	radSq := radius * radius

	res := make([]V2, 0, int(math.Ceil((math.Pi*radius*radius)/4)))
	for z := 0; z <= rad; z++ {
		targetX := int(radSq - float64(z*z))
		for x := 1; x*x <= targetX; x++ {
			res = append(res, V2{x, z})
		}
	}

	// results contains one quarter of all results, excluding the center,
	// so we need to rotate all coordinates three times, then add
	// the center
	newRes := make([]V2, len(res)*4+1)
	for i, r := range res {
		newRes[i*4] = center.Add(r)
		newRes[i*4+1] = center.Add(V2{r.Z, -r.X})
		newRes[i*4+2] = center.Add(V2{-r.X, -r.Z})
		newRes[i*4+3] = center.Add(V2{-r.Z, r.X})
	}
	newRes[len(newRes)-1] = center
	return newRes
}
