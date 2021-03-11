package mj

import (
	"bytes"
	"fmt"
	"sort"
)

// Hand is an ordered sequence of tiles, representing a mahjong hand.
// Except for Swap, the methods of Hand are guaranteed to not mutate the sequence.
type Hand []Tile

// Valid returns true if all the tiles are valid.
func (h Hand) Valid() bool {
	for _, t := range h {
		if !t.Valid() {
			return false
		}
	}
	return true
}

// See sort.Interface.
func (h Hand) Len() int {
	return len(h)
}

// See sort.Interface.
func (h Hand) Less(i, j int) bool {
	if h[i].Suit != h[j].Suit {
		return h[i].Suit < h[j].Suit
	}
	return h[i].Value < h[j].Value
}

// See sort.Interface.
func (h Hand) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// sort returns the Hand directly if it is sorted, or else it returns a sorted copy of the Hand.
func (h Hand) sort() Hand {
	if sort.IsSorted(h) {
		return h
	}
	hSorted := make(Hand, len(h))
	copy(hSorted, h)
	sort.Sort(hSorted)
	return hSorted
}

// String returns the unicode string representation of this Hand.
// It is always sorted and therefore suitable for comparison.
// Note: one tile requires up to 7 bytes in utf-8 encoding.
// See Repr() for a more efficient representation.
func (h Hand) String() string {
	var b bytes.Buffer
	for _, t := range h.sort() {
		b.WriteString(t.String())
	}
	return b.String()
}

// Repr returns a space-efficient encoding of this Hand.
// No sorting is performed before encoding. For an encoding that
// is suitable for comparison, use the sort.Sort() function on the
// Hand first.
func (h Hand) Repr() string {
	b := make([]byte, len(h))
	for i, t := range h {
		b[i] = t.Repr()
	}
	// don't return []byte directly, it is not suitable as a map key
	return string(b)
}

// Remove returns a copy of this Hand with the tile at index i removed.
func (h Hand) Remove(i int) Hand {
	if len(h) == 0 {
		return nil
	}

	if i < 0 || i > len(h)-1 {
		panic(fmt.Sprintf("remove out of bounds: %d len(h)=%d", i, len(h)))
	}

	hNew := make(Hand, len(h)-1)
	iNew := 0
	for j, t := range h {
		if j == i {
			continue
		}
		hNew[iNew] = t
		iNew++
	}
	return hNew
}

// ToCount converts this Hand to a Counter. The result is completely independent
// of this Hand (i.e. no aliasing).
func (h Hand) ToCount() Counter {
	c := make(map[Tile]int)
	for _, t := range h {
		c[t]++
	}
	hNew := make(Hand, len(h))
	copy(hNew, h)
	return Counter{c, len(h), hNew}
}

// Split splits a Hand into sub-Hands that each contain the tiles belonging to the same suit.
func (h Hand) Split(sorted bool) map[Suit]Hand {
	out := make(map[Suit]Hand)

	var hUsed Hand
	if sorted {
		hUsed = h.sort()
	} else {
		hUsed = h
	}

	for _, t := range hUsed {
		out[t.Suit] = append(out[t.Suit], t)
	}
	return out
}

// CanMeld returns true if all the tiles in this Hand may be used in melds.
func (h Hand) CanMeld() bool {
	for _, t := range h {
		if !t.CanMeld() {
			return false
		}
	}
	return true
}

// IsPeng returns true if this Hand contains only 3 identical melding tiles.
func (h Hand) IsPeng() bool {
	return h.IsPengAt(0)
}

// IsChi returns true if this Hand contains only 3 consecutively increasing melding tiles.
func (h Hand) IsChi() bool {
	return h.IsChiAt(0)
}

// IsPair returns true if this Hand contains only 2 identical melding tiles.
func (h Hand) IsPair() bool {
	return h.IsPairAt(0)
}

// IsPengAt returns true if this Hand contains 3 identical and consecutive melding tiles starting
// at the index i.
func (h Hand) IsPengAt(i int) bool {
	if i >= len(h)-2 || len(h) < 3 || !h.CanMeld() {
		return false
	}

	return h[i] == h[i+1] && h[i+1] == h[i+2]
}

// IsChiAt returns true is this Hand contains 3 increasing basic tiles starting
// at the index i. The hand must be sorted first.
func (h Hand) IsChiAt(i int) bool {
	if i >= len(h)-2 || len(h) < 3 {
		return false
	}
	// This requires linear search: the next tile in the set may not be the next tile in the hand
	// Consider the hand b1 b2 b2 b3, IsChiAt(0) should return true because [0]b1 [1/2]b2 [3]b3
	// forms the set.

	t1 := h[i]
	if !t1.IsBasic() || t1.Value > 7 {
		return false
	}

	t2 := Tile{t1.Suit, t1.Value+1}
	t3 := Tile{t1.Suit, t1.Value+2}

	var foundt2, foundt3 bool

	for _, t := range h[i+1:] {
		if t == t2 {
			foundt2 = true
		}
		if t == t3 {
			foundt3 = true
		}
		if foundt2 && foundt3 {
			return true
		}
	}

	return false
}

// IsPairAt returns true if this Hand contains 2 identical and consecutive melding tiles starting
// at the index i.
func (h Hand) IsPairAt(i int) bool {
	if i >= len(h)-1 || len(h) < 2 || !h.CanMeld() {
		return false
	}

	return h[i] == h[i+1]
}
