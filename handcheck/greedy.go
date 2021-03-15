package handcheck

import (
	"fmt"
	"sort"

	"github.com/nik0sc/mj"
)

type GreedyChecker struct {
	// Split=true breaks the guarantee that if we return ok=false there
	// is no possible winning interpretation of the hand
	Split    bool
	FailFast bool
}

type gstate struct {
	res GreedyResult

	h     mj.Hand
	build mj.Hand

	shared *gshared
}

type gshared struct {
	stepCount int
}

func (c GreedyChecker) Check(hand mj.Hand) GreedyResult {
	h := make(mj.Hand, len(hand))
	copy(h, hand)
	sort.Sort(h)

	if c.Split {
		hsplit := h.Split(false)
		rs := make([]GreedyResult, 0, len(hsplit))
		for _, hs := range hsplit {
			rs = append(rs, c.start(hs))
		}

		r := GreedyResult{Ok: true}
		for _, rsub := range rs {
			if !rsub.Ok {
				if c.FailFast {
					r = GreedyResult{}
					break
				} else {
					r.Ok = false
				}
			}
			r.Peng += rsub.Peng
			r.Chi += rsub.Chi
			r.Pair += rsub.Pair
		}
		return r
	} else {
		return c.start(h)
	}
}

func (c GreedyChecker) start(h mj.Hand) GreedyResult {
	_ = c
	s := gstate{GreedyResult{}, h, nil, &gshared{}}

	r := s.step()
	if writeMetrics {
		fmt.Printf("shared: steps=%d\n", s.shared.stepCount)
	}

	return r
}

func (s gstate) step() GreedyResult {
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
			return GreedyResult{true, s.res.Peng, s.res.Chi, 1, 0}
		} else {
			return GreedyResult{false, s.res.Peng, s.res.Chi, 0, 0}
		}
	}

	for i, t := range s.h {
		hNew := s.h.Remove(i)

		buildNew := make(mj.Hand, len(s.build)+1)
		copy(buildNew, s.build)
		buildNew[len(s.build)] = t

		next := gstate{
			res:    s.res,
			h:      hNew,
			build:  buildNew,
			shared: s.shared,
		}

		if len(next.build) == 3 {
			if next.build.IsPeng() {
				next.res.Peng++
				next.build = nil
			} else if next.build.IsChi() {
				next.res.Chi++
				next.build = nil
			} else {
				// Failed build
				continue
			}
		}

		result := next.step()
		if result.Ok {
			return result
		}
	}
	return GreedyResult{}
}
