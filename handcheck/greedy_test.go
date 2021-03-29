package handcheck

import (
	"reflect"
	"testing"

	"github.com/nik0sc/mj"
)

func Test_GreedyChecker_Check(t *testing.T) {
	type args struct {
		split    bool
		failfast bool
	}
	tests := []struct {
		name string
		hand string
		args args
		want Group
	}{
		{
			name: "all p",
			hand: "b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1",
			args: args{false, false},
			want: Group{
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
			args{false, false},
			Group{
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
			args{true, true},
			Group{},
		},
		{
			"not simple, full result",
			"w1 b7 w4 c5 b9 he w5 hf w5 c3 b8 hf hn hf",
			args{false, false},
			Group{},
		},
		{
			"not simple either",
			"c1 c2 c3 c3 c3 c4 c5 c6",
			args{false, false},
			Group{
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
			args{false, false},
			Group{
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
			args{false, false},
			Group{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GreedyChecker{tt.args.split, tt.args.failfast}

			h, err := mj.ParseHand(tt.hand)
			if err != nil {
				t.Fatalf("invalid hand: " + err.Error())
			}
			t.Logf("%s repr: %x", h.String(), h.Marshal())

			if got := c.Check(h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_GreedyChecker_AllP(b *testing.B) {
	hand, _ := mj.ParseHand("b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1")
	benchmark_GreedyChecker(b, hand)
}

func Benchmark_GreedyChecker_AllPReal(b *testing.B) {
	hand, _ := mj.ParseHand("b1 b1 b1 b2 b2 b2 b3 b3 b3 b4 b4 b4 b5 b5")
	benchmark_GreedyChecker(b, hand)
}

func Benchmark_GreedyChecker_AllC(b *testing.B) {
	hand, _ := mj.ParseHand("b1 b2 b3 b3 b4 b5 b5 b6 b7 b7 b8 b9 b9 b9")
	benchmark_GreedyChecker(b, hand)
}

func Benchmark_GreedyChecker_NS(b *testing.B) {
	hand, _ := mj.ParseHand("w1 b7 w4 c5 b9 he w5 hf w5 c3 b8 hf hn hf")
	benchmark_GreedyChecker(b, hand)
}

func benchmark_GreedyChecker(b *testing.B, h mj.Hand) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = GreedyChecker{}.Check(h)
	}
}
