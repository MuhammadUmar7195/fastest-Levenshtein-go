# Benchmark Methodology

All benchmarks were run on an Intel i5-6200U @ 2.30GHz.
We used the original library's benchmark suite methodology: 
Generate 1000 random strings of length N. Compute the distance for 500 consecutive pairs.
Measure operations per second (ops/sec).

The Go benchmark (go test -bench .) was used to gather memory allocation (B/op) and speed (ns/op).
No heap allocations occur for lengths <= 64 due to the native uint64 optimization.
