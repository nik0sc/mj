package special

import (
	"sort"

	"github.com/nik0sc/mj"
)

var thirteenPure = mj.MustParseHand("b1 b9 c1 c9 w1 w9 he hs hw hn hz hf hb")

// IsThirteenOrphans returns true if the hand is in a waiting state for the
// Thirteen Orphans win (all honours, all terminals, and one other tile to make a pair).
// Additionally, the waited tile is returned; if the hand is "pure", the zero Tile is returned.
func IsThirteenOrphans(hand mj.Hand) (ok bool, wait mj.Tile) {
	if len(hand) != len(thirteenPure) {
		return
	}

	h := make(mj.Hand, len(hand))
	copy(h, hand)
	sort.Sort(h)

	for i := 0; i < len(thirteenPure); i++ {
		if h[i] != thirteenPure[i] {
			if wait.Valid() {
				// too many tiles different
				return false, mj.Tile{}
			}
			wait = thirteenPure[i]
		}
	}

	ok = true
	return
}
