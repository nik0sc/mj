package win

import (
	"errors"

	"github.com/nik0sc/mj"
	"github.com/nik0sc/mj/handcheck"
	"github.com/nik0sc/mj/wait"
)

const (
	NumTilesInHand = 13
)

type Edit struct {
	Old mj.Tile
	New mj.Tile
}

type EditSequence []Edit

// FindPath takes an input mj.Hand representing the player's hand, including
// opened melds, and an mj.Counter representing the tiles available for play.
// It returns the shortest sequence of edits required to get to a hand in the
// waiting state (one tile from winning), and the tiles that can be waited for.
//
// If `depth` > 0, at most `depth` number of edits is allowed (ie. limits
// search depth).
func FindPath(h mj.Hand, available mj.Counter, depth int) (EditSequence, []mj.Tile, error) {
	if h.Len() != NumTilesInHand {
		return nil, nil, errors.New("not enough tiles")
	}

	if !h.Valid() || !available.Valid() {
		return nil, nil, errors.New("failed validation")
	}

	c := handcheck.OptHandRLEChecker{UseMemo: true}
	group := c.Check(h)

	waits := wait.Find(group)
	if waits != nil {
		return nil, waits, nil
	}
	panic("unimplemented")
}

// TODO: Refactor wait.Find so that the free->meld promotion code can be used here
//   ... or maybe move both that and this function into one package?
