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

func (s *shared) setMemo(repr string, r Result) {
	if s.memo == nil {
		return
	}
	if rOld, ok := s.memo[repr]; ok {
		// memo should not be updated like this! because the memo result should already be optimal for
		// the currently free tiles
		panic(fmt.Sprintf("updating memo: repr=%x rOld=%+v r=%+v", repr, rOld, r))
	}
	// sort the result first
	store := r.Copy(true)

	s.memo[repr] = store.Marshal()
}

func (s *shared) getMemo(repr string) (Result, bool) {
	//if s.memo == nil {
	//	return Result{}, false
	//}
	if r, ok := s.memo[repr]; ok {
		if writeMetrics {
			s.memoHits++
		}
		return UnmarshalResult(r), true
	}
	return Result{}, false
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
