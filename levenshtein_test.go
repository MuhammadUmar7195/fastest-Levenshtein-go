package levenshtein

import (
	"fmt"
	"math/rand"
	"testing"
)

func referenceLevenshtein(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)
	if len(ra) == 0 {
		return len(rb)
	}
	if len(rb) == 0 {
		return len(ra)
	}

	if len(ra) > len(rb) {
		ra, rb = rb, ra
	}

	row := make([]int, len(ra)+1)
	for i := range row {
		row[i] = i
	}

	for i := 1; i <= len(rb); i++ {
		prev := i
		for j := 1; j <= len(ra); j++ {
			var val int
			if rb[i-1] == ra[j-1] {
				val = row[j-1]
			} else {
				min1 := row[j-1] + 1
				min2 := prev + 1
				min3 := row[j] + 1
				val = min1
				if min2 < val {
					val = min2
				}
				if min3 < val {
					val = min3
				}
			}
			row[j-1] = prev
			prev = val
		}
		row[len(ra)] = prev
	}

	return row[len(ra)]
}

func makeID(length int) string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}

func TestCompare(t *testing.T) {
	for i := 0; i < 1000; i++ {
		num1 := rand.Intn(1000)
		num2 := rand.Intn(1000)
		s1 := makeID(num1)
		s2 := makeID(num2)
		actual := Distance(s1, s2)
		expected := referenceLevenshtein(s1, s2)
		if actual != expected {
			t.Errorf("Distance(%q, %q) = %d; expected %d", s1, s2, actual, expected)
		}
	}
}

func TestFind(t *testing.T) {
	actual := Closest("fast", []string{"slow", "faster", "fastest"})
	expected := "faster"
	if actual != expected {
		t.Errorf("Closest = %q; expected %q", actual, expected)
	}

	if Closest("fast", []string{}) != "" {
		t.Errorf("Expected empty string for empty array")
	}
}

func TestSimilarity(t *testing.T) {
	if Similarity("", "") != 1 {
		t.Errorf("Similarity('','') should be 1")
	}
	if Similarity("same", "same") != 1 {
		t.Errorf("Similarity('same','same') should be 1")
	}
	if Similarity("fast", "faster") < 0.5 {
		t.Errorf("Similarity('fast','faster') should be >= 0.5")
	}
	if Similarity("abc", "xyz") != 0 {
		t.Errorf("Similarity('abc','xyz') should be 0")
	}
	if Similarity("a", "") != 0 {
		t.Errorf("Similarity('a','') should be 0")
	}
}

func TestClosestParallel(t *testing.T) {
	// Test small array (< 200 items)
	smallArr := []string{"apple", "banana", "cherry"}
	if ClosestParallel("app", smallArr) != "apple" {
		t.Errorf("ClosestParallel small array failed")
	}

	if ClosestParallel("app", []string{}) != "" {
		t.Errorf("ClosestParallel empty array failed")
	}

	// Test large array (>= 200 items)
	largeArr := make([]string, 500)
	for i := range largeArr {
		largeArr[i] = fmt.Sprintf("word-%d", i)
	}
	largeArr[250] = "targetword"
	if ClosestParallel("target", largeArr) != "targetword" {
		t.Errorf("ClosestParallel large array failed")
	}
}
