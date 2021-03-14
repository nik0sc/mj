package mj

import (
	"bytes"
	"fmt"
	"sort"
)

// Counter is a counter of tiles. It is an alternative unordered representation of Hand.
// It offers constant-time lookup of tile counts. However creating and modifying it is
// more expensive. The methods of Counter are guaranteed not to mutate the struct or cause
// memory aliasing.
//
// Counter does not have a Marshal method, since the underlying representation is unordered.
// Use Counter.ToHand(true).Marshal().
//
// You can bypass the methods of Counter by converting to a map[Tile]int with Counter.Map(),
// modifying the map, and passing it to NewCounter(), which will also verify your map.
type Counter struct {
	m map[Tile]int
	n int
	// Unlike Hand, this is an opaque type
}

// CountEntry is a pair of a Tile and its count.
type CountEntry struct {
	Tile  Tile
	Count int16
	// This fits nicely in 4 bytes and should not require padding in most architectures.
}

// NewCounter creates a new Counter from a map of tiles to their counts.
func NewCounter(m map[Tile]int) (Counter, error) {
	c := Counter{m: make(map[Tile]int, len(m))}

	for t, n := range m {
		if !t.Valid() {
			return Counter{}, fmt.Errorf("tile %+v is invalid", t)
		}
		if n < 0 {
			return Counter{}, fmt.Errorf("invalid count for tile %+v: %d", t, n)
		}
		if n == 0 {
			continue
		}
		c.m[t] += n
		c.n += n
	}

	return c, nil
}

// Valid returns true if the Counter is valid and all the tiles in the Counter are valid.
// The zero Counter causes Valid to return false.
func (c Counter) Valid() bool {
	if c.m == nil {
		return false
	}

	for t := range c.m {
		if !t.Valid() {
			return false
		}
	}
	return true
}

// Len returns the number of tiles in the Counter.
func (c Counter) Len() int {
	return c.n
}

// Get returns the count of a tile.
func (c Counter) Get(t Tile) int {
	return c.m[t]
}

// Map returns a map of Tiles to their counts in the Counter. It is
// the inverse of NewCounter(). This map is a copy, so changes to
// the Counter won't be reflected in this map, or vice versa.
func (c Counter) Map() map[Tile]int {
	if c.m == nil {
		return nil
	}

	m := make(map[Tile]int, len(c.m))
	for t, n := range c.m {
		m[t] = n
	}
	return m
}

// Entries returns all tile-count pairs in the Counter.
//
// Warning: CountEntry.Count is int16, not int. This means the maximum count
// for a tile is 32767. If you really need to count more tiles than that,
// use Get() instead.
func (c Counter) Entries() []CountEntry {
	// If we have Map, what's the point of Entries?
	es := make([]CountEntry, 0, len(c.m))
	for t, cnt := range c.m {
		es = append(es, CountEntry{Tile: t, Count: int16(cnt)})
	}
	return es
}

// ToHand converts this Counter to a Hand. If sorted is true, the hand is
// guaranteed to be in sorted order.
func (c Counter) ToHand(sorted bool) Hand {
	h := make(Hand, c.n)
	i := 0
	for t, n := range c.m {
		for j := 0; j < n; j++ {
			h[i] = t
			i++
		}
	}

	if sorted {
		sort.Sort(h)
	}

	return h
}

// String returns the human-readable representation of this Counter.
// The caveats of Hand.String() apply here, with one additional: the
// order of tiles is undefined, but the same tiles will appear together.
func (c Counter) String() string {
	var b bytes.Buffer
	for t, n := range c.m {
		for i := 0; i < n; i++ {
			b.WriteString(t.String())
		}
	}
	return b.String()
}

// Remove deep-copies this Counter and removes a tile from it.
// It panics if the tile isn't in the counter.
func (c Counter) Remove(t Tile) Counter {
	if c.m[t] <= 0 {
		panic("no tiles to remove: " + t.String())
	}

	cNew := c.Copy()
	cNew.m[t]--
	cNew.n--

	return cNew
}

// Copy returns a deep copy of this Counter.
func (c Counter) Copy() Counter {
	m := make(map[Tile]int)
	for t, n := range c.m {
		m[t] = n
	}

	return Counter{m, c.n}
}

// TryPeng attempts to form a peng with the given tile. If it succeeds, it
// returns (a new Counter with the peng removed, true). Otherwise, it
// returns (a zero Counter, false).
//
// It is possible to return (zero Counter, true) if the 3 tiles to be removed
// are the only tiles in the original Counter.
func (c Counter) TryPeng(t Tile) (Counter, bool) {
	return c.tryMeldRun(t, 3)
}

// TryChi attempts to form a chi with the given tile as the first in the set.
// If it succeeds, it returns a new Counter with one of each of the given tile,
// the next tile, and the one after that, all removed. Otherwise, it returns
// (a zero Counter, false).
//
// For example: (not the real syntax)
//   Counter{B1:1 B2:2 B3:1 B4:1}.TryChi(B1) -> Counter{B2:1 B4:1}
// Note that one B1, one B2 and one B3 were removed.
//
// It is possible to return (zero Counter, true) if the 3 tiles to be removed
// are the only tiles in the original Counter.
func (c Counter) TryChi(t Tile) (Counter, bool) {
	if !t.IsBasic() {
		return Counter{}, false
	}
	t2 := Tile{Suit: t.Suit, Value: t.Value + 1}
	t3 := Tile{Suit: t.Suit, Value: t.Value + 2}
	if !t2.Valid() || !t3.Valid() {
		return Counter{}, false
	}

	if c.m[t] <= 0 || c.m[t2] <= 0 || c.m[t3] <= 0 {
		return Counter{}, false
	}

	if c.n == 3 && c.m[t] == 1 && c.m[t2] == 1 && c.m[t3] == 1 {
		return Counter{}, true
	}

	// don't bother with Copy
	nNew := c.n - 3
	mNew := make(map[Tile]int)
	for tt, n := range c.m {
		if tt == t || tt == t2 || tt == t3 {
			if n > 1 {
				mNew[tt] = n - 1
			} // else, don't copy into the new map
		} else {
			mNew[tt] = n
		}
	}

	return Counter{mNew, nNew}, true
}

// TryPeng attempts to form a pair with the given tile. If it succeeds, it
// returns (a new Counter with the pair removed, true). Otherwise, it
// returns (a zero Counter, false).
//
// It is possible to return (zero Counter, true) if the 2 tiles to be removed
// are the only tiles in the original Counter.
func (c Counter) TryPair(t Tile) (Counter, bool) {
	return c.tryMeldRun(t, 2)
}

// tryMeldRun generalises TryPeng and TryPair.
func (c Counter) tryMeldRun(t Tile, n int) (Counter, bool) {
	if c.m[t] < n {
		return Counter{}, false
	} else if c.n == n && c.m[t] == n {
		return Counter{}, true
	}

	cNew := c.Copy()
	cNew.m[t] -= n
	cNew.n -= n

	return cNew, true
}
