package mj

import (
	"fmt"
	"sort"
	"strings"
)

// HandRLE is a run-length encoded version of Hand. It combines the best of Hand and Counter:
// like Hand, it is compact, contiguous in memory and preserves order, but like Counter,
// it stores tiles and their counts.
//
// Unfortunately, the public API is also a weird mix of Hand and Counter.
// Tile lookup now uses binary search, since the tile-count pairs are stored in order.
type HandRLE struct {
	es []CountEntry
	n  int
}

// NewHandRLE creates a new HandRLE from one or more CountEntry values.
func NewHandRLE(entries ...CountEntry) (HandRLE, error) {
	if len(entries) == 0 {
		return HandRLE{}, nil
	}

	h := HandRLE{es: make([]CountEntry, len(entries))}
	copy(h.es, entries)

	var err error
	h.n, err = h.valid(false)
	if err != nil {
		return HandRLE{}, err
	}
	h.sort()

	return h, nil
}

func (h HandRLE) valid(checkN bool) (int, error) {
	n := 0
	seen := make(map[Tile]bool)
	for _, e := range h.es {
		if !e.Tile.Valid() {
			return 0, fmt.Errorf("invalid tile: %+v", e.Tile)
		}
		if e.Count <= 0 {
			return 0, fmt.Errorf("invalid count: %d", e.Count)
		}
		dup := seen[e.Tile]
		if dup {
			return 0, fmt.Errorf("duplicated tile in entry: %+v", e.Tile)
		}
		seen[e.Tile] = true
		n += int(e.Count)
	}
	if checkN && h.n != n {
		return 0, fmt.Errorf("wrong total: n=%d, counted=%d", h.n, n)
	}

	return n, nil
}

// Valid returns true if the HandRLE is valid. A valid hand is not necessarily in sorted order,
// but a valid hand cannot contain an entry meeting any of these conditions:
//  - count is 0 or less
//  - tile is invalid
//  - tile previously occurred in the HandRLE
//
func (h HandRLE) Valid() bool {
	_, err := h.valid(true)
	return err == nil
}

// Copy deep-copies this HandRLE.
func (h HandRLE) Copy() HandRLE {
	esNew := make([]CountEntry, len(h.es))
	copy(esNew, h.es)
	return HandRLE{es: esNew, n: h.n}
}

// sort sorts HandRLE.es in-place.
func (h HandRLE) sort() {
	sort.Slice(h.es, func(i, j int) bool {
		return h.es[i].Tile.Less(h.es[j].Tile)
	})
}

// Len returns the number of tiles in the HandRLE.
func (h HandRLE) Len() int {
	return h.n
}

// Get returns the count of a tile.
func (h HandRLE) Get(t Tile) int {
	sl := h.es

	// Obviously, this only works if HandRLE is always sorted (which it is)
	for len(sl) > 0 {
		i := len(sl) / 2
		if sl[i].Tile == t {
			return int(sl[i].Count)
		} else if t.Less(sl[i].Tile) {
			sl = sl[:i]
		} else {
			sl = sl[i+1:]
		}
	}

	return 0
}

// Entries returns all tile-count pairs in the HandRLE.
func (h HandRLE) Entries() []CountEntry {
	esNew := make([]CountEntry, len(h.es))
	copy(esNew, h.es)
	return esNew
}

// ForEach allows iteration over the tile-count pairs without the extra copying of Entries.
// Using HandRLE.ForEach() instead of ranging over the result of HandRLE.Entries() can save
// a lot of time and memory. The passed-in function should accept an index and a CountEntry
// and return whether or not to continue the iteration.
func (h HandRLE) ForEach(f func(int, CountEntry) bool) {
	for i, e := range h.es {
		if !f(i, e) {
			break
		}
	}
}

// String returns the unicode string representation of this HandRLE.
// It is always sorted and therefore suitable for comparison.
// Note: one tile requires up to 7 bytes in utf-8 encoding.
// See Marshal() for a more efficient representation.
func (h HandRLE) String() string {
	var sb strings.Builder
	predicted := 4
	if uniUseVS16 {
		predicted += 3
	}
	sb.Grow(predicted * h.n)

	for _, e := range h.es {
		s := e.Tile.String()
		for i := 0; i < int(e.Count); i++ {
			sb.WriteString(s)
		}
	}

	return sb.String()
}

// Marshal returns a space-efficient encoding of this HandRLE.
// It is suitable for comparison, because the output is always sorted in the same order.
func (h HandRLE) Marshal() string {
	var sb strings.Builder
	sb.Grow(len(h.es) * 2)
	for _, e := range h.es {
		sb.WriteByte(e.Tile.Marshal())
		if e.Count <= 0x7f {
			sb.WriteByte(byte(e.Count))
		} else {
			panic("unimplemented: e.Count > 0x7f")
		}
	}
	return sb.String()
}

// TryPengAt attempts to form a peng with the tile at the given index.
// If it succeeds, it returns (a new HandRLE with those tiles removed, true).
// Otherwise, it returns (zero HandRLE, false).
//
// It is possible to return (zero HandRLE, true) if the 3 tiles to be removed
// are the only tiles in the original HandRLE.
func (h HandRLE) TryPengAt(i int) (HandRLE, bool) {
	return h.tryMeldRunAt(i, 3)
}

// TryPairAt attempts to form a pair with the tile at the given index.
// If it succeeds, it returns (a new HandRLE with the pair removed, true).
// Otherwise, it returns (zero HandRLE, false).
//
// It is possible to return (zero HandRLE, true) if the 2 tiles to be removed
// are the only tiles in the original HandRLE.
func (h HandRLE) TryPairAt(i int) (HandRLE, bool) {
	return h.tryMeldRunAt(i, 2)
}

func (h HandRLE) tryMeldRunAt(i, n int) (HandRLE, bool) {
	if n < 2 {
		panic("tryMeldRunAt: n < 2")
	}

	if len(h.es) == 1 && i == 0 && h.n == n && h.es[0].Tile.CanMeld() {
		// or equivalently, h.es[0].Count == n
		return HandRLE{}, true
	}

	if i >= len(h.es) {
		return HandRLE{}, false
	}

	if !h.es[i].Tile.CanMeld() {
		return HandRLE{}, false
	}

	if int(h.es[i].Count) < n {
		return HandRLE{}, false
	}

	hrNew := h.Copy()
	hrNew.n -= n

	if int(h.es[i].Count) > n {
		hrNew.es[i].Count -= int16(n)
	} else {
		hrNew.es = append(hrNew.es[:i], hrNew.es[i+1:]...)
	}
	return hrNew, true
}

// TryChiAt attempts to form a chi starting with the tile at index i.
// If it succeeds, it returns (a new HandRLE with the chi tiles removed, true).
// Otherwise, it returns (zero HandRLE, false).
//
// For example: (not the real syntax)
//   HandRLE{B1:1 B2:2 B3:1 B4:1}.TryChi(B1) -> HandRLE{B2:1 B4:1}
// Note that one B1, one B2 and one B3 were removed.
//
// It is possible to return (zero HandRLE, true) if the 3 tiles to be removed
// are the only tiles in the original HandRLE.
func (h HandRLE) TryChiAt(i int) (HandRLE, bool) {
	if i >= len(h.es)-2 || len(h.es) < 3 {
		return HandRLE{}, false
	}

	// this might be buggy
	if len(h.es) == 3 && h.n == 3 {
		if i != 0 {
			panic("TryChiAt: buggy guard for easy case")
		}
		if h.es[0].Count != 1 || h.es[1].Count != 1 || h.es[2].Count != 1 {
			panic("TryChiAt: wrong (impossible?) tile counts for easy case")
		}

		if !h.es[0].Tile.IsBasic() {
			return HandRLE{}, false
		}
		tRepr := h.es[0].Tile.Marshal()
		if tRepr+1 == h.es[1].Tile.Marshal() && tRepr+2 == h.es[2].Tile.Marshal() {
			return HandRLE{}, true
		} else {
			return HandRLE{}, false
		}
	}

	t1 := h.es[i].Tile
	if !t1.IsBasic() || t1.Value > 7 || h.es[i].Count < 1 {
		return HandRLE{}, false
	}

	t2 := Tile{t1.Suit, t1.Value + 1}
	if h.es[i+1].Tile != t2 {
		return HandRLE{}, false
	}

	t3 := Tile{t1.Suit, t1.Value + 2}
	if h.es[i+2].Tile != t3 {
		return HandRLE{}, false
	}

	// not the most efficient
	esNew := make([]CountEntry, i, len(h.es))
	copy(esNew, h.es[:i])

	copied := 0
	for _, e := range h.es[i : i+3] {
		if e.Count > 1 {
			// e is a copy, this is ok
			e.Count--
			esNew = append(esNew, e)
			copied++
		}
	}

	esNew = append(esNew, h.es[i+3:]...)
	return HandRLE{es: esNew, n: h.n - 3}, true
}

func UnmarshalHandRLE(s string) HandRLE {
	b := []byte(s)
	if len(b)%2 != 0 {
		panic("odd number of bytes in s")
	}
	hr := HandRLE{es: make([]CountEntry, len(b)/2)}

	for i := 0; i < len(b)/2; i++ {
		j := i * 2
		cnt := b[j+1]
		if cnt > 0x7f {
			panic("unimplemented: marshal tile count > 0x7f")
		}
		hr.es[i] = CountEntry{Tile: UnmarshalTile(b[j]), Count: int16(cnt)}
		hr.n += int(b[j+1])
	}

	return hr
}
