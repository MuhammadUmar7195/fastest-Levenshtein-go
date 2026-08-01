package levenshtein

import (
	"math"
	"runtime"
)

// peqSize is the number of distinct byte values in a string.
//
// The original TypeScript implementation indexes its equality-bit table by
// UTF-16 code units (0..0xFFFF), which requires a 0x10000-element array.
// Go strings are byte sequences, so byte values are always in 0..255.
// Indexing by byte therefore shrinks the hot lookup table from 256 KiB to
// 1 KiB (a 256x reduction), which eliminates the per-call allocation that
// dominated the original port's memory profile.
const peqSize = 256

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

// Similarity returns a normalized similarity score in [0, 1] between two
// strings, where 1 means identical and 0 means completely different.
//
// It is a convenience layer on top of Distance that the original
// TypeScript package does not provide.
func Similarity(a, b string) float64 {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	if maxLen == 0 {
		return 1
	}
	return 1 - float64(Distance(a, b))/float64(maxLen)
}

// Closest finds the closest string in an array to the target string.
func Closest(target string, array []string) string {
	if len(array) == 0 {
		return ""
	}
	minDistance := math.MaxInt
	minIndex := 0
	for i, s := range array {
		dist := Distance(target, s)
		if dist < minDistance {
			minDistance = dist
			minIndex = i
		}
	}
	return array[minIndex]
}

// ClosestParallel finds the closest string in an array concurrently using goroutines.
// This is an architectural optimization leveraging Go's multi-core concurrency model.
func ClosestParallel(target string, array []string) string {
	if len(array) == 0 {
		return ""
	}
	if len(array) < 200 {
		return Closest(target, array)
	}

	numWorkers := runtime.NumCPU()
	chunkSize := (len(array) + numWorkers - 1) / numWorkers

	type workerResult struct {
		minDist int
		minIdx  int
	}

	results := make(chan workerResult, numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(array) {
			end = len(array)
		}

		go func(subset []string, offset int) {
			minDist := math.MaxInt
			minIdx := 0
			for i, s := range subset {
				dist := Distance(target, s)
				if dist < minDist {
					minDist = dist
					minIdx = offset + i
				}
			}
			results <- workerResult{minDist, minIdx}
		}(array[start:end], start)
	}

	globalMinDist := math.MaxInt
	globalMinIdx := 0

	for i := 0; i < numWorkers; i++ {
		res := <-results
		if res.minDist < globalMinDist {
			globalMinDist = res.minDist
			globalMinIdx = res.minIdx
		}
	}

	return array[globalMinIdx]
}

func myers32(a, b string) int {
	var equalityBits [peqSize]uint32
	aLength := len(a)
	bLength := len(b)
	lastBit := uint32(1 << uint(aLength-1))
	var plusVector uint32 = 0xFFFFFFFF
	var minusVector uint32 = 0
	score := aLength
	i := aLength

	for i > 0 {
		i--
		equalityBits[a[i]] |= 1 << uint(i)
	}

	for i = 0; i < bLength; i++ {
		equality := equalityBits[b[i]]
		xVector := equality | minusVector
		equality |= ((equality & plusVector) + plusVector) ^ plusVector
		minusVector |= ^(equality | plusVector)
		plusVector &= equality
		if (minusVector & lastBit) != 0 {
			score++
		}
		if (plusVector & lastBit) != 0 {
			score--
		}
		minusVector = (minusVector << 1) | 1
		plusVector = (plusVector << 1) | ^(xVector | minusVector)
		minusVector &= xVector
	}

	return score
}

func myersX(longer, shorter string) int {
	var equalityBits [peqSize]uint32
	shorterLength := len(shorter)
	longerLength := len(longer)
	horizontalSize := (shorterLength + 31) / 32
	verticalSize := (longerLength + 31) / 32
	plusCarry := make([]uint32, horizontalSize)
	minusCarry := make([]uint32, horizontalSize)

	for i := 0; i < horizontalSize; i++ {
		plusCarry[i] = 0xFFFFFFFF
		minusCarry[i] = 0
	}

	j := 0
	for ; j < verticalSize-1; j++ {
		var minusVector uint32 = 0
		var plusVector uint32 = 0xFFFFFFFF
		start := j * 32
		verticalEnd := int(math.Min(32, float64(longerLength))) + start
		for k := start; k < verticalEnd; k++ {
			equalityBits[longer[k]] |= 1 << uint(k&31)
		}
		for i := 0; i < shorterLength; i++ {
			equality := equalityBits[shorter[i]]
			plusBit := (plusCarry[i/32] >> uint(i&31)) & 1
			minusBit := (minusCarry[i/32] >> uint(i&31)) & 1
			xVector := equality | minusVector
			xHorizontal := ((((equality | minusBit) & plusVector) + plusVector) ^ plusVector) | equality | minusBit
			plusHorizontal := minusVector | ^(xHorizontal | plusVector)
			minusHorizontal := plusVector & xHorizontal
			if ((plusHorizontal >> 31) ^ plusBit) != 0 {
				plusCarry[i/32] ^= 1 << uint(i&31)
			}
			if ((minusHorizontal >> 31) ^ minusBit) != 0 {
				minusCarry[i/32] ^= 1 << uint(i&31)
			}
			plusHorizontal = (plusHorizontal << 1) | plusBit
			minusHorizontal = (minusHorizontal << 1) | minusBit
			plusVector = minusHorizontal | ^(xVector | plusHorizontal)
			minusVector = plusHorizontal & xVector
		}
		for k := start; k < verticalEnd; k++ {
			equalityBits[longer[k]] = 0
		}
	}

	var minusVector uint32 = 0
	var plusVector uint32 = 0xFFFFFFFF
	start := j * 32
	verticalEnd := int(math.Min(32, float64(longerLength-start))) + start
	for k := start; k < verticalEnd; k++ {
		equalityBits[longer[k]] |= 1 << uint(k&31)
	}
	score := longerLength
	for i := 0; i < shorterLength; i++ {
		equality := equalityBits[shorter[i]]
		plusBit := (plusCarry[i/32] >> uint(i&31)) & 1
		minusBit := (minusCarry[i/32] >> uint(i&31)) & 1
		xVector := equality | minusVector
		xHorizontal := ((((equality | minusBit) & plusVector) + plusVector) ^ plusVector) | equality | minusBit
		plusHorizontal := minusVector | ^(xHorizontal | plusVector)
		minusHorizontal := plusVector & xHorizontal
		score += int((plusHorizontal >> uint((longerLength-1)&31)) & 1)
		score -= int((minusHorizontal >> uint((longerLength-1)&31)) & 1)
		if ((plusHorizontal >> 31) ^ plusBit) != 0 {
			plusCarry[i/32] ^= 1 << uint(i&31)
		}
		if ((minusHorizontal >> 31) ^ minusBit) != 0 {
			minusCarry[i/32] ^= 1 << uint(i&31)
		}
		plusHorizontal = (plusHorizontal << 1) | plusBit
		minusHorizontal = (minusHorizontal << 1) | minusBit
		plusVector = minusHorizontal | ^(xVector | plusHorizontal)
		minusVector = plusHorizontal & xVector
	}
	for k := start; k < verticalEnd; k++ {
		equalityBits[longer[k]] = 0
	}
	return score
}
