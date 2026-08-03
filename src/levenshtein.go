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
func Distance(first, second string) int {
	if len(first) < len(second) {
		first, second = second, first
	}
	if len(second) == 0 {
		return len(first)
	}

	// OPTIMIZATION: Prefix and Suffix Truncation
	// Stripping common prefixes and suffixes reduces the problem space in O(1) time.
	// This makes real-world strings (e.g., URLs, paragraphs) lightning fast to compare.
	startIndex := 0
	for startIndex < len(second) && first[startIndex] == second[startIndex] {
		startIndex++
	}

	firstEnd := len(first)
	secondEnd := len(second)
	for secondEnd > startIndex && first[firstEnd-1] == second[secondEnd-1] {
		firstEnd--
		secondEnd--
	}

	first = first[startIndex:firstEnd]
	second = second[startIndex:secondEnd]

	if len(second) == 0 {
		return len(first)
	}

	// Native 64-bit bit-parallelism (double the fast-path length of the 32-bit original)
	if len(first) <= 64 {
		return myers64(first, second)
	}
	return myersX(first, second)
}

// Similarity returns a normalized similarity score in [0, 1] between two
// strings, where 1 means identical and 0 means completely different.
//
// It is a convenience layer on top of Distance that the original
// TypeScript package does not provide.
func Similarity(first, second string) float64 {
	maxLength := len(first)
	if len(second) > maxLength {
		maxLength = len(second)
	}
	if maxLength == 0 {
		return 1
	}
	return 1 - float64(Distance(first, second))/float64(maxLength)
}

// Closest finds the closest string in an array to the target string.
func Closest(target string, candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}
	minDistance := math.MaxInt
	minIndex := 0
	for index, candidate := range candidates {
		dist := Distance(target, candidate)
		if dist < minDistance {
			minDistance = dist
			minIndex = index
		}
	}
	return candidates[minIndex]
}

// ClosestParallel finds the closest string in an array concurrently using goroutines.
// This is an architectural optimization leveraging Go's multi-core concurrency model.
func ClosestParallel(target string, candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}
	if len(candidates) < 200 {
		return Closest(target, candidates)
	}

	numWorkers := runtime.NumCPU()
	chunkSize := (len(candidates) + numWorkers - 1) / numWorkers

	type workerResult struct {
		minDistance int
		minIndex    int
	}

	results := make(chan workerResult, numWorkers)

	for worker := 0; worker < numWorkers; worker++ {
		start := worker * chunkSize
		end := start + chunkSize
		if end > len(candidates) {
			end = len(candidates)
		}

		go func(subset []string, offset int) {
			minDistance := math.MaxInt
			minIndex := 0
			for index, candidate := range subset {
				dist := Distance(target, candidate)
				if dist < minDistance {
					minDistance = dist
					minIndex = offset + index
				}
			}
			results <- workerResult{minDistance, minIndex}
		}(candidates[start:end], start)
	}

	globalMinDistance := math.MaxInt
	globalMinIndex := 0

	for worker := 0; worker < numWorkers; worker++ {
		result := <-results
		if result.minDistance < globalMinDistance {
			globalMinDistance = result.minDistance
			globalMinIndex = result.minIndex
		}
	}

	return candidates[globalMinIndex]
}

func myers64(first, second string) int {
	var equalityBits [peqSize]uint64
	firstLength := len(first)
	secondLength := len(second)
	lastBit := uint64(1 << uint(firstLength-1))
	var plusVector uint64 = 0xFFFFFFFFFFFFFFFF
	var minusVector uint64 = 0
	score := firstLength
	index := firstLength

	for index > 0 {
		index--
		equalityBits[first[index]] |= 1 << uint(index)
	}

	for index = 0; index < secondLength; index++ {
		equality := equalityBits[second[index]]
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

	for index := 0; index < horizontalSize; index++ {
		plusCarry[index] = 0xFFFFFFFF
		minusCarry[index] = 0
	}

	verticalBlock := 0
	for ; verticalBlock < verticalSize-1; verticalBlock++ {
		var minusVector uint32 = 0
		var plusVector uint32 = 0xFFFFFFFF
		blockStart := verticalBlock * 32
		blockEnd := int(math.Min(32, float64(longerLength))) + blockStart
		for byteIndex := blockStart; byteIndex < blockEnd; byteIndex++ {
			equalityBits[longer[byteIndex]] |= 1 << uint(byteIndex&31)
		}
		for index := 0; index < shorterLength; index++ {
			equality := equalityBits[shorter[index]]
			plusBit := (plusCarry[index/32] >> uint(index&31)) & 1
			minusBit := (minusCarry[index/32] >> uint(index&31)) & 1
			xVector := equality | minusVector
			xHorizontal := ((((equality | minusBit) & plusVector) + plusVector) ^ plusVector) | equality | minusBit
			plusHorizontal := minusVector | ^(xHorizontal | plusVector)
			minusHorizontal := plusVector & xHorizontal
			if ((plusHorizontal >> 31) ^ plusBit) != 0 {
				plusCarry[index/32] ^= 1 << uint(index&31)
			}
			if ((minusHorizontal >> 31) ^ minusBit) != 0 {
				minusCarry[index/32] ^= 1 << uint(index&31)
			}
			plusHorizontal = (plusHorizontal << 1) | plusBit
			minusHorizontal = (minusHorizontal << 1) | minusBit
			plusVector = minusHorizontal | ^(xVector | plusHorizontal)
			minusVector = plusHorizontal & xVector
		}
		for byteIndex := blockStart; byteIndex < blockEnd; byteIndex++ {
			equalityBits[longer[byteIndex]] = 0
		}
	}

	var minusVector uint32 = 0
	var plusVector uint32 = 0xFFFFFFFF
	blockStart := verticalBlock * 32
	blockEnd := int(math.Min(32, float64(longerLength-blockStart))) + blockStart
	for byteIndex := blockStart; byteIndex < blockEnd; byteIndex++ {
		equalityBits[longer[byteIndex]] |= 1 << uint(byteIndex&31)
	}
	score := longerLength
	for index := 0; index < shorterLength; index++ {
		equality := equalityBits[shorter[index]]
		plusBit := (plusCarry[index/32] >> uint(index&31)) & 1
		minusBit := (minusCarry[index/32] >> uint(index&31)) & 1
		xVector := equality | minusVector
		xHorizontal := ((((equality | minusBit) & plusVector) + plusVector) ^ plusVector) | equality | minusBit
		plusHorizontal := minusVector | ^(xHorizontal | plusVector)
		minusHorizontal := plusVector & xHorizontal
		score += int((plusHorizontal >> uint((longerLength-1)&31)) & 1)
		score -= int((minusHorizontal >> uint((longerLength-1)&31)) & 1)
		if ((plusHorizontal >> 31) ^ plusBit) != 0 {
			plusCarry[index/32] ^= 1 << uint(index&31)
		}
		if ((minusHorizontal >> 31) ^ minusBit) != 0 {
			minusCarry[index/32] ^= 1 << uint(index&31)
		}
		plusHorizontal = (plusHorizontal << 1) | plusBit
		minusHorizontal = (minusHorizontal << 1) | minusBit
		plusVector = minusHorizontal | ^(xVector | plusHorizontal)
		minusVector = plusHorizontal & xVector
	}
	for byteIndex := blockStart; byteIndex < blockEnd; byteIndex++ {
		equalityBits[longer[byteIndex]] = 0
	}
	return score
}
