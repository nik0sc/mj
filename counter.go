package mj

import (
	"bytes"
	"sort"
)

// Counter is a counter of tiles. It is an alternative unordered representation of Hand.
// It offers constant-time lookup of tile counts. However creating and modifying it is
// more expensive. The methods of Counter are guaranteed not to mutate the struct or cause
// memory aliasing.
type Counter struct {
	m map[Tile]int
	n int
	h Hand
}

// CountEntry is a pair of a Tile and its count.
type CountEntry struct {
	Tile  Tile
	Count int16
	// This fits nicely in 4 bytes and should not require padding in most architectures.
}

// Valid returns true if the Counter is valid and all the tiles in the Counter are valid.
func (c Counter) Valid() bool {
	if c.m == nil || c.h == nil {
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

// Entries returns all tile-count pairs in the Counter.
//
// Warning: CountEntry.Count is int16, not int. This means the maximum count
// for a tile is 32767. If you really need to count more tiles than that,
// use Get() instead.
func (c Counter) Entries() []CountEntry {
	es := make([]CountEntry, 0, len(c.m))
	for t, cnt := range c.m {
		es = append(es, CountEntry{Tile: t, Count: int16(cnt)})
	}
	return es
}

// ToHand converts this Counter to a Hand. The result is not guaranteed to be sorted.
func (c Counter) ToHand() Hand {
	h := make(Hand, c.n)
	copy(h, c.h)
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

	remove := 0
	for i, tNew := range cNew.h {
		if tNew == t {
			remove = i
			break
		}
	}
	cNew.h = append(cNew.h[:remove], cNew.h[remove+1:]...)

	return cNew
}

// Copy returns a deep copy of this Counter.
func (c Counter) Copy() Counter {
	m := make(map[Tile]int)
	for t, n := range c.m {
		m[t] = n
	}

	h := make(Hand, len(c.h))
	copy(h, c.h)

	return Counter{m, c.n, h}
}

// TryPeng TODO
func (c Counter) TryPeng(t Tile) Counter {
	if c.m[t] < 3 {
		return Counter{}
	}

	cNew := c.Copy()
	cNew.m[t] -= 3
	cNew.n -= 3

	sort.Sort(cNew.h)
	remove := -1
	for i, ht := range cNew.h {
		if ht == t {
			remove = i
			break
		}
	}
	if remove == -1 {
		panic("unreachable")
	}
	cNew.h = append(cNew.h[:remove], cNew.h[remove+3:]...)

	return cNew
}

// TryChi TODO
func (c Counter) TryChi(t Tile) Counter {
	if !t.IsBasic() {
		return Counter{}
	}
	t2 := Tile{Suit: t.Suit, Value: t.Value+1}
	t3 := Tile{Suit: t.Suit, Value: t.Value+2}
	if !t2.Valid() || !t3.Valid() {
		return Counter{}
	}

	if c.m[t] <= 0 || c.m[t2] <= 0 || c.m[t3] <= 0 {
		return Counter{}
	}

	// don't bother with Copy
	nNew := c.n - 3
	mNew := make(map[Tile]int)
	for tt, n := range c.m {
		if (tt == t || tt == t2 || tt == t3) && n > 1 {
			mNew[tt] = n-1
		} else {
			mNew[tt] = n
		}
	}

	hNew := make(Hand, nNew)
	hIndex := 0
	for tt, n := range mNew {
		for i := 0; i < n; i++ {
			hNew[hIndex] = tt
			hIndex++
		}
	}
	if !hNew.Valid() {
		panic("hNew not created correctly")
	}

	return Counter{m: mNew, n: nNew, h: hNew}
}

// TryPair TODO
func (c Counter) TryPair(t Tile) Counter {
	if c.m[t] < 2 {
		return Counter{}
	}

	cNew := c.Copy()
	cNew.m[t] -= 2
	cNew.n -= 2

	sort.Sort(cNew.h)
	remove := -1
	for i, ht := range cNew.h {
		if ht == t {
			remove = i
			break
		}
	}
	if remove == -1 {
		panic("unreachable")
	}
	cNew.h = append(cNew.h[:remove], cNew.h[remove+2:]...)

	return cNew
}

// _iter is a cool idea, because it lets us for-range over the Counter,
// but it starts a goroutine and allocates a channel. Just use Entries() instead.
func (c Counter) _iter() (iter chan CountEntry, stop chan struct{}) {
	iter = make(chan CountEntry)
	stop = make(chan struct{})

	go func() {
	iterFor:
		for t, c := range c.m {
			e := CountEntry{
				Tile:  t,
				Count: int16(c),
			}
			select {
			case iter <- e:
				continue
			case <-stop:
				break iterFor
			}
		}
		close(iter)
	}()

	return
}
