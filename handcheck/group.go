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

// Group is an allocation of tiles in a hand to melds.
type Group struct {
	// Each tile represents a meld of 3 identical tiles.
	Pengs mj.Hand
	// Each tile is the first of 3 consecutive tiles.
	Chis mj.Hand
	// Each tile represents a pair.
	Pairs mj.Hand
	// All the leftover tiles.
	Free mj.Hand
}

// ToHand expands the Pengs, Chis and Pairs into their full tile sequences and
// recreates the original Hand.
func (g Group) ToHand() mj.Hand {
	var h mj.Hand

	for _, t := range g.Pengs {
		for i := 0; i < 3; i++ {
			h = append(h, t)
		}
	}

	for _, t := range g.Chis {
		h = append(h, t)
		t.Value++
		h = append(h, t)
		t.Value++
		h = append(h, t)
	}

	for _, t := range g.Pairs {
		for i := 0; i < 2; i++ {
			h = append(h, t)
		}
	}

	h = append(h, g.Free...)
	return h
}

// ToCount expands the Pengs, Chis and Pairs into their full tile sequences and
// returns a Counter of the full hand.
func (g Group) ToCount() mj.Counter {
	m := make(map[mj.Tile]int)

	for _, t := range g.Pengs {
		m[t] += 3
	}

	for _, t := range g.Chis {
		m[t]++
		t.Value++
		m[t]++
		t.Value++
		m[t]++
	}

	for _, t := range g.Pairs {
		m[t] += 2
	}

	for _, t := range g.Free {
		m[t]++
	}

	cnt, err := mj.NewCounter(m)
	if err != nil {
		panic("cannot build counter from result: " + err.Error())
	}

	return cnt
}

// String returns the human-readable representation of this Group, in the order
// Pengs, Chis, Pairs and Free.
func (g Group) String() string {
	return g.ToHand().String()
}

// Marshal returns a space-efficient encoding of this Group, suitable for comparison
// and map keys. For a stable representation, sort each field first.
func (g Group) Marshal() string {
	var b bytes.Buffer

	for _, t := range g.Pengs {
		b.WriteByte(t.Marshal())
	}
	b.WriteByte(',')

	for _, t := range g.Chis {
		b.WriteByte(t.Marshal())
	}
	b.WriteByte(',')

	for _, t := range g.Pairs {
		b.WriteByte(t.Marshal())
	}
	b.WriteByte(',')

	for _, t := range g.Free {
		b.WriteByte(t.Marshal())
	}

	return b.String()
}

// Copy deep-copies the Group. The new Group's fields may also be sorted.
func (g Group) Copy(sorted bool) Group {
	var gNew Group

	if g.Pengs != nil {
		gNew.Pengs = make(mj.Hand, len(g.Pengs))
		copy(gNew.Pengs, g.Pengs)
	}

	if g.Chis != nil {
		gNew.Chis = make(mj.Hand, len(g.Chis))
		copy(gNew.Chis, g.Chis)
	}

	if g.Pairs != nil {
		gNew.Pairs = make(mj.Hand, len(g.Pairs))
		copy(gNew.Pairs, g.Pairs)
	}

	if g.Free != nil {
		gNew.Free = make(mj.Hand, len(g.Free))
		copy(gNew.Free, g.Free)
	}

	if sorted {
		gNew.sort()
	}

	return gNew
}

// Score is used to determine the optimality of groupings. A higher
// score is better. This only considers the hand and not the context
// of the surrounding game.
func (g Group) Score() int {
	// b1 b1 b1 b1 b1 b1:
	// - 2 peng: Score=8
	// - 3 pair: Score=6
	// The zero Group has score 0, this is intentional
	// Effectively, a free tile is worth nothing,
	// a tile in a pair is worth 1,
	// and a tile in a peng/chi is worth 1.333...
	return 4*len(g.Pengs) + 4*len(g.Chis) + 2*len(g.Pairs)
	// a good compiler would turn that into left shifts and adds
}

// sort sorts the groups in-place
func (g Group) sort() {
	sort.Sort(g.Pengs)
	sort.Sort(g.Chis)
	sort.Sort(g.Pairs)
	sort.Sort(g.Free)
}

// free derives the value of Free by subtracting the formed groups
// from a map of tiles to counts.
func (g Group) free(cmap map[mj.Tile]int) (mj.Hand, error) {
	for _, t := range g.Pengs {
		cmap[t] -= 3
	}

	for _, t := range g.Pairs {
		cmap[t] -= 2
	}

	for _, t := range g.Chis {
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

// UnmarshalGroup is the inverse of Group.Marshal().
func UnmarshalGroup(repr string) Group {
	var g Group

	reprs := strings.Split(repr, ",")
	if len(reprs) != 4 {
		panic(fmt.Sprintf("wrong number of fields: %d", len(reprs)))
	}

	g.Pengs = mj.UnmarshalHand(reprs[0])
	g.Chis = mj.UnmarshalHand(reprs[1])
	g.Pairs = mj.UnmarshalHand(reprs[2])
	g.Free = mj.UnmarshalHand(reprs[3])

	return g
}
