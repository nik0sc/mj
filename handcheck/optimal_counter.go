package handcheck

import (
	"fmt"
	"sort"

	"github.com/nik0sc/mj"
)

type OptCountChecker struct {
	// cache map[string]Result

	// optimisations

	// Split hand into sub hands (melds are restricted to one hand only)
	// This is unused for now.
	Split bool
	// use memoisation to avoid O(2^n) running time
	UseMemo bool
}

type ocstate struct {
	// Unlike OptChecker, the result always has a nil Free field, because
	// the ocstate.free field already has that information.
	res    Result
	free   mj.Counter
	shared *ocshared
}

type ocshared struct {
	memo      map[string]string
	stepCount int
	memoHits  int
}

// Check finds the optimal grouping for a hand.
func (c OptCountChecker) Check(hand mj.Hand) Result {
	h := make(mj.Hand, len(hand))
	copy(h, hand)
	sort.Sort(h)

	shared := ocshared{}
	if c.UseMemo {
		shared.memo = make(map[string]string)
	}
	cnt := h.ToCount()
	s := ocstate{Result{}, cnt, &shared}

	r := s.step()
	if writeMetrics {
		fmt.Printf("shared: len(memo)=%d steps=%d memohits=%d\n",
			len(s.shared.memo), s.shared.stepCount, s.shared.memoHits)
	}

	r.sort()

	// Reconstruct the free tiles
	cmap := cnt.Map()
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
		panic("cannot recreate Counter from map: " + err.Error())
	}
	r.Free = freecnt.ToHand(true)

	return r
}

func (s ocstate) step() Result {
	if writeMetrics {
		s.shared.stepCount++
	}
	if traceSteps {
		fmt.Printf("at %s\n", s.free.String())
	}

	if s.free.Len() == 0 {
		return s.res
	}

	repr := s.free.ToHand(true).Marshal()
	if s.shared.memo != nil {
		// use memoization: this problem has optimal substructure and
		// overlapping subproblems, making it a good use for DP
		if r, ok := s.shared.memo[repr]; ok {
			if writeMetrics {
				s.shared.memoHits++
			}
			return UnmarshalResult(r)
		}
	}

	best := s.res
	for _, e := range s.free.Entries() {
		if nextFree, ok := s.free.TryPeng(e.Tile); ok {
			// build the state that results from building a peng with this tile
			nextPengs := make([]mj.Tile, len(s.res.Pengs)+1)

			copy(nextPengs, s.res.Pengs)
			nextPengs[len(nextPengs)-1] = e.Tile

			if traceSteps {
				fmt.Printf("peng: %s x%d\n", e.Tile, e.Count)
			}

			r := ocstate{Result{
				Pengs: nextPengs,
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs,
			}, nextFree, s.shared}.step() // the recursion

			// If this state results in an improvement, keep it
			if r.score() > best.score() {
				best = r
			}
		}

		if nextFree, ok := s.free.TryPair(e.Tile); ok {
			nextPairs := make([]mj.Tile, len(s.res.Pairs)+1)

			copy(nextPairs, s.res.Pairs)
			nextPairs[len(nextPairs)-1] = e.Tile

			if traceSteps {
				fmt.Printf("pair: %s x%d\n", e.Tile, e.Count)
			}

			r := ocstate{Result{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis,
				Pairs: nextPairs,
			}, nextFree, s.shared}.step()

			if r.score() > best.score() {
				best = r
			}
		}

		if nextFree, ok := s.free.TryChi(e.Tile); ok {
			nextChis := make([]mj.Tile, len(s.res.Chis)+1)

			copy(nextChis, s.res.Chis)
			nextChis[len(nextChis)-1] = e.Tile

			if traceSteps {
				fmt.Printf("chi: %s x%d\n", e.Tile, e.Count)
			}

			r := ocstate{Result{
				Pengs: s.res.Pengs,
				Chis:  nextChis,
				Pairs: s.res.Pairs,
			}, nextFree, s.shared}.step()

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
		// sort the result first
		result := best.Copy()
		result.sort()

		s.shared.memo[repr] = result.Marshal()
	}

	return best
}
