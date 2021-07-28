package handcheck

import (
	"os"
	"sort"

	lru "github.com/hashicorp/golang-lru"
	"github.com/nik0sc/mj"
	"github.com/nik0sc/mj/internal/cache"
)

// OptChecker implements an optimal hand checker.
// The optimal result minimises the number of tiles not participating in a meld,
// then maximises the number of 3-tile melds.
// Note that it may be possible for a hand to have multiple optimal solutions,
// only one will be returned in that case.
// The zero value is safe to use immediately (without any optimisations).
//
// Under the hood, this uses mj.Hand to represent the free tiles at each subproblem.
// While this has a higher branching factor, it actually has lower runtime and memory
// usage in average cases.
type OptChecker struct {
	// If OptChecker is reused for multiple hands (perhaps in a mahjong-playing AI agent),
	// we can cache the results.
	cache cache.Cache
	// optimisations

	// Split hands by suit into sub-hands. Since melds are restricted to one suit only,
	// this should reduce the search space without too much effort.
	Split bool
	// UseMemo enables memoisation of repeated subproblems when solving a hand.
	// This should really always be on.
	UseMemo bool
}

// NewOptCheckerWithCache creates an OptChecker with cache and all optimisations enabled.
// size is the maximum number of solutions to be cached.
func NewOptCheckerWithCache(size int) (OptChecker, error) {
	c := OptChecker{
		Split:   true,
		UseMemo: true,
	}
	var err error

	c.cache, err = lru.New(size)
	if err != nil {
		return OptChecker{}, err
	}

	return c, nil
}

type ostate struct {
	res    mj.Group
	shared *shared
}

// Check finds the optimal grouping for a hand.
func (c OptChecker) Check(hand mj.Hand) mj.Group {
	h := make(mj.Hand, len(hand))
	copy(h, hand)
	// very important, when we search for melds we depend on sorted order
	sort.Sort(h)

	var hrepr string

	if c.cache != nil {
		// did we solve this hand before?
		hrepr = h.Marshal()
		if v, ok := c.cache.Get(hrepr); ok {
			return v.(*mj.Group).Copy(false)
		}
	}

	var r mj.Group
	if c.Split {
		// no need to sort again
		hsplit := h.Split(false)
		rs := make([]mj.Group, 0, len(hsplit))
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

	if c.cache != nil {
		rCopy := r.Copy(true)
		c.cache.Add(hrepr, &rCopy)
	}
	return r
}

func (c OptChecker) start(h mj.Hand) mj.Group {
	shr := shared{}
	if c.UseMemo {
		shr.memo = make(map[string]string)
	}
	// at first, the entire hand is free
	s := ostate{mj.Group{Free: h}, &shr}

	r := s.step()
	shr.writeSummary(os.Stdout)

	return r
}

func (s ostate) step() mj.Group {
	s.shared.enterStep(os.Stdout, s.res.Free)
	// invariant: s.res.Free is always in sorted order

	// base case
	if len(s.res.Free) == 0 {
		return s.res
	}

	repr := s.res.Free.Marshal()
	// use memoization: this problem has optimal substructure and
	// overlapping subproblems, making it a good use for DP
	if r, ok := s.shared.getMemo(repr); ok {
		return r
	}

	// The best result so far is the one from our parent
	best := s.res
	for i, t := range s.res.Free {
		// try and build a set with this tile
		// the hand is always kept in sorted order, this vastly simplifies building
		if nextFree, ok := s.res.Free.TryPengAt(i); ok {
			// build the state that results from building a peng with this tile
			r := ostate{mj.Group{
				Pengs: s.res.Pengs.Append(t),
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs,
				Free:  nextFree,
			}, s.shared}.step() // the recursion

			// If this state results in an improvement, keep it
			if r.Score() > best.Score() {
				best = r
			}
		}

		// A possible optimisation: Try pair first, and only if it succeeds, try peng
		// Tried it, causes test "all c" to fail on the fast but not on the slow version
		if nextFree, ok := s.res.Free.TryPairAt(i); ok {
			r := ostate{mj.Group{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs.Append(t),
				Free:  nextFree,
			}, s.shared}.step()

			if r.Score() > best.Score() {
				best = r
			}
		}

		if nextFree, ok := s.res.Free.TryChiAt(i); ok {
			r := ostate{mj.Group{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis.Append(t),
				Pairs: s.res.Pairs,
				Free:  nextFree,
			}, s.shared}.step()

			if r.Score() > best.Score() {
				best = r
			}
		}
	}

	s.shared.setMemo(repr, best)
	return best
}
