package handcheck

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nik0sc/mj"
)

// TODO: move to top level package, rename to "Grouping" or similar

// Result is the optimal solution found by OptChecker and friends.
type Result struct {
	// Each tile represents a meld of 3 identical tiles.
	Pengs mj.Hand
	// Each tile is the first of 3 consecutive tiles.
	Chis mj.Hand
	// Each tile represents a pair.
	Pairs mj.Hand
	// All the leftover tiles.
	Free mj.Hand
}

func (r Result) ToHand() mj.Hand {
	var h mj.Hand

	for _, t := range r.Pengs {
		for i := 0; i < 3; i++ {
			h = append(h, t)
		}
	}

	for _, t := range r.Chis {
		h = append(h, t)
		t.Value++
		h = append(h, t)
		t.Value++
		h = append(h, t)
	}

	for _, t := range r.Pairs {
		for i := 0; i < 2; i++ {
			h = append(h, t)
		}
	}

	h = append(h, r.Free...)
	return h
}

func (r Result) ToCount() mj.Counter {
	m := make(map[mj.Tile]int)

	for _, t := range r.Pengs {
		m[t] += 3
	}

	for _, t := range r.Chis {
		m[t]++
		t.Value++
		m[t]++
		t.Value++
		m[t]++
	}

	for _, t := range r.Pairs {
		m[t] += 2
	}

	for _, t := range r.Free {
		m[t]++
	}

	cnt, err := mj.NewCounter(m)
	if err != nil {
		panic("cannot build counter from result: " + err.Error())
	}

	return cnt
}

// String returns the human-readable representation of this result, in the order
// Pengs, Chis, Pairs and Free.
func (r Result) String() string {
	return r.ToHand().String()
}

// Marshal returns a space-efficient encoding of this result, suitable for comparison
// and map keys. For a stable representation, sort each result field first.
func (r Result) Marshal() string {
	var b bytes.Buffer

	for _, t := range r.Pengs {
		b.WriteByte(t.Marshal())
	}
	b.WriteByte(',')

	for _, t := range r.Chis {
		b.WriteByte(t.Marshal())
	}
	b.WriteByte(',')

	for _, t := range r.Pairs {
		b.WriteByte(t.Marshal())
	}
	b.WriteByte(',')

	for _, t := range r.Free {
		b.WriteByte(t.Marshal())
	}

	return b.String()
}

// Copy deep-copies the Result. The new Result's fields may also be sorted.
func (r Result) Copy(sorted bool) Result {
	var rNew Result

	if r.Pengs != nil {
		rNew.Pengs = make(mj.Hand, len(r.Pengs))
		copy(rNew.Pengs, r.Pengs)
	}

	if r.Chis != nil {
		rNew.Chis = make(mj.Hand, len(r.Chis))
		copy(rNew.Chis, r.Chis)
	}

	if r.Pairs != nil {
		rNew.Pairs = make(mj.Hand, len(r.Pairs))
		copy(rNew.Pairs, r.Pairs)
	}

	if r.Free != nil {
		rNew.Free = make(mj.Hand, len(r.Free))
		copy(rNew.Free, r.Free)
	}

	if sorted {
		rNew.sort()
	}

	return rNew
}

// Score is used to determine the optimality of solutions. A higher
// score is better. This only considers the hand and not the context
// of the surrounding game.
func (r Result) Score() int {
	// b1 b1 b1 b1 b1 b1:
	// - 2 peng: Score=8
	// - 3 pair: Score=6
	// The zero Result has score 0, this is intentional
	// Effectively, a free tile is worth nothing,
	// a tile in a pair is worth 1,
	// and a tile in a peng/chi is worth 1.333...
	return 4*len(r.Pengs) + 4*len(r.Chis) + 2*len(r.Pairs)
	// a good compiler would turn that into left shifts and adds
}

// sort sorts the groups in the result in-place
func (r Result) sort() {
	sort.Sort(r.Pengs)
	sort.Sort(r.Chis)
	sort.Sort(r.Pairs)
	sort.Sort(r.Free)
}

// free derives the value of Free by subtracting the formed groups
// from a map of tiles to counts.
func (r Result) free(cmap map[mj.Tile]int) (mj.Hand, error) {
	for _, t := range r.Pengs {
		cmap[t] -= 3
	}

	for _, t := range r.Pairs {
		cmap[t] -= 2
	}

	for _, t := range r.Chis {
		t2 := t
		t2.Value++

		t3 := t
		t3.Value += 2

		cmap[t]--
		cmap[t2]--
		cmap[t3]--
	}

	freecnt, err := mj.NewCounter(cmap)
	if err != nil {
		return nil, errors.New("cannot recreate Counter from map: " + err.Error())
	}
	return freecnt.ToHand(true), nil
}

// UnmarshalResult is the inverse of Result.Marshal().
func UnmarshalResult(repr string) Result {
	var r Result

	reprs := strings.Split(repr, ",")
	if len(reprs) != 4 {
		panic(fmt.Sprintf("wrong number of fields: %d", len(reprs)))
	}

	r.Pengs = mj.UnmarshalHand(reprs[0])
	r.Chis = mj.UnmarshalHand(reprs[1])
	r.Pairs = mj.UnmarshalHand(reprs[2])
	r.Free = mj.UnmarshalHand(reprs[3])

	return r
}
