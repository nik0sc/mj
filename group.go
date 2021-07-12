package mj

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Group is an allocation of tiles in a hand to melds.
type Group struct {
	// Each tile represents a meld of 3 identical tiles.
	Pengs Hand
	// Each tile is the first of 3 consecutive tiles.
	Chis Hand
	// Each tile represents a pair.
	Pairs Hand
	// All the leftover tiles.
	Free Hand
}

// ToHand expands the Pengs, Chis and Pairs into their full tile sequences and
// recreates the original Hand.
func (g Group) ToHand() Hand {
	var h Hand

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
func (g Group) ToCount() Counter {
	m := make(map[Tile]int)

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

	cnt, err := NewCounter(m)
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
		gNew.Pengs = make(Hand, len(g.Pengs))
		copy(gNew.Pengs, g.Pengs)
	}

	if g.Chis != nil {
		gNew.Chis = make(Hand, len(g.Chis))
		copy(gNew.Chis, g.Chis)
	}

	if g.Pairs != nil {
		gNew.Pairs = make(Hand, len(g.Pairs))
		copy(gNew.Pairs, g.Pairs)
	}

	if g.Free != nil {
		gNew.Free = make(Hand, len(g.Free))
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

	// A winning hand (including the 14th tile) when grouped has a score of 18.
	// If waiting to complete a pair, the score is 4*4 = 16.
	// If waiting to complete a peng, the score is 4*3 + 2*2 = also 16.
	// If waiting to complete a chi, the score is 4*3 + 2*1 = 14.
	//
	// 7 pairs has a lower score than a winning hand, and 6 pairs has a lower score than
	// any waiting hand. This may be fixed in a future scoring algorithm.
	// Maybe some non-linearity in the growth of the score's pair portion?

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

// UnmarshalGroup is the inverse of Group.Marshal().
func UnmarshalGroup(repr string) Group {
	var g Group

	reprs := strings.Split(repr, ",")
	if len(reprs) != 4 {
		panic(fmt.Sprintf("wrong number of fields: %d", len(reprs)))
	}

	g.Pengs = UnmarshalHand(reprs[0])
	g.Chis = UnmarshalHand(reprs[1])
	g.Pairs = UnmarshalHand(reprs[2])
	g.Free = UnmarshalHand(reprs[3])

	return g
}
