package handcheck

import (
	"fmt"
	"io"
)

const (
	// Record metrics in shared struct.
	writeMetrics = true
	// Trace execution of each step. Very slow.
	traceSteps = false
)

type shared struct {
	memo      map[string]string
	stepCount int
	memoHits  int
}

func (s *shared) setMemo(repr string, g Group) {
	if s.memo == nil {
		return
	}
	if rOld, ok := s.memo[repr]; ok {
		// memo should not be updated like this! because the memo result should already be optimal for
		// the currently free tiles
		panic(fmt.Sprintf("updating memo: repr=%x rOld=%+v g=%+v", repr, rOld, g))
	}
	// sort the result first
	store := g.Copy(true)

	s.memo[repr] = store.Marshal()
}

func (s *shared) getMemo(repr string) (Group, bool) {
	//if s.memo == nil {
	//	return Group{}, false
	//}
	if g, ok := s.memo[repr]; ok {
		if writeMetrics {
			s.memoHits++
		}
		return UnmarshalGroup(g), true
	}
	return Group{}, false
}

func (s *shared) enterStep(w io.Writer, at fmt.Stringer) {
	if writeMetrics {
		s.stepCount++
	}
	if traceSteps {
		_, _ = fmt.Fprintf(w, "at %s\n", at.String())
	}
}

func (s *shared) writeSummary(writer io.Writer) {
	if writeMetrics {
		_, _ = fmt.Fprintf(writer, "shared: len(memo)=%d steps=%d memohits=%d\n",
			len(s.memo), s.stepCount, s.memoHits)
	}
}
