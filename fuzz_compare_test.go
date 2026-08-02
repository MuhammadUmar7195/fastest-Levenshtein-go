package levenshtein

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// tsRef bundles the esbuild-transpiled original TypeScript module and a
// persistent directory so we can run Node many times without re-transpiling.
type tsRef struct {
	dir string
}

// compileOnce builds the ORIGINAL ka-weihe/fastest-levenshtein mod.ts to
// CommonJS a single time per test process, then reuses the artifact.
var compileOnce = sync.OnceValues(func() (*tsRef, error) {
	srcDir := filepath.Join("testdata", "original")
	modSrc, err := os.ReadFile(filepath.Join(srcDir, "mod.ts"))
	if err != nil {
		return nil, fmt.Errorf("read original mod.ts: %w", err)
	}

	dir, err := os.MkdirTemp("", "fl-ts-ref-*")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "mod.ts"), modSrc, 0o600); err != nil {
		return nil, err
	}

	esbuild := exec.Command("npx", "--yes", "esbuild", "mod.ts", "--bundle", "--format=cjs", "--platform=node", "--outfile=mod.cjs")
	esbuild.Dir = dir
	if out, err := esbuild.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("esbuild: %v: %s", err, out)
	}
	return &tsRef{dir: dir}, nil
})

// distance evaluates the real TypeScript implementation on the given pairs.
func (r *tsRef) distance(pairs [][2]string) ([]int, error) {
	payload, err := json.Marshal(pairs)
	if err != nil {
		return nil, err
	}

	runner := fmt.Sprintf(`
const m = require('./mod.cjs');
const pairs = %s;
const out = pairs.map(p => m.distance(p[0], p[1]));
process.stdout.write(JSON.stringify(out));
`, payload)
	if err := os.WriteFile(filepath.Join(r.dir, "runner.js"), []byte(runner), 0o600); err != nil {
		return nil, err
	}

	cmd := exec.Command("node", "runner.js")
	cmd.Dir = r.dir
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("node: %w", err)
	}

	var res []int
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if len(res) != len(pairs) {
		return nil, fmt.Errorf("result count mismatch: got %d want %d", len(res), len(pairs))
	}
	return res, nil
}

func fuzzString(rng *rand.Rand, maxLen int) string {
	if maxLen <= 1 {
		return ""
	}
	n := rng.Intn(maxLen)
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	out := make([]byte, n)
	for i := range out {
		out[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(out)
}

// TestDifferentialAgainstOriginal compares the Go port against the REAL
// TypeScript implementation on randomized inputs, including edge cases the
// original port's own Wagner-Fischer reference would not expose:
//   - empty strings, single chars, boundary lengths (31/32/33, 63/64/65)
//   - long strings (up to 4096) to exercise myersX multi-block paths
func TestDifferentialAgainstOriginal(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not installed; skipping differential test against original TS")
	}
	ref, err := compileOnce()
	if err != nil {
		t.Skipf("could not compile original TS: %v", err)
	}

	rng := rand.New(rand.NewSource(42))
	pairs := make([][2]string, 0, 30000)

	// Deterministic edge cases.
	edges := []string{"", "a", "ab", "abc", "abcdefghij", strings.Repeat("a", 31), strings.Repeat("a", 32),
		strings.Repeat("a", 33), strings.Repeat("b", 63), strings.Repeat("b", 64),
		strings.Repeat("c", 65), strings.Repeat("ab", 512), strings.Repeat("x", 4096)}
	for _, a := range edges {
		for _, b := range edges {
			pairs = append(pairs, [2]string{a, b})
		}
	}

	// Randomized fuzz across multiple length regimes.
	for i := 0; i < 20000; i++ {
		a := fuzzString(rng, rng.Intn(6))
		b := fuzzString(rng, rng.Intn(6))
		pairs = append(pairs, [2]string{a, b})
	}
	for i := 0; i < 8000; i++ {
		pairs = append(pairs, [2]string{fuzzString(rng, 4096), fuzzString(rng, 4096)})
	}

	// Batch through Node in chunks to bound memory.
	const chunk = 2000
	mismatches := 0
	for start := 0; start < len(pairs); start += chunk {
		end := start + chunk
		if end > len(pairs) {
			end = len(pairs)
		}
		refRes, err := ref.distance(pairs[start:end])
		if err != nil {
			t.Fatalf("reference failed: %v", err)
		}
		for i, pair := range pairs[start:end] {
			got := Distance(pair[0], pair[1])
			if got != refRes[i] {
				mismatches++
				if mismatches <= 10 {
					t.Errorf("MISMATCH Distance(%q, %q) = %d; original TS = %d", pair[0], pair[1], got, refRes[i])
				}
			}
		}
	}

	if mismatches > 0 {
		t.Fatalf("differential mismatch: %d of %d cases disagree with original TS", mismatches, len(pairs))
	}
}
