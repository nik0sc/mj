// Package handcheck contains several hand optimisers. The purpose of a hand optimiser
// is to find the optimal grouping of tiles into melds and pairs for a given hand. An
// optimal grouping minimises free tiles, then minimises pairs. There may be more than
// one solution for a hand.
//
// The optimisers will not detect special hands like thirteen orphans or all pairs.
// They will also not group tiles into gangs/kongs. This is by design.
package handcheck
