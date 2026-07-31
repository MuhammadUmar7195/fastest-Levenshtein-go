# Architecture & Implementation Decisions (`DECISIONS.md`)

## Project Overview

Port of `fastest-levenshtein` from TypeScript to Go.

- **Source Language:** TypeScript (`ka-weihe/fastest-levenshtein`)
- **Target Language:** Go (`github.com/umar/levenshtein`)
- **Goal:** High-performance Levenshtein distance calculation with 100% behavioral equivalence and test parity.

---

## Key Implementation Trade-offs

### 1. Algorithm Selection (Myers' Bit-Parallel Algorithm)

Rather than a naive $O(m \times n)$ Wagner-Fischer dynamic programming matrix (which is memory-heavy and slow for larger strings), we ported Myers' bit-parallel algorithm (`myers32` for strings $\le 32$ length and `myersX` for longer strings). 

- **Trade-off:** Bit-parallel algorithms require careful handling of 32-bit unsigned integers (`uint32`) and bitwise operator semantics (`<<`, `>>`, `~`) to match JavaScript's 32-bit integer coercion behavior.

### 2. Bitwise Semantics & JS Interoperability

- **Challenge:** JavaScript bitwise operations automatically coerce numbers to 32-bit signed integers and mask shift amounts by 31 (`& 31`).
- **Solution:** In Go, we explicitly applied bitwise masking (`& 31`) on shift amounts (e.g., `1 << uint(k&31)`) to ensure exact behavioral equivalence with the original TypeScript engine across strings of any length (up to 1000+ characters).

### 3. Memory & Thread Safety

- **Design:** Allocating a local frequency table (`peq [0x10000]uint32`) on the stack/heap per call ensures complete thread safety and prevents cross-call contamination, avoiding race conditions in concurrent execution.

---

## Test Parity & Validation

- **Differential Testing:** We implemented a rigorous test suite (`levenshtein_test.go`) running 1,000 randomized test cases comparing our Go port against a reference Wagner-Fischer Levenshtein implementation across variable-length strings (0 to 1000 characters).
- **Result:** 100% test parity achieved across all 1,000 iterations.
