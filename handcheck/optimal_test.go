package handcheck

import (
	"sort"
	"testing"

	"github.com/nik0sc/mj"
)

func Test_OptChecker_Check(t *testing.T) {
	tests := []struct {
		name string
		hand string
		want mj.Group
	}{
		{
			"all p",
			"b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1",
			mj.Group{
				Pengs: []mj.Tile{
					{Suit: mj.Bamboo, Value: 1},
					{Suit: mj.Bamboo, Value: 1},
					{Suit: mj.Bamboo, Value: 1},
					{Suit: mj.Bamboo, Value: 1},
				},
				Chis:  nil,
				Pairs: []mj.Tile{{Suit: mj.Bamboo, Value: 1}},
				Free:  nil,
			},
		},
		{
			"all c",
			"b1 b2 b3 b3 b4 b5 b5 b6 b7 b7 b8 b9 b9 b9",
			mj.Group{
				Pengs: nil,
				Chis: []mj.Tile{
					{Suit: mj.Bamboo, Value: 1},
					{Suit: mj.Bamboo, Value: 3},
					{Suit: mj.Bamboo, Value: 5},
					{Suit: mj.Bamboo, Value: 7},
				},
				Pairs: []mj.Tile{{Suit: mj.Bamboo, Value: 9}},
				Free:  nil,
			},
		},
		{
			"not simple",
			"w1 b7 w4 c5 b9 he w5 hf w5 c3 b8 hf hn hf",
			mj.Group{
				Pengs: []mj.Tile{{Suit: mj.Honour, Value: mj.Fa}},
				Chis:  []mj.Tile{{Suit: mj.Bamboo, Value: 7}},
				Pairs: []mj.Tile{{Suit: mj.Wan, Value: 5}},
				Free: mj.Hand{
					{Suit: mj.Coin, Value: 3},
					{Suit: mj.Coin, Value: 5},
					{Suit: mj.Wan, Value: 1},
					{Suit: mj.Wan, Value: 4},
					{Suit: mj.Honour, Value: mj.East},
					{Suit: mj.Honour, Value: mj.North},
				},
			},
		},
		{
			"not simple either",
			"c1 c2 c3 c3 c3 c4 c5 c6",
			mj.Group{
				Pengs: nil,
				Chis: []mj.Tile{
					{Suit: mj.Coin, Value: 1},
					{Suit: mj.Coin, Value: 4},
				},
				Pairs: []mj.Tile{{Suit: mj.Coin, Value: 3}},
				Free:  nil,
			},
		},
		{
			"not simple 2",
			"c1 c1 c1 c2 c3 c4 c5 c6",
			mj.Group{
				Pengs: nil,
				Chis: []mj.Tile{
					{Suit: mj.Coin, Value: 1},
					{Suit: mj.Coin, Value: 4},
				},
				Pairs: []mj.Tile{{Suit: mj.Coin, Value: 1}},
				Free:  nil,
			},
		},
		{
			"degen",
			"b1 b3 b5 b7 b9 c1 c3 c5 c7 c9 w1 w3 w5 w7",
			mj.Group{
				Pengs: nil,
				Chis:  nil,
				Pairs: nil,
				Free: mj.Hand{
					{Suit: mj.Bamboo, Value: 1},
					{Suit: mj.Bamboo, Value: 3},
					{Suit: mj.Bamboo, Value: 5},
					{Suit: mj.Bamboo, Value: 7},
					{Suit: mj.Bamboo, Value: 9},
					{Suit: mj.Coin, Value: 1},
					{Suit: mj.Coin, Value: 3},
					{Suit: mj.Coin, Value: 5},
					{Suit: mj.Coin, Value: 7},
					{Suit: mj.Coin, Value: 9},
					{Suit: mj.Wan, Value: 1},
					{Suit: mj.Wan, Value: 3},
					{Suit: mj.Wan, Value: 5},
					{Suit: mj.Wan, Value: 7},
				},
			},
		},
		{
			"pairs or chis",
			"b1 b2 b3 b1 b2 b3",
			mj.Group{
				Pengs: nil,
				Chis: []mj.Tile{
					{Suit: mj.Bamboo, Value: 1},
					{Suit: mj.Bamboo, Value: 1},
				},
				Pairs: nil,
				Free:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := mj.ParseHand(tt.hand)
			if err != nil {
				t.Fatalf("invalid hand: " + err.Error())
			}
			t.Logf("%s repr: %x", h.String(), h.Marshal())

			rFast := OptChecker{Split: true, UseMemo: true}.Check(h)
			// order may be wrong
			sort.Sort(rFast.Free)

			if tt.want.Marshal() != rFast.Marshal() {
				t.Errorf("fast: want %v, got %v", tt.want, rFast)
			}

			rSlow := OptChecker{}.Check(h)
			sort.Sort(rSlow.Free)

			if tt.want.Marshal() != rSlow.Marshal() {
				t.Errorf("slow: want %v, got %v", tt.want, rSlow)
			}
		})
	}
}

func Benchmark_OptChecker_AllP(b *testing.B) {
	hand, _ := mj.ParseHand("b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1")
	benchmark_OptChecker(b, hand)

}

func Benchmark_OptChecker_AllPReal(b *testing.B) {
	hand, _ := mj.ParseHand("b1 b1 b1 b2 b2 b2 b3 b3 b3 b4 b4 b4 b5 b5")
	benchmark_OptChecker(b, hand)
}

func Benchmark_OptChecker_AllC(b *testing.B) {
	hand, _ := mj.ParseHand("b1 b2 b3 b3 b4 b5 b5 b6 b7 b7 b8 b9 b9 b9")
	benchmark_OptChecker(b, hand)

}

func Benchmark_OptChecker_NS(b *testing.B) {
	hand, _ := mj.ParseHand("w1 b7 w4 c5 b9 he w5 hf w5 c3 b8 hf hn hf")
	benchmark_OptChecker(b, hand)
}

func benchmark_OptChecker(b *testing.B, h mj.Hand) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = OptChecker{UseMemo: true}.Check(h)
	}
}
