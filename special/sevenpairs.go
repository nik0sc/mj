package special

import "github.com/nik0sc/mj"

// IsSevenPairs returns true if the hand is in a waiting state for
// the Seven Pairs win (six pairs and one tile that could be waited).
// If allowRepeat is true, pairs may be repeated once.
// Additionally, the waited tile is returned.
func IsSevenPairs(hand mj.Hand, allowRepeat bool) (ok bool, wait mj.Tile) {
	if len(hand) != 13 {
		return
	}

	// could this be done without Counter?
	ok = true
	hand.ToCount().ForEach(func(t mj.Tile, n int) bool {
		switch n {
		case 4:
			if !allowRepeat {
				ok = false
				return false
			}
			fallthrough
		case 2:
			return true
		case 1, 3:
			if !wait.Valid() {
				wait = t
				return true
			}
			fallthrough
		default:
			ok = false
			return false
		}
	})

	if !ok {
		// for hands with multiple single tiles
		wait = mj.Tile{}
	} else if !wait.Valid() {
		panic("ok but invalid wait tile")
	}

	return
}
