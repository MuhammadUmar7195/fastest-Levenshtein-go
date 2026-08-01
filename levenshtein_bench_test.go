package levenshtein

import (
	"fmt"
	"testing"
)

func benchmarkLen(b *testing.B, n int) {
	s1 := makeID(n)
	s2 := makeID(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Distance(s1, s2)
	}
}

func BenchmarkLen4(b *testing.B)    { benchmarkLen(b, 4) }
func BenchmarkLen8(b *testing.B)    { benchmarkLen(b, 8) }
func BenchmarkLen16(b *testing.B)   { benchmarkLen(b, 16) }
func BenchmarkLen32(b *testing.B)   { benchmarkLen(b, 32) }
func BenchmarkLen64(b *testing.B)   { benchmarkLen(b, 64) }
func BenchmarkLen128(b *testing.B)  { benchmarkLen(b, 128) }
func BenchmarkLen256(b *testing.B)  { benchmarkLen(b, 256) }
func BenchmarkLen512(b *testing.B)  { benchmarkLen(b, 512) }
func BenchmarkLen1024(b *testing.B) { benchmarkLen(b, 1024) }

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
