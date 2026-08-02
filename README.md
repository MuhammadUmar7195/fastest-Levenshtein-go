<p align="center">
  <img src="images/logo.svg" alt="Fastest Levenshtein Go Logo" width="300">
</p>

# Fastest Levenshtein (Go Port)

> A fast, bit-parallel Levenshtein distance implementation in Go. Ported from [ka-weihe/fastest-levenshtein](https://github.com/ka-weihe/fastest-levenshtein) and optimized for Go's concurrency and memory model.

[![Port Mortem 2026](https://img.shields.io/badge/Hackathon-Port%20Mortem%202026-red)](https://coderesurrection.com/2026)
[![Repo pool](https://img.shields.io/badge/Repo%20pool-ka--weihe%2Ffastest--levenshtein-blue)](https://github.com/ka-weihe/fastest-levenshtein)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.md)

---

## Improvements Over the Original

1. **Myers' Bit-Parallel Algorithm**: Translated directly from the original repo, ensuring the bitwise math perfectly matches JavaScript's behavior.
2. **Concurrent Batch Processing (`ClosestParallel`)**: Uses Go goroutines and channels to split large array searches across multiple CPU cores, which the original single-threaded version doesn't do.
3. **Smaller Memory Footprint**: Changed the lookup table from indexing 16-bit characters (which takes 256 KiB per call) to indexing bytes (which takes just 1 KiB).
4. **64-bit Support**: Updated the core bitwise engine to use Go's native `uint64`. This lets us process strings up to 64 characters long entirely on the stack (0 allocations), doubling the original 32-character limit.
5. **Prefix and Suffix Trimming**: We scan and remove identical start and end characters before running the main algorithm. This saves a lot of time on real-world strings like file paths and URLs.

---

## Benchmarks & Performance

### 1. Direct Comparison: Go Port vs. Original TypeScript (`ka-weihe/fastest-levenshtein`)

Benchmarked on an Intel i5-6200U @ 2.30GHz using the original methodology: 1,000 random strings of length $N$, distance computed across 500 consecutive pairs, reported as ops/sec.

| Input Size ($N$) | Go Port (ops/sec) | Original TypeScript (ops/sec) | Speedup   |
| ---------------- | ----------------- | ----------------------------- | --------- |
| $N=4$            | 18,994            | 13,575                        | **1.40x** |
| $N=8$            | 15,039            | 4,892                         | **3.07x** |
| $N=16$           | 8,259             | 3,260                         | **2.53x** |
| $N=32$           | 5,227             | 1,651                         | **3.17x** |
| $N=64$           | 2,305             | 252                           | **9.14x** |
| $N=128$          | 205               | 129                           | **1.58x** |
| $N=256$          | 59                | 35                            | **1.68x** |
| $N=512$          | 12                | 8                             | **1.45x** |
| $N=1024$         | 4                 | 1                             | **3.88x** |

<p align="center">
   <img src="images/benchmark.svg" alt="Levenshtein-go vs TypeScript Benchmark Chart" width="100%">
</p>

---

### 2. Comparison with Other Go Libraries

Here is how this port compares against other major Go Levenshtein libraries based on algorithm type, memory allocations, and multi-core support:

| Library / Implementation                                                                          | Language / Stack     | Primary Algorithm                  | Est. Relative Speed (N=32) | Memory Allocations ($\le 32$ chars) | Multi-Core Parallel Processing |   UTF-8 / Rune Support   |
| :------------------------------------------------------------------------------------------------ | :------------------- | :--------------------------------- | :------------------------: | :---------------------------------: | :----------------------------: | :----------------------: |
| **`MuhammadUmar7195/fastest-Levenshtein-go` (This Port)**                                         | **Go 1.21+**         | **Myers' Bit-Parallel (`uint32`)** |    **1.00x (Baseline)**    |             **0 B/op**              |  **Yes (`ClosestParallel`)**   | Byte-level (ASCII exact) |
| [`ka-weihe/fastest-levenshtein`](https://github.com/ka-weihe/fastest-levenshtein) (Repo pool repo)| TypeScript / Node.js | Myers' Bit-Parallel (V8)           |   ~0.31x (3.17x slower)    |           V8 GC Overhead            |      No (Single-threaded)      |    UTF-16 Code Units     |
| `agnivade/levenshtein`                                                                            | Go                   | Standard Dynamic Programming       |    ~0.16x (~6x slower)     |            $O(N)$ Slices            |               No               |     Full UTF-8 Runes     |
| `gnames/levenshtein`                                                                              | Go                   | Myers' Bit-Parallel Variant        |   ~0.60x (~1.6x slower)    |              Low Heap               |               No               |       UTF-8 Runes        |
| `eaxis/levenshtein`                                                                               | Go                   | Matrix Dynamic Programming         |    ~0.08x (~12x slower)    |           Dynamic Slices            |               No               |       Byte / Rune        |
| `pollen5/go-levenshtein`                                                                          | Go                   | Wagner-Fischer Matrix              |    ~0.07x (~14x slower)    |       $O(M \times N)$ Slices        |               No               |       Byte / Rune        |

> **Key Takeaways:**
>
> - **No Heap Allocations (up to 64 chars)**: By using `uint64` and a smaller byte table, we avoid heap allocations for strings under 65 characters.
> - **Prefix/Suffix Trimming**: Stripping matching characters from the ends makes a big difference for real-world strings before the main algorithm even runs.
> - **Bit-Parallel Matrix**: Computes up to 64 matrix cells in a single 64-bit bitwise operation.
> - **Scalability**: `ClosestParallel` breaks large lists into smaller chunks and processes them across multiple cores.

---

## Quick Start & Usage

```go
package main

import (
 "fmt"
 "github.com/MuhammadUmar7195/fastest-Levenshtein-go"
)

func main() {
 // Calculate distance
 dist := levenshtein.Distance("fast", "faster")
 fmt.Println("Distance:", dist)

 // Find closest string
 closest := levenshtein.Closest("fast", []string{"slow", "faster", "fastest"})
 fmt.Println("Closest:", closest)
}
```

---

## Testing & Verification

We tested this directly against the original TypeScript code. The test suite compiles the actual `ka-weihe/fastest-levenshtein` code using esbuild and runs it through Node.

It compares the results against our Go port across **~29,000 randomized cases** (including edge cases like lengths of 31/32/33, empty strings, and long 4 KB strings). The behavior is 100% identical.

```bash
# Full test suite + coverage
go test -v -cover ./...

# Differential test against the real TypeScript implementation
go test -run TestDifferentialAgainstOriginal -v
```

> ![Test Coverage Screenshot](images/test-coverage.svg)

---

## Documentation

See [DECISIONS.md](DECISIONS.md) for more details on the engineering trade-offs, bitwise shift differences between Go and JS, and the fuzzing strategy.
