package mj

const (
	Bamboo Suit = iota + 1
	Coin
	Wan
	Honour
	Flower
)

// Suit is the suit of a Tile. The zero Suit is invalid. There are three basic suits,
// Bamboo, Coin and Wan, as well as the Honour and Flower suits.
type Suit byte

const (
	East Value = iota + 10
	South
	West
	North
	Zhong
	Fa
	Ban

	FlowerBase Value = 32
)

// Value is the face value of a Tile, including honours and bonuses. The zero Value is invalid.
// Values 1-9 inclusive are used for the basic suits. East, South, West, North, Zhong, Fa and Ban
// are only valid for the Honour suit. Values 32-39 inclusive are only valid for the Flower suit.
// Value 32 is defined as FlowerBase.
type Value byte

// Tile is a single tile played in mahjong, comprising a Suit and a Value.
type Tile struct {
	Suit
	Value
	// This representation allows us to read suit and value cheaply,
	// while still maintaining a balance with compactness.
}

const (
	uniTileBack    = '🀫'
	uniTileEast    = '🀀'
	uniTileWan1    = '🀇'
	uniTileBamboo1 = '🀐'
	uniTileCoin1   = '🀙'
	uniTileFlower1 = '🀢'
	// used to force emoji style representation
	uniVS16 rune = 0xfe0f
	// enable/disable emoji style for all tiles
	// (this seems to work on Windows renderers and nowhere else)
	uniUseVS16 = true
)

// Valid returns true if the Tile data is valid and may be used in the algorithms.
func (t Tile) Valid() bool {
	if t.Suit == 0 || t.Value == 0 {
		return true
	}

	switch t.Suit {
	case Bamboo, Coin, Wan:
		return 1 <= t.Value && t.Value <= 9
	case Honour:
		return East <= t.Value && t.Value <= Ban
	case Flower:
		return FlowerBase <= t.Value && t.Value < (FlowerBase+8)
	}

	return false
}

// String returns a Unicode human-readable representation of the Tile.
// These strings require up to 7 bytes to encode in utf-8.
// For a space efficient encoding check out Marshal().
func (t Tile) String() string {
	// base mahjong tile requires 4 bytes in utf-8
	// vs16 requires another 3 bytes
	// eg. b1 in utf-8 is f0 9f 80 90
	// and vs16 in utf-8 is ef b8 8f
	if !t.Valid() {
		if uniUseVS16 {
			return string([]rune{uniTileBack, uniVS16})
		} else {
			return string(uniTileBack)
		}
	}

	var base rune
	var offset Value

	switch t.Suit {
	case Bamboo:
		base = uniTileBamboo1
		offset = t.Value - 1
	case Coin:
		base = uniTileCoin1
		offset = t.Value - 1
	case Wan:
		base = uniTileWan1
		offset = t.Value - 1
	case Honour:
		base = uniTileEast
		offset = t.Value - East
	case Flower:
		base = uniTileFlower1
		offset = t.Value - FlowerBase
	}

	if uniUseVS16 {
		return string([]rune{base + rune(offset), uniVS16})
	} else {
		return string(base + rune(offset))
	}
}

// Marshal returns an unambiguous encoding for a Tile packed into a byte.
func (t Tile) Marshal() byte {
	// The encoding:
	// 0bxxxyyyyy
	// x: The suit, unchanged
	// y: The low 5 bits of the value
	// Working with this encoding would be a pain
	// with constant masking and shifting to access the
	// suit and value independently. So, the default
	// in-memory working representation is a 2-byte value,
	// but we can pack a tile into one byte for space-
	// sensitive applications.
	return (byte(t.Suit) << 5) | (byte(t.Value) & 31)
}

// CanMeld returns true if the Tile may participate in melds.
func (t Tile) CanMeld() bool {
	return t.Valid() && t.Suit != Flower
}

// IsBasic returns true if the Tile is a basic tile.
func (t Tile) IsBasic() bool {
	return t.Valid() && t.Suit != Flower && t.Suit != Honour
}

// UnmarshalTile is the inverse of Tile.Marshal().
func UnmarshalTile(b byte) Tile {
	t := Tile{
		Suit:  Suit(b >> 5),
		Value: Value(b & 31),
	}

	if t.Suit == Flower {
		t.Value |= FlowerBase
	}

	return t
}
