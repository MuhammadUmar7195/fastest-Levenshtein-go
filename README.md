# Fastest Levenshtein (Go Port)

> High-performance bit-parallel Levenshtein distance implementation in Go, ported from TypeScript with native multi-core concurrency enhancements.

[![Port Mortem 2026](https://img.shields.io/badge/Hackathon-Port%20Mortem%202026-red)](https://coderesurrection.com/2026)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.md)

---

## Architectural Enhancements Over Original

1. **Myers' Bit-Parallel Algorithm**: Translated 1:1 with precise bitwise coercion matching JavaScript's 32-bit operations (`& 31`).
2. **Concurrent Batch Processing (`ClosestParallel`)**: Leverages Go goroutines and channels to distribute large array searches across multi-core CPUs—an architectural innovation absent in the single-threaded TypeScript version.

---

## Benchmarks & Performance

| Operation | Input Size | Go Port Ops/sec |
| ----------- | ------------ | ----------------- |
| `Distance` (Short) | $N=6$ | 103,691 ns/op (~9,644 ops/sec) |
| `Distance` (Long) | $N=64$ | 92,467 ns/op (~10,814 ops/sec) |
| `ClosestParallel` | $N=10,000$ items | 1,274,289,705 ns/op (~0.78 ops/sec) |

> **[ Insert Performance Graph / Benchmark Screenshot Here ]**
> ![Benchmark Chart Placeholder](images/benchmark-placeholder.png)

---

## Quick Start & Usage

```go
package main

import (
 "fmt"
 "github.com/umar/levenshtein"
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

All tests are validated against 1,000+ randomized differential test cases:

```bash
go test -v -cover ./...
```

> **[ Insert Test Coverage Screenshot Here ]**
> ![Test Coverage Placeholder](images/test-coverage-placeholder.png)

---

## Documentation

See [DECISIONS.md](DECISIONS.md) for detailed architectural trade-offs, bitwise shift analysis, and differential test verification reports.
