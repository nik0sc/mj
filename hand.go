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
// See Marshal() for a more efficient representation.
func (h Hand) String() string {
	var b bytes.Buffer
	for _, t := range h.sort() {
		b.WriteString(t.String())
	}
	return b.String()
}

// Marshal returns a space-efficient encoding of this Hand.
// No sorting is performed before encoding. For an encoding that
// is suitable for comparison, use the sort.Sort() function on the
// Hand first.
func (h Hand) Marshal() string {
	b := make([]byte, len(h))
	for i, t := range h {
		b[i] = t.Marshal()
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

// Append returns a copy of this Hand with the tile appended to the end.
func (h Hand) Append(t Tile) Hand {
	hNew := make(Hand, len(h)+1)
	copy(hNew, h)
	hNew[len(h)] = t
	return hNew
}

// ToCount converts this Hand to a Counter. The result is completely independent
// of this Hand (i.e. no aliasing).
func (h Hand) ToCount() Counter {
	c := make(map[Tile]int)
	for _, t := range h {
		c[t]++
	}
	return Counter{c, len(h)}
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

// TryPengAt attempts to form a peng with the tile at the given index.
// If it succeeds, it returns (a new Hand with those tiles removed, true).
// Otherwise, it returns (nil, false). The hand should be sorted first.
//
// It is possible to return (nil, true) if i == 0, len(h) == 3 and h[0]
// == h[1] == h[2].
func (h Hand) TryPengAt(i int) (Hand, bool) {
	return h.tryMeldRunAt(i, 3)
}

// TryPairAt attempts to form a pair with the tile at the given index.
// If it succeeds, it returns (a new Hand with those tiles removed, true).
// Otherwise, it returns (nil, false). The hand should be sorted first.
//
// It is possible to return (nil, true) if i == 0, len(h) == 2 and h[0]
// == h[1].
func (h Hand) TryPairAt(i int) (Hand, bool) {
	return h.tryMeldRunAt(i, 2)
}

// tryMeldRunAt generalises TryPair/PengAt (and maybe TryGangAt in the future).
func (h Hand) tryMeldRunAt(i, n int) (Hand, bool) {
	if n < 2 {
		panic("tryRunAt: n < 2")
	}

	// if len(h)=4, i<=1 and n=3, we should pass
	// but if len(h)=4, i>=2 and n=3, we should fail
	if i > len(h)-n || len(h) < n {
		return nil, false
	}

	t := h[i]
	if !t.CanMeld() {
		return nil, false
	}

	// rely on sorted hand
	for j := 1; j < n; j++ {
		if t != h[i+j] {
			return nil, false
		}
	}

	if len(h) == n {
		// the else branch would result in a zero-length
		// slice, so just return nil quickly
		return nil, true
	} else {
		hNew := make(Hand, len(h))
		copy(hNew, h)
		hNew = append(hNew[:i], hNew[i+n:]...)

		return hNew, true
	}
}

// TryChiAt attempts to form a chi starting with the tile at index i.
// If it succeeds, it returns (a new Hand with the chi tiles removed, true).
// Otherwise, it returns (nil, false). The hand should be sorted first.
//
// It is possible to return (nil, true) if i == 0, len(h) == 3 and the
// 3 tiles in the hand form a chi by themselves.
func (h Hand) TryChiAt(i int) (Hand, bool) {
	if i >= len(h)-2 || len(h) < 3 {
		return nil, false
	}

	// easy case
	if len(h) == 3 {
		if i != 0 {
			panic("TryChiAt: buggy guard for easy case")
		}
		if !h[0].IsBasic() {
			return nil, false
		}
		tRepr := h[0].Marshal()
		if tRepr+1 == h[1].Marshal() && tRepr+2 == h[2].Marshal() {
			return nil, true
		} else {
			return nil, false
		}
	}

	t1 := h[i]
	if !t1.IsBasic() || t1.Value > 7 {
		return nil, false
	}

	t2 := Tile{t1.Suit, t1.Value + 1}
	t3 := Tile{t1.Suit, t1.Value + 2}

	i2, i3 := -1, -1
	// This requires linear search: the next tile in the set may not be the next tile in the hand
	// Consider the hand b1 b2 b2 b3, TryChiAt(0) should return true because [0]b1 [1/2]b2 [3]b3
	// forms the set.
	for j, t := range h[i+1:] {
		if t == t2 {
			// need to offset the current index since we are
			// iterating over a slice of the original
			i2 = i + j + 1
		} else if t == t3 {
			i3 = i + j + 1
			// assumes sorted hand!
			break
		}
	}

	if i2 >= 0 && i3 >= 0 {
		// make the copy now, allocate as late as possible
		// to reduce memory and gc pressure at the
		// expense of the happy (but infrequent) path
		iNew := 0
		hNew := make(Hand, len(h)-3)
		for j := range h {
			if j == i || j == i2 || j == i3 {
				continue
			}
			hNew[iNew] = h[j]
			iNew++
		}
		return hNew, true
	}

	return nil, false
}

func (h Hand) IsPair() bool {
	if len(h) != 2 {
		return false
	}
	if !h[0].CanMeld() {
		return false
	}
	return h[0] == h[1]
}

func (h Hand) IsPeng() bool {
	if len(h) != 3 {
		return false
	}
	if !h[0].CanMeld() {
		return false
	}
	return h[0] == h[1] && h[1] == h[2]
}

func (h Hand) IsChi() bool {
	if len(h) != 3 {
		return false
	}
	if !h[0].IsBasic() {
		return false
	}
	tRepr := h[0].Marshal()
	if tRepr+1 == h[1].Marshal() && tRepr+2 == h[2].Marshal() {
		return true
	} else {
		return false
	}
}

func UnmarshalHand(s string) Hand {
	repr := []byte(s)
	h := make(Hand, len(repr))

	for i, b := range repr {
		h[i] = UnmarshalTile(b)
	}

	return h
}
