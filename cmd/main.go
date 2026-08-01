package main

import (
	"fmt"
	"os"

	"github.com/MuhammadUmar7195/fastest-Levenshtein-go"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("=== Fastest Levenshtein CLI ===")
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/main.go distance <word1> <word2>")
		fmt.Println("  go run cmd/main.go closest <target> <word1> <word2> <word3> ...")
		return
	}

	command := os.Args[1]
	switch command {
	case "distance":
		if len(os.Args) < 4 {
			fmt.Println("Error: Please provide two words.")
			return
		}
		word1 := os.Args[2]
		word2 := os.Args[3]
		distance := levenshtein.Distance(word1, word2)
		fmt.Printf("Distance between %q and %q = %d\n", word1, word2, distance)

	case "closest":
		if len(os.Args) < 4 {
			fmt.Println("Error: Please provide a target word and candidates.")
			return
		}
		target := os.Args[2]
		candidates := os.Args[3:]
		best := levenshtein.ClosestParallel(target, candidates)
		fmt.Printf("Closest match to %q in %v = %q\n", target, candidates, best)

	default:
		fmt.Println("Unknown command:", command)
	}
}
