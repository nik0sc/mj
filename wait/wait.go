package wait

import (
	"fmt"
	"sort"

	"github.com/nik0sc/mj"
	"github.com/nik0sc/mj/handcheck"
)

// Find takes an input handcheck.Result and determines what tiles
// the player could wait for to win. Tile counts within the hand are
// considered, but there is no consideration of discarded tile counts.
func Find(result handcheck.Result) []mj.Tile {
	meldsets := result.Chis.Len() + result.Pengs.Len()
	cnt := result.ToCount()
	var waits []mj.Tile

	if result.Free.Len() == 0 && meldsets == 3 && result.Pairs.Len() == 2 {
		// check for peng from pairs
		// Either pair could be possible
		for _, t := range result.Pairs {
			if cnt.Get(t) < 4 && t.CanMeld() {
				waits = append(waits, t)
			}
		}
		return waits
	} else if result.Free.Len() == 2 && meldsets == 3 && result.Pairs.Len() == 1 {
		// check for chi
		if !result.Free[0].IsBasic() || !result.Free[1].IsBasic() {
			return nil
		}

		if result.Free[0].Suit != result.Free[1].Suit {
			return nil
		}

		// mutates result.Free passed in: this is probably not a big deal
		sort.Sort(result.Free)

		if result.Free[0].Value+1 == result.Free[1].Value {
			// check both ends
			if result.Free[0].Value > 1 {
				t := result.Free[0]
				t.Value--
				if !t.Valid() {
					panic(fmt.Sprintf("invalid low tile: %+v", t))
				}
				if cnt.Get(t) < 4 {
					waits = append(waits, t)
				}
			}
			if result.Free[1].Value < 9 {
				t := result.Free[1]
				t.Value++
				if !t.Valid() {
					panic(fmt.Sprintf("invalid high tile: %+v", t))
				}
				if cnt.Get(t) < 4 {
					waits = append(waits, t)
				}
			}
			return waits
		} else if result.Free[0].Value+2 == result.Free[1].Value {
			// check middle
			t := result.Free[0]
			t.Value++
			if !t.Valid() {
				panic(fmt.Sprintf("invalid middle tile: %+v", t))
			}
			if cnt.Get(t) < 4 {
				waits = append(waits, t)
			}
			return waits
		}
		return nil
	} else if result.Free.Len() == 1 && meldsets == 4 && result.Pairs.Len() == 0 {
		// check for pair
		t := result.Free[0]
		if cnt.Get(t) < 4 && t.CanMeld() {
			waits = append(waits, t)
		}
		return waits
	} else {
		return nil
	}
}
