package handcheck

import (
	"fmt"
	"sort"

	"mj"
)

// OptChecker implements an optimal hand checker.
// The optimal result minimises the number of tiles not participating in a meld,
// then maximises the number of 3-tile melds.
// Note that it may be possible for a hand to have multiple optimal solutions,
// only one will be returned in that case.
//
// The zero value is safe to use immediately (without any optimisations).
type OptChecker struct {
	// If OptChecker is reused for multiple hands (perhaps in a mahjong-playing AI agent),
	// we can cache the results.
	// TODO: Need some kind of cache eviction policy in that case.
	cache map[string]Result

	// optimisations

	// Split hands by suit into sub-hands. Since melds are restricted to one suit only,
	// this should reduce the search space without too much effort.
	Split bool
	// UseMemo enables memoisation of repeated subproblems when solving a hand.
	// This should really always be on.
	UseMemo bool
}

type ostate struct {
	res    Result
	shared *oshared
}

type oshared struct {
	memo      map[string]Result
	stepCount int
	memoHits  int
}

// Check finds the optimal grouping for a hand.
func (c OptChecker) Check(hand mj.Hand) Result {
	h := make(mj.Hand, len(hand))
	copy(h, hand)
	// very important, when we search for melds we depend on sorted order
	sort.Sort(h)

	// did we solve this hand before?
	hrepr := h.Repr()
	if c.cache == nil {
		c.cache = make(map[string]Result)
	}
	if r, ok := c.cache[hrepr]; ok {
		return r
	}

	var r Result
	if c.Split {
		// no need to sort again
		hsplit := h.Split(false)
		rs := make([]Result, 0, len(hsplit))
		for _, hs := range hsplit {
			// improvement: could start in goroutines
			rs = append(rs, c.start(hs))
		}

		for _, rsub := range rs {
			r.Pengs = append(r.Pengs, rsub.Pengs...)
			r.Chis = append(r.Chis, rsub.Chis...)
			r.Pairs = append(r.Pairs, rsub.Pairs...)
			r.Free = append(r.Free, rsub.Free...)
		}
	} else {
		r = c.start(h)
	}

	c.cache[hrepr] = r
	return r
}

func (c OptChecker) start(h mj.Hand) Result {
	shared := oshared{}
	if c.UseMemo {
		shared.memo = make(map[string]Result)
	}
	// at first, the entire hand is free
	s := ostate{Result{Free: h}, &shared}

	r := s.step()
	if writeMetrics {
		fmt.Printf("shared: len(memo)=%d steps=%d memohits=%d\n",
			len(s.shared.memo), s.shared.stepCount, s.shared.memoHits)
	}

	return r
}

func (s ostate) step() Result {
	if writeMetrics {
		s.shared.stepCount++
	}
	if traceSteps {
		fmt.Printf("at %s\n", s.res.Free.String())
	}
	// invariant: s.res.Free is always in sorted order

	numFree := len(s.res.Free)

	// base case
	if numFree == 0 {
		return s.res
	}

	repr := s.res.Free.Repr()
	if s.shared.memo != nil {
		// use memoization: this problem has optimal substructure and
		// overlapping subproblems, making it a good use for DP
		if r, ok := s.shared.memo[repr]; ok {
			if writeMetrics {
				s.shared.memoHits++
			}
			return r
		}
	}

	// The best result so far is the one from our parent
	best := s.res
	for i, t := range s.res.Free {
		// try and build a set with this tile
		// the hand is always kept in sorted order, this vastly simplifies building
		// peng is checked first, then chi, then pair, this is intentional
		if s.res.Free.IsPengAt(i) {
			// build the state that results from building a peng with this tile
			nextPengs := make([]mj.Tile, len(s.res.Pengs)+1)
			nextFree := make(mj.Hand, numFree)

			copy(nextPengs, s.res.Pengs)
			nextPengs[len(nextPengs)-1] = t

			copy(nextFree, s.res.Free)
			nextFree = append(nextFree[:i], nextFree[i+3:]...)

			r := ostate{Result{
				Pengs: nextPengs,
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs,
				Free:  nextFree,
			}, s.shared}.step() // the recursion

			// If this state results in an improvement, keep it
			if r.score() > best.score() {
				best = r
			}
		}

		if s.res.Free.IsChiAt(i) {
			nextChis := make([]mj.Tile, len(s.res.Chis)+1)
			nextFree := make(mj.Hand, numFree-3)

			copy(nextChis, s.res.Chis)
			nextChis[len(nextChis)-1] = t

			// This is a bit more complicated, since the tiles in a chi are not necessarily
			// consecutive in the hand
			var removed1, removed2, removed3 bool
			nextFreeIdx := 0
			// For faster comparison and manipulation
			tRepr := t.Repr()

			for _, tFree := range s.res.Free {
				freeRepr := tFree.Repr()
				if tRepr == freeRepr && !removed1 {
					removed1 = true
				} else if tRepr+1 == freeRepr && !removed2 {
					// Adding to tRepr: the lowest bits of repr are the value of the tile.
					// So we can just increment the repr to increment the value.
					removed2 = true
				} else if tRepr+2 == freeRepr && !removed3 {
					removed3 = true
				} else {
					// only in this case, copy the tile in
					nextFree[nextFreeIdx] = tFree
					nextFreeIdx++
				}
			}

			r := ostate{Result{
				Pengs: s.res.Pengs,
				Chis:  nextChis,
				Pairs: s.res.Pairs,
				Free:  nextFree,
			}, s.shared}.step()

			if r.score() > best.score() {
				best = r
			}
		}

		if s.res.Free.IsPairAt(i) {
			nextPairs := make([]mj.Tile, len(s.res.Pairs)+1)
			nextFree := make(mj.Hand, numFree)

			copy(nextPairs, s.res.Pairs)
			nextPairs[len(nextPairs)-1] = t

			copy(nextFree, s.res.Free)
			nextFree = append(nextFree[:i], nextFree[i+2:]...)

			r := ostate{Result{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis,
				Pairs: nextPairs,
				Free:  nextFree,
			}, s.shared}.step()

			if r.score() > best.score() {
				best = r
			}
		}
	}

	if s.shared.memo != nil {
		if rOld, ok := s.shared.memo[repr]; ok {
			// memo should not be updated like this! because the memo result should already be optimal for
			// the currently free tiles
			fmt.Printf("updating memo? repr=%x rOld=%+v best=%+v", repr, rOld, best)
		}
		s.shared.memo[repr] = best
	}

	return best
}
