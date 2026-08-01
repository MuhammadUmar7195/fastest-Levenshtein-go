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

| Operation         | Input Size ($N$) | Execution Time (ns/op) | Go Port Throughput (Ops/sec) |
| ----------------- | ---------------- | ---------------------- | ---------------------------- |
| `Distance`        | $N=4$            | 174,895 ns/op          | **5,717,662 ops/sec**        |
| `Distance`        | $N=8$            | 191,059 ns/op          | **5,233,989 ops/sec**        |
| `Distance`        | $N=16$           | 125,345 ns/op          | **7,977,981 ops/sec**        |
| `Distance`        | $N=32$           | 173,131 ns/op          | **5,775,984 ops/sec**        |
| `Distance`        | $N=64$           | 171,635 ns/op          | **5,826,317 ops/sec**        |
| `Distance`        | $N=128$          | 277,191 ns/op          | **3,607,261 ops/sec**        |
| `Distance`        | $N=256$          | 348,124 ns/op          | **2,872,545 ops/sec**        |
| `Distance`        | $N=512$          | 588,815 ns/op          | **1,698,328 ops/sec**        |
| `Distance`        | $N=1024$         | 1,434,066 ns/op        | **697,320 ops/sec**          |
| `ClosestParallel` | $N=10,000$ items | 2,117,898,665 ns/op    | **0.47 ops/sec**             |

> **Performance Benchmark Chart (Execution Time ns/op vs String Length N):**
> ![Levenshtein-go Benchmark Chart](images/benchmark.svg)

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

All tests are validated against 1,000+ randomized differential test cases:

```bash
go test -v -cover ./...
```

> **Test Execution & Coverage Report (99.3% Package Coverage & 100% CLI Coverage across 1,000+ differential tests):**
> ![Test Coverage Screenshot](images/test-coverage.svg)

---

## Documentation

See [DECISIONS.md](DECISIONS.md) for detailed architectural trade-offs, bitwise shift analysis, and differential test verification reports.
