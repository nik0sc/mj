package handcheck

import (
	"fmt"
	"os"
	"sort"

	"github.com/nik0sc/mj"
)

// OptHandRLEChecker implements an optimal hand checker.
// The optimal result minimises the number of tiles not participating in a meld,
// then maximises the number of 3-tile melds.
// Note that it may be possible for a hand to have multiple optimal solutions,
// only one will be returned in that case.
// The zero value is safe to use immediately (without any optimisations).
//
// Under the hood, this uses mj.HandRLE to represent the free tiles at each subproblem.
type OptHandRLEChecker struct {
	// cache map[string]Result

	// optimisations

	// Split hand into sub hands (melds are restricted to one hand only)
	// This is unused for now.
	Split bool
	// use memoisation to avoid O(2^n) running time
	UseMemo bool
}

type ohrstate struct {
	res    Result
	free   mj.HandRLE
	shared *shared
}

// Check finds the optimal grouping for a hand.
func (c OptHandRLEChecker) Check(hand mj.Hand) Result {
	h := make(mj.Hand, len(hand))
	copy(h, hand)
	sort.Sort(h)

	shared := shared{}
	if c.UseMemo {
		shared.memo = make(map[string]string)
	}
	// This will be useful later for recreating the free tiles
	cnt := h.ToCount()
	hr, err := mj.NewHandRLE(cnt.Entries()...)
	if err != nil {
		panic("Counter and HandRLE don't agree on entries: " + err.Error())
	}
	s := ohrstate{Result{}, hr, &shared}

	r := s.step()
	shared.writeSummary(os.Stdout)

	r.sort()

	// Reconstruct the free tiles
	r.Free, err = r.free(cnt.Map())
	if err != nil {
		panic(err)
	}

	return r
}

func (s ohrstate) step() Result {
	s.shared.enterStep(os.Stdout, s.free)

	if s.free.Len() == 0 {
		return s.res
	}

	repr := s.free.Marshal()
	// use memoization: this problem has optimal substructure and
	// overlapping subproblems, making it a good use for DP
	if r, ok := s.shared.getMemo(repr); ok {
		return r
	}

	best := s.res
	s.free.ForEach(func(i int, e mj.CountEntry) {
		if nextFree, ok := s.free.TryPengAt(i); ok {
			// build the state that results from building a peng with this tile
			if traceSteps {
				fmt.Printf("peng: %s x%d\n", e.Tile, e.Count)
			}

			r := ohrstate{Result{
				Pengs: s.res.Pengs.Append(e.Tile),
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs,
			}, nextFree, s.shared}.step() // the recursion

			// If this state results in an improvement, keep it
			if r.Score() > best.Score() {
				best = r
			}
		}

		if nextFree, ok := s.free.TryPairAt(i); ok {
			if traceSteps {
				fmt.Printf("pair: %s x%d\n", e.Tile, e.Count)
			}

			r := ohrstate{Result{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs.Append(e.Tile),
			}, nextFree, s.shared}.step()

			if r.Score() > best.Score() {
				best = r
			}
		}

		if nextFree, ok := s.free.TryChiAt(i); ok {
			if traceSteps {
				fmt.Printf("chi: %s x%d\n", e.Tile, e.Count)
			}

			r := ohrstate{Result{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis.Append(e.Tile),
				Pairs: s.res.Pairs,
			}, nextFree, s.shared}.step()

			if r.Score() > best.Score() {
				best = r
			}
		}
	})

	s.shared.setMemo(repr, best)
	return best
}
