package mj

import (
	"errors"
	"strings"
)

var suitParse = map[uint8]Suit{
	'b': Bamboo,
	'c': Coin,
	'w': Wan,
	'h': Honour,
	'f': Flower,
}

var honourParse = map[uint8]Value{
	'e': East,
	's': South,
	'w': West,
	'n': North,
	'z': Zhong,
	'f': Fa,
	'b': Ban,
}

// ParseTile turns a 2-character string into a Tile.
// The first character is the Suit and may be one of the characters "bcwhf" (for
// Bamboo, Coin, Wan, Honour and Flower).
// The second character is the Value and its permissible range depends on the Suit:
//  - Bamboo, Coin and Wan: a digit between 1-9 inclusive.
//  - Honour: one of the characters "eswnzfb" (for East, South, West, North,
//      Zhong, Fa and Ban).
//  - Flower: a digit between 1-8 inclusive.
// Parsing errors are returned in err.
func ParseTile(s string) (t Tile, err error) {
	var ok bool
	if len(s) != 2 {
		err = errors.New("tile representation must be 2 characters long")
		return
	}

	s = strings.ToLower(s)

	t.Suit, ok = suitParse[s[0]]
	if !ok {
		err = errors.New("unrecognised suit: " + string(s[0]))
		return
	}

	switch t.Suit {
	case Bamboo, Coin, Wan:
		t.Value = Value(s[1] - '0')
		if t.Value < 1 || t.Value > 9 {
			err = errors.New("invalid value for simple tile: " + string(s[1]))
		}
	case Honour:
		t.Value, ok = honourParse[s[1]]
		if !ok {
			err = errors.New("invalid value for honour suit: " + string(s[1]))
		}
	case Flower:
		t.Value = Value(s[1] - '0')
		if t.Value < 1 || t.Value > 8 {
			err = errors.New("invalid value for flower tile: " + string(s[1]))
		}
		t.Value |= FlowerBase
	default:
		panic("ParseTile: unreachable")
	}
	return
}

// ParseHand turns a space-separated string of 2-character sequences into a Hand, in order.
// Each 2-character sequence is passed to ParseTile.
func ParseHand(s string) (h Hand, err error) {
	s = strings.TrimSpace(s)
	ss := strings.Split(s, " ")
	h = make(Hand, len(ss))
	for i, t := range ss {
		h[i], err = ParseTile(t)
		if err != nil {
			return
		}
	}
	return
}
