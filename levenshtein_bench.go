package levenshtein

import (
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