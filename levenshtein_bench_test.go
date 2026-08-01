package levenshtein

import (
	"fmt"
	"testing"
)

// benchmarkPairs generates 1000 random strings of length n and computes
// distance across 500 consecutive pairs, mirroring the original
// fastest-levenshtein bench.ts methodology so ops/sec are directly
// comparable to the published TypeScript numbers.
func benchmarkPairs(b *testing.B, n int) {
	const arrSize = 1000
	arr := make([]string, arrSize)
	for i := range arr {
		arr[i] = makeID(n)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < arrSize-1; j += 2 {
			Distance(arr[j], arr[j+1])
		}
	}
}

func BenchmarkPairs4(b *testing.B)    { benchmarkPairs(b, 4) }
func BenchmarkPairs8(b *testing.B)    { benchmarkPairs(b, 8) }
func BenchmarkPairs16(b *testing.B)   { benchmarkPairs(b, 16) }
func BenchmarkPairs32(b *testing.B)   { benchmarkPairs(b, 32) }
func BenchmarkPairs64(b *testing.B)   { benchmarkPairs(b, 64) }
func BenchmarkPairs128(b *testing.B)  { benchmarkPairs(b, 128) }
func BenchmarkPairs256(b *testing.B)  { benchmarkPairs(b, 256) }
func BenchmarkPairs512(b *testing.B)  { benchmarkPairs(b, 512) }
func BenchmarkPairs1024(b *testing.B) { benchmarkPairs(b, 1024) }

func BenchmarkClosestParallel(b *testing.B) {
	target := "benchmark"
	arr := make([]string, 10000)
	for i := range arr {
		arr[i] = fmt.Sprintf("string-%d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClosestParallel(target, arr)
	}
}
