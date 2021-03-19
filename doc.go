// Package mj contains types and tools for working with the game of Mahjong.
//
// It contains data types and structures that can represent tiles and collections of tiles.
// It also contains hand optimisers that finds the optimal grouping of tiles into sets.
//
// All types have sensible zero values representing an empty hand, and the methods will
// behave accordingly. However, with the exception of Hand, the types are also immutable once
// declared and the (exported) methods will not mutate their receivers or cause aliasing.
package mj
