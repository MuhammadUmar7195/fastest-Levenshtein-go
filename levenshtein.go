package levenshtein

import (
	"math"
	"runtime"
)

// Distance calculates the Levenshtein distance between two strings.
func Distance(a, b string) int {
	if len(a) < len(b) {
		a, b = b, a
	}
	if len(b) == 0 {
		return len(a)
	}
	if len(a) <= 32 {
		return myers32(a, b)
	}
	return myersX(a, b)
}

// Closest finds the closest string in an array to the target string.
func Closest(str string, arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	minDistance := math.MaxInt
	minIndex := 0
	for i, s := range arr {
		dist := Distance(str, s)
		if dist < minDistance {
			minDistance = dist
			minIndex = i
		}
	}
	return arr[minIndex]
}

// ClosestParallel finds the closest string in an array concurrently using goroutines.
// This is an architectural optimization leveraging Go's multi-core concurrency model.
func ClosestParallel(str string, arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	if len(arr) < 200 {
		return Closest(str, arr)
	}

	numWorkers := runtime.NumCPU()
	chunkSize := (len(arr) + numWorkers - 1) / numWorkers

	type result struct {
		minDist int
		minIdx  int
	}

	ch := make(chan result, numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		if start >= len(arr) {
			break
		}

		go func(subset []string, offset int) {
			minDist := math.MaxInt
			minIdx := 0
			for i, s := range subset {
				dist := Distance(str, s)
				if dist < minDist {
					minDist = dist
					minIdx = offset + i
				}
			}
			ch <- result{minDist, minIdx}
		}(arr[start:end], start)
	}

	globalMinDist := math.MaxInt
	globalMinIdx := 0
	activeWorkers := 0
	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		if start < len(arr) {
			activeWorkers++
		}
	}

	for i := 0; i < activeWorkers; i++ {
		res := <-ch
		if res.minDist < globalMinDist {
			globalMinDist = res.minDist
			globalMinIdx = res.minIdx
		}
	}

	return arr[globalMinIdx]
}

func myers32(a, b string) int {
	var peq [0x10000]uint32
	n := len(a)
	m := len(b)
	lst := uint32(1 << uint(n-1))
	var pv uint32 = 0xFFFFFFFF
	var mv uint32 = 0
	sc := n
	i := n

	for i > 0 {
		i--
		peq[a[i]] |= 1 << uint(i)
	}

	for i = 0; i < m; i++ {
		eq := peq[b[i]]
		xv := eq | mv
		eq |= ((eq & pv) + pv) ^ pv
		mv |= ^(eq | pv)
		pv &= eq
		if (mv & lst) != 0 {
			sc++
		}
		if (pv & lst) != 0 {
			sc--
		}
		mv = (mv << 1) | 1
		pv = (pv << 1) | ^(xv | mv)
		mv &= xv
	}

	return sc
}

func myersX(longer, shorter string) int {
	var peq [0x10000]uint32
	n := len(shorter)
	m := len(longer)
	hsize := (n + 31) / 32
	vsize := (m + 31) / 32
	phc := make([]uint32, hsize)
	mhc := make([]uint32, hsize)

	for i := 0; i < hsize; i++ {
		phc[i] = 0xFFFFFFFF
		mhc[i] = 0
	}

	j := 0
	for ; j < vsize-1; j++ {
		var mv uint32 = 0
		var pv uint32 = 0xFFFFFFFF
		start := j * 32
		vlen := int(math.Min(32, float64(m))) + start
		for k := start; k < vlen; k++ {
			peq[longer[k]] |= 1 << uint(k&31)
		}
		for i := 0; i < n; i++ {
			eq := peq[shorter[i]]
			pb := (phc[i/32] >> uint(i&31)) & 1
			mb := (mhc[i/32] >> uint(i&31)) & 1
			xv := eq | mv
			xh := ((((eq | mb) & pv) + pv) ^ pv) | eq | mb
			ph := mv | ^(xh | pv)
			mh := pv & xh
			if ((ph >> 31) ^ pb) != 0 {
				phc[i/32] ^= 1 << uint(i&31)
			}
			if ((mh >> 31) ^ mb) != 0 {
				mhc[i/32] ^= 1 << uint(i&31)
			}
			ph = (ph << 1) | pb
			mh = (mh << 1) | mb
			pv = mh | ^(xv | ph)
			mv = ph & xv
		}
		for k := start; k < vlen; k++ {
			peq[longer[k]] = 0
		}
	}

	var mv uint32 = 0
	var pv uint32 = 0xFFFFFFFF
	start := j * 32
	vlen := int(math.Min(32, float64(m-start))) + start
	for k := start; k < vlen; k++ {
		peq[longer[k]] |= 1 << uint(k&31)
	}
	score := m
	for i := 0; i < n; i++ {
		eq := peq[shorter[i]]
		pb := (phc[i/32] >> uint(i&31)) & 1
		mb := (mhc[i/32] >> uint(i&31)) & 1
		xv := eq | mv
		xh := ((((eq | mb) & pv) + pv) ^ pv) | eq | mb
		ph := mv | ^(xh | pv)
		mh := pv & xh
		score += int((ph >> uint((m-1)&31)) & 1)
		score -= int((mh >> uint((m-1)&31)) & 1)
		if ((ph >> 31) ^ pb) != 0 {
			phc[i/32] ^= 1 << uint(i&31)
		}
		if ((mh >> 31) ^ mb) != 0 {
			mhc[i/32] ^= 1 << uint(i&31)
		}
		ph = (ph << 1) | pb
		mh = (mh << 1) | mb
		pv = mh | ^(xv | ph)
		mv = ph & xv
	}
	for k := start; k < vlen; k++ {
		peq[longer[k]] = 0
	}
	return score
}
