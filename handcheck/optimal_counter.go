package handcheck

import (
	"fmt"

	"mj"
)

type OptCountChecker struct {
	memo map[string]Result

	// optimisations

	// Split hand into sub hands (melds are restricted to one hand only)
	Split bool
	// use memoisation to avoid O(2^n) running time
	UseMemo bool
}

type ocstate struct {
	res Result
	freeCount mj.Counter
	shared *ocshared
}

type ocshared struct {
	memo map[string]Result
	stepCount int
	memoHits int
}

func (c OptCountChecker) Check(hand mj.Hand) Result {
	panic("unimplemented")
}

func (c OptCountChecker) start(h mj.Hand) Result {
	shared := ocshared{}
	if c.UseMemo {
		shared.memo = make(map[string]Result)
	}
	s := ocstate{Result{Free: h}, h.ToCount(), &shared}

	r := s.step()
	if writeMetrics {
		fmt.Printf("shared: len(memo)=%d steps=%d memohits=%d\n",
			len(s.shared.memo), s.shared.stepCount, s.shared.memoHits)
	}

	return r
}

func (s ocstate) step() Result {
	panic("unimplemented")
}