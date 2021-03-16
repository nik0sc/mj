package handcheck

import "fmt"

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
