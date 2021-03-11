package handcheck

import (
	"reflect"
	"testing"

	"mj"
)

func Test_handChecker_Check(t *testing.T) {
	type args struct {
		split    bool
		failfast bool
	}
	tests := []struct {
		name string
		hand string
		args args
		want GreedyResult
	}{
		{
			name: "all p",
			hand: "b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1",
			args: args{false, false},
			want: GreedyResult{true, 4, 0, 1, 0},
		},
		{
			"all c",
			"b1 b2 b3 b3 b4 b5 b5 b6 b7 b7 b8 b9 b9 b9",
			args{false, false},
			GreedyResult{true, 0, 4, 1, 0},
		},
		{
			"not simple",
			"w1 b7 w4 c5 b9 he w5 hf w5 c3 b8 hf hn hf",
			args{false, true},
			GreedyResult{},
		},
		{
			"not simple, full result",
			"w1 b7 w4 c5 b9 he w5 hf w5 c3 b8 hf hn hf",
			args{false, false},
			GreedyResult{false, 1, 1, 1, 0},
		},
		{
			"not simple either",
			"c1 c2 c3 c3 c3 c4 c5 c6",
			args{false, false},
			GreedyResult{true, 0, 2, 1, 0},
		},
		{
			"not simple 2",
			"c1 c1 c1 c2 c3 c4 c5 c6",
			args{false, false},
			GreedyResult{true, 0, 2, 1, 0},
		},
		{
			"degen",
			"b1 b3 b5 b7 b9 c1 c3 c5 c7 c9 w1 w3 w5 w7",
			args{false, false},
			GreedyResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GreedyChecker{tt.args.split, tt.args.failfast}

			h, err := mj.ParseHand(tt.hand)
			if err != nil {
				t.Fatalf("invalid hand: " + err.Error())
			}
			t.Logf("%s repr: %x", h.String(), h.Repr())

			if got := c.Check(h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
