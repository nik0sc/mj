package special

import (
	"reflect"
	"testing"

	"github.com/nik0sc/mj"
)

func TestIsThirteenOrphans(t *testing.T) {
	tests := []struct {
		name     string
		hand     mj.Hand
		wantOk   bool
		wantWait mj.Tile
	}{
		{
			"Pure",
			mj.MustParseHand("b1 b9 c1 c9 w1 w9 he hs hw hn hz hf hb"),
			true,
			mj.Tile{},
		},
		{
			"Impure 1",
			mj.MustParseHand("b1 b1 c1 c9 w1 w9 he hs hw hn hz hf hb"),
			true,
			mj.Tile{Suit: mj.Bamboo, Value: 9},
		},
		{
			"Impure 2",
			mj.MustParseHand("b1 b9 c1 c9 w1 w9 he hs hw hn hz hf hf"),
			true,
			mj.Tile{Suit: mj.Honour, Value: mj.Ban},
		},
		{
			"Impure 3",
			mj.MustParseHand("b1 b9 c1 c9 w1 w9 he hs hw hn hz hb hb"),
			true,
			mj.Tile{Suit: mj.Honour, Value: mj.Fa},
		},
		{
			"Wrong 1",
			mj.MustParseHand("b1 b2 b9 c1 c9 w1 he hs hw hn hz hf hb"),
			false,
			mj.Tile{},
		},
		{
			"Wrong 2",
			mj.MustParseHand("b1 b2 b3 b4 b5 b6 b7 b8 b9 b1 b2 b3 b4"),
			false,
			mj.Tile{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotWait := IsThirteenOrphans(tt.hand)
			if gotOk != tt.wantOk {
				t.Errorf("IsThirteenOrphans() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotWait, tt.wantWait) {
				t.Errorf("IsThirteenOrphans() gotWait = %v, want %v", gotWait, tt.wantWait)
			}
		})
	}
}
