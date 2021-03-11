package handcheck

import (
	"bytes"
	"strings"

	"mj"
)

const (
	// Record metrics in shared struct.
	writeMetrics = true
	// Trace execution of each step. Very slow.
	traceSteps = false
)

// GreedyResult is the solution found by GreedyChecker. It contains
// totals of melds and unused tiles, not the actual tiles.
type GreedyResult struct {
	// Ok is true if this hand is optimal (has no unused tiles).
	Ok     bool
	Peng   int
	Chi    int
	Pair   int
	Unused int
}

// Result is the optimal solution found by OptChecker and friends.
type Result struct {
	// Each tile represents a meld of 3 identical tiles.
	Pengs []mj.Tile
	// Each tile is the first of 3 consecutive tiles.
	Chis []mj.Tile
	// Each tile represents a pair.
	Pairs []mj.Tile
	// All the leftover tiles.
	Free mj.Hand
}

// String returns the human-readable representation of this result, in the order
// Pengs, Chis, Pairs and Free.
func (r Result) String() string {
	var ss []string

	for _, t := range r.Pengs {
		ss = append(ss, strings.Repeat(t.String(), 3))
	}

	for _, t := range r.Chis {
		t2 := t
		t2.Value++
		t3 := t2
		t3.Value++
		ss = append(ss, t.String()+t2.String()+t3.String())
	}

	for _, t := range r.Pairs {
		ss = append(ss, strings.Repeat(t.String(), 2))
	}

	ss = append(ss, r.Free.String())

	return strings.Join(ss, " ")
}

// Repr returns a space-efficient encoding of this result, suitable for comparison
// and map keys. For a stable representation, sort each result field first.
func (r Result) Repr() string {
	var b bytes.Buffer

	for _, t := range r.Pengs {
		b.WriteByte(t.Repr())
	}
	b.WriteByte(',')

	for _, t := range r.Chis {
		b.WriteByte(t.Repr())
	}
	b.WriteByte(',')

	for _, t := range r.Pairs {
		b.WriteByte(t.Repr())
	}
	b.WriteByte(',')

	for _, t := range r.Free {
		b.WriteByte(t.Repr())
	}

	return b.String()
}

// score is used to determine the optimality of solutions. It is an
// implementation detail and calling code should not know about it.
func (r Result) score() int {
	// b1 b1 b1 b1 b1 b1:
	// - 2 peng: score=8
	// - 3 pair: score=6
	// The zero Result has score 0, this is intentional
	// Effectively, a free tile is worth nothing,
	// a tile in a pair is worth 1,
	// and a tile in a peng/chi is worth 1.333...
	return 4*len(r.Pengs)+4*len(r.Chis)+2*len(r.Pairs)
	// a good compiler would turn that into left shifts and adds
}
