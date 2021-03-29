package wait

import (
	"sort"
	"testing"

	"github.com/nik0sc/mj"
)

func Test_Find(t *testing.T) {
	tests := []struct {
		name string
		res  mj.Group
		want []mj.Tile
	}{
		{
			"empty",
			mj.Group{},
			[]mj.Tile{},
		},
		{
			"peng",
			mj.Group{
				Pengs: mj.MustParseHand("b1 b2 b3"),
				Pairs: mj.MustParseHand("b4 b5"),
			},
			mj.MustParseHand("b4 b5"),
		},
		{
			"peng impossible",
			mj.Group{
				// b1 b2 b3 b3 b4 b5 b3 b4 b5
				Chis: mj.MustParseHand("b1 b3 b3"),
				// b4 b4 b5 b5
				Pairs: mj.MustParseHand("b4 b5"),
				// b1:1 b2:1 b3:3 b4:4 b5:b4
				// It's not possible to wait for another b4 or b5.
				// This could also be Gang:{b4 b5} Chi:{b1} Pair:{b3}
				// and a human player would probably prefer this grouping.
				// But we are algorithms.
			},
			[]mj.Tile{},
		},
		{
			"chi",
			mj.Group{
				// b1 b2 b3 b2 b3 b4 b3 b4 b5
				Chis: mj.MustParseHand("b1 b2 b3"),
				// b5 b5
				Pairs: mj.MustParseHand("b5"),
				Free:  mj.MustParseHand("b7 b8"),
				// b1:1 b2:2 b3:3 b4:2 b5:3 b7:1 b8:1
			},
			mj.MustParseHand("b6 b9"),
		},
		{
			"chi high",
			mj.Group{
				Pengs: mj.MustParseHand("b1 b2 b3"),
				Pairs: mj.MustParseHand("b5"),
				Free:  mj.MustParseHand("b8 b9"),
			},
			mj.MustParseHand("b7"),
		},
		{
			"chi low",
			mj.Group{
				Pengs: mj.MustParseHand("b7 b8 b9"),
				Pairs: mj.MustParseHand("b5"),
				Free:  mj.MustParseHand("b1 b2"),
			},
			mj.MustParseHand("b3"),
		},
		{
			"chi impossible",
			mj.Group{
				Chis:  mj.MustParseHand("b3 b3 b4"),
				Pairs: mj.MustParseHand("b3"),
				Free:  mj.MustParseHand("b1 b2"),
				// b1:1 b2:1 b3:4 b4:3 b5:3 b6:1
			},
			[]mj.Tile{},
		},
		{
			"chi middle",
			mj.Group{
				Pengs: mj.MustParseHand("b2 b3 b7"),
				Pairs: mj.MustParseHand("b1"),
				Free:  mj.MustParseHand("b4 b6"),
			},
			mj.MustParseHand("b5"),
		},
		{
			"chi wrong suit 1",
			mj.Group{
				Chis:  mj.MustParseHand("b1 b2 b3"),
				Pairs: mj.MustParseHand("b5"),
				Free:  mj.MustParseHand("b7 c8"),
			},
			[]mj.Tile{},
		},
		{
			"chi wrong suit 2",
			mj.Group{
				Chis:  mj.MustParseHand("b1 b2 b3"),
				Pairs: mj.MustParseHand("b5"),
				Free:  mj.MustParseHand("hz hf"),
			},
			[]mj.Tile{},
		},
		{
			"chi too far",
			mj.Group{
				Chis:  mj.MustParseHand("b1 b2 b3"),
				Pairs: mj.MustParseHand("b5"),
				Free:  mj.MustParseHand("b6 b9"),
			},
			[]mj.Tile{},
		},
		{
			"pair",
			mj.Group{
				Pengs: mj.MustParseHand("b2 b3 b4 b5"),
				Free:  mj.MustParseHand("b1"),
			},
			mj.MustParseHand("b1"),
		},
		{
			"pair impossible",
			mj.Group{
				Pengs: mj.MustParseHand("b1 b2 b3 b4"),
				Free:  mj.MustParseHand("b1"),
				// b1:4 b2:3 b3:3 b4:3
			},
			[]mj.Tile{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Find(tt.res)
			// sort both
			goth := mj.Hand(got)
			wanth := mj.Hand(tt.want)
			sort.Sort(goth)
			sort.Sort(wanth)

			if goth.Marshal() != wanth.Marshal() {
				t.Fatalf("want %s, got %s", wanth.String(), goth.String())
			}
		})
	}
}
