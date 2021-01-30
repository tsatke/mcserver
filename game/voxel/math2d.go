package voxel

import (
	"math"
)

func CircleAround(center V2, radius float64) []V2 {
	rad := int(radius)
	radSq := radius * radius

	/*
		preallocation gives huge performance boost:

		name                          old time/op  new time/op  delta
		_circleAround/radius=10-16    1.61µs ± 0%  1.03µs ± 0%  -36.11%  (p=1.000 n=1+1)
		_circleAround/radius=20-16    5.61µs ± 0%  3.87µs ± 0%  -31.07%  (p=1.000 n=1+1)
		_circleAround/radius=30-16    12.5µs ± 0%   9.1µs ± 0%  -27.54%  (p=1.000 n=1+1)
		_circleAround/radius=40-16    20.3µs ± 0%  14.7µs ± 0%  -27.37%  (p=1.000 n=1+1)
		_circleAround/radius=50-16    35.3µs ± 0%  22.7µs ± 0%  -35.84%  (p=1.000 n=1+1)
		_circleAround/radius=60-16    49.5µs ± 0%  32.3µs ± 0%  -34.88%  (p=1.000 n=1+1)
		_circleAround/radius=100-16    153µs ± 0%    87µs ± 0%  -43.22%  (p=1.000 n=1+1)
		_circleAround/radius=150-16    335µs ± 0%   187µs ± 0%  -44.19%  (p=1.000 n=1+1)
		_circleAround/radius=1000-16  12.9ms ± 0%   7.5ms ± 0%  -41.44%  (p=1.000 n=1+1)

	*/
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
	newRes := make([]V2, 0, len(res)*4+1)
	for _, r := range res {
		newRes = append(newRes, center.Add(r))
		newRes = append(newRes, center.Add(V2{r.Z, -r.X}))
		newRes = append(newRes, center.Add(V2{-r.X, -r.Z}))
		newRes = append(newRes, center.Add(V2{-r.Z, r.X}))
	}
	newRes = append(newRes, center)
	return newRes
}
