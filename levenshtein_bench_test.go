package levenshtein

import (
	"fmt"
	"testing"
)

func BenchmarkDistanceShort(b *testing.B) {
	s1 := "fast"
	s2 := "faster"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Distance(s1, s2)
	}
}

func BenchmarkDistanceLong(b *testing.B) {
	s1 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	s2 := "abcdefghijklmnopqrstuvwxyz9876543210"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Distance(s1, s2)
	}
}

func BenchmarkClosest(b *testing.B) {
	target := "fast"
	arr := []string{"slow", "faster", "fastest", "faste", "fas"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Closest(target, arr)
	}
}

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
