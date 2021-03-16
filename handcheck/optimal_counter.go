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
	shared *shared
}

// Check finds the optimal grouping for a hand.
func (c OptCountChecker) Check(hand mj.Hand) Result {
	h := make(mj.Hand, len(hand))
	copy(h, hand)
	sort.Sort(h)

	shared := shared{}
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
	// use memoization: this problem has optimal substructure and
	// overlapping subproblems, making it a good use for DP
	if r, ok := s.shared.getMemo(repr); ok {
		return r
	}

	best := s.res
	for _, e := range s.free.Entries() {
		if nextFree, ok := s.free.TryPeng(e.Tile); ok {
			// build the state that results from building a peng with this tile
			if traceSteps {
				fmt.Printf("peng: %s x%d\n", e.Tile, e.Count)
			}

			r := ocstate{Result{
				Pengs: s.res.Pengs.Append(e.Tile),
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs,
			}, nextFree, s.shared}.step() // the recursion

			// If this state results in an improvement, keep it
			if r.score() > best.score() {
				best = r
			}
		}

		if nextFree, ok := s.free.TryPair(e.Tile); ok {
			if traceSteps {
				fmt.Printf("pair: %s x%d\n", e.Tile, e.Count)
			}

			r := ocstate{Result{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs.Append(e.Tile),
			}, nextFree, s.shared}.step()

			if r.score() > best.score() {
				best = r
			}
		}

		if nextFree, ok := s.free.TryChi(e.Tile); ok {
			if traceSteps {
				fmt.Printf("chi: %s x%d\n", e.Tile, e.Count)
			}

			r := ocstate{Result{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis.Append(e.Tile),
				Pairs: s.res.Pairs,
			}, nextFree, s.shared}.step()

			if r.score() > best.score() {
				best = r
			}
		}
	}

	s.shared.setMemo(repr, best)
	return best
}
