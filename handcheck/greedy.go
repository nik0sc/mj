package handcheck

import (
	"fmt"
	"sort"

	"github.com/nik0sc/mj"
)

// GreedyChecker greedily builds 3-tile melds and then returns a solution if it can
// build a pair with the last 2 tiles in the hand. While it is not optimal, it can
// be much faster than the Opt* checkers with certain hands.
//
// You probably don't want to use this checker for most cases.
type GreedyChecker struct {
	// Split=true breaks the guarantee that if we return ok=false there
	// is no possible winning interpretation of the hand
	Split    bool
	FailFast bool
}

type gstate struct {
	res Result

	h     mj.Hand
	build mj.Hand

	shared *gshared
}

type gshared struct {
	stepCount int
}

func (c GreedyChecker) Check(hand mj.Hand) Result {
	h := make(mj.Hand, len(hand))
	copy(h, hand)
	sort.Sort(h)

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

	return r
}

func (c GreedyChecker) start(h mj.Hand) Result {
	_ = c
	s := gstate{h: h, shared: &gshared{}}

	r, ok := s.step()
	if writeMetrics {
		fmt.Printf("shared: ok=%t, steps=%d\n", ok, s.shared.stepCount)
	}

	if ok {
		return r
	} else {
		return Result{Free: h}
	}
}

func (s gstate) step() (Result, bool) {
	if writeMetrics {
		s.shared.stepCount++
	}
	if traceSteps {
		fmt.Printf("at %s\n", s.h.String())
	}

	if len(s.h) < 2 {
		panic("short hand: " + s.h.String())
	}

	if len(s.h) == 2 {
		if s.h.IsPair() {
			// a winner!
			r := Result{
				Pengs: s.res.Pengs,
				Chis:  s.res.Chis,
				Pairs: s.res.Pairs.Append(s.h[0]),
				Free:  nil,
			}
			return r, true
		} else {
			return Result{}, false
		}
	}

	for i, t := range s.h {
		next := gstate{
			res:    s.res,
			h:      s.h.Remove(i),
			build:  s.build.Append(t),
			shared: s.shared,
		}

		if len(next.build) == 3 {
			if next.build.IsPeng() {
				next.res.Pengs = next.res.Pengs.Append(next.build[0])
				next.build = nil
			} else if next.build.IsChi() {
				next.res.Chis = next.res.Chis.Append(next.build[0])
				next.build = nil
			} else {
				// Failed build
				continue
			}
		}

		result, ok := next.step()
		if ok {
			return result, ok
		}
	}
	return Result{}, false
}
