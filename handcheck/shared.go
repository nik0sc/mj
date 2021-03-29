package handcheck

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/nik0sc/mj"
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

func (s *shared) setMemo(repr string, g mj.Group) {
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

func (s *shared) getMemo(repr string) (mj.Group, bool) {
	//if s.memo == nil {
	//	return Group{}, false
	//}
	if g, ok := s.memo[repr]; ok {
		if writeMetrics {
			s.memoHits++
		}
		return mj.UnmarshalGroup(g), true
	}
	return mj.Group{}, false
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

// postprocessCountGroup does some cleanup that is common to the count-type checkers.
// It sorts the peng, chi and pair groups, then derives the value of Free by
// subtracting the formed groups from a map of tiles to counts.
func postprocessCountGroup(g *mj.Group, cmap map[mj.Tile]int) error {
	sort.Sort(g.Pengs)
	sort.Sort(g.Chis)
	sort.Sort(g.Pairs)

	for _, t := range g.Pengs {
		cmap[t] -= 3
	}

	for _, t := range g.Pairs {
		cmap[t] -= 2
	}

	for _, t := range g.Chis {
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
		return errors.New("cannot recreate Counter from map: " + err.Error())
	}
	g.Free = freecnt.ToHand(true)
	return nil
}
