package special

import (
	"reflect"
	"testing"

	"github.com/nik0sc/mj"
)

func TestIsSevenPairs(t *testing.T) {
	type args struct {
		hand        mj.Hand
		allowRepeat bool
	}
	tests := []struct {
		name     string
		args     args
		wantOk   bool
		wantWait mj.Tile
	}{
		{
			"No repeat 1",
			args{
				mj.MustParseHand("b1 b1 b2 b2 b3 b3 b4 b4 b5 b5 b6 b6 b7"),
				false,
			},
			true,
			mj.Tile{Suit: mj.Bamboo, Value: 7},
		},
		{
			"No repeat 2",
			args{
				mj.MustParseHand("b1 b1 b2 b2 b3 b3 b3 b3 b4 b4 b5 b5 b6"),
				false,
			},
			false,
			mj.Tile{},
		},
		{
			"No repeat 3",
			args{
				mj.MustParseHand("b1 b1 b1 b1 b2 b2 b3 b3 b4 b4 b5 b5 b6"),
				false,
			},
			false,
			mj.Tile{},
		},
		{
			"Repeat 1",
			args{
				mj.MustParseHand("c1 c1 c2 c2 c2 c2 w1 w1 hf hf w2 w2 w3"),
				true,
			},
			true,
			mj.Tile{Suit: mj.Wan, Value: 3},
		},
		{
			"Repeat 2",
			args{
				mj.MustParseHand("c1 c1 c2 c2 c2 c2 w1 w1 hf hf w2 w2 w2"),
				true,
			},
			true,
			mj.Tile{Suit: mj.Wan, Value: 2},
		},
		{
			"Wrong 1",
			args{
				mj.MustParseHand("c1 c2 c3 c4 c5 c6 c7 c8 c9 w1 w2 w3 w4"),
				false,
			},
			false,
			mj.Tile{},
		},
		{
			"Wrong 2",
			args{
				mj.MustParseHand("b1 b1 b1 b2 b2 b2 b3 b3 b4 b4 b5 b5 b6"),
				true,
			},
			false,
			mj.Tile{},
		},
		{
			"Wrong 3",
			args{
				mj.MustParseHand("b1 b1 b1 b1 b1 b2 b2 b3 b3 b4 b4 b5 b5"),
				true,
			},
			false,
			mj.Tile{},
		},
		{
			"Short 1",
			args{
				mj.Hand{},
				true,
			},
			false,
			mj.Tile{},
		},
		{
			"Short 2",
			args{
				mj.MustParseHand("b1 b1 b1 b1"),
				true,
			},
			false,
			mj.Tile{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotWait := IsSevenPairs(tt.args.hand, tt.args.allowRepeat)
			if gotOk != tt.wantOk {
				t.Errorf("IsSevenPairs() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotWait, tt.wantWait) {
				t.Errorf("IsSevenPairs() gotWait = %v, want %v", gotWait, tt.wantWait)
			}
		})
	}
}
