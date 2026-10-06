package main

import (
	"fmt"
	"strings"
)

/*

  File:    evaluate.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

type MeldType int

type Melds []*Meld

const (
	MeldPair MeldType = iota
	MeldPon
	MeldKan
	MeldChi
)

func (m MeldType) String() string {
	switch m {
	case MeldPair:
		return "pair" // Toitsu
	case MeldPon:
		return "pon" // Koutsu
	case MeldKan:
		return "kan" // Kantsu
	case MeldChi:
		return "chi" // Shuntsu
	}
	panic("invalid MeldType")
}

type Meld struct {
	kind  MeldType
	open  bool
	tiles Tiles // includes discard
}

func (m *Meld) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Meld: %s, %d", m.kind.String(), len(m.tiles)))
	for i := range len(m.tiles) {
		s.WriteString(fmt.Sprintf(" %s", m.tiles[i].Code))
	}
	return s.String()
}

func NewMeld(kind MeldType, open bool, tiles Tiles) *Meld {
	return &Meld{
		kind:  kind,
		open:  open,
		tiles: tiles,
	}
}

func (m *Meld) Score() int {
	switch m.kind {
	case MeldPair:
		return 20
	case MeldPon:
		return 40
	case MeldKan:
		return 80
	case MeldChi:
		return 30
	default:
		return 0
	}
}

// create a Meld if the opponent's discard completes it.
func DetectMeld(hasPair bool, hand Tiles, discard *Tile) *Meld {
	// same rank melds
	var matches Tiles
	for _, h := range hand {
		if !h.Dummy && discard.Suit == h.Suit && discard.Rank == h.Rank {
			matches = append(matches, h)
		}
	}
	// 4 in hand isnt valid
	if len(matches) >= 4 {
		return nil
	}
	// Kan - 3 in hand + discard
	if len(matches) == 3 {
		matches = append(matches, discard)
		return NewMeld(MeldKan, false, matches)
	}
	// Pon - 2 in hand + discard
	if len(matches) == 2 {
		matches = append(matches, discard)
		return NewMeld(MeldPon, false, matches)
	}
	// Pair - 1 in hand + discard - only 1 allowed
	if len(matches) == 1 && !hasPair {
		matches = append(matches, discard)
		return NewMeld(MeldPair, false, matches)
	}

	// Chi same-suit tiles by rank
	if discard.Suit == Manzu || discard.Suit == Pinzu || discard.Suit == Souzu {
		rank := discard.Rank
		ranks := make(map[int]Tiles)
		for _, h := range hand {
			if h.Suit == discard.Suit {
				ranks[h.Rank] = append(ranks[h.Rank], h)
			}
		}
		// left (rank-2, rank-1)
		if len(ranks[rank-2]) > 0 && len(ranks[rank-1]) > 0 {
			tiles := Tiles{discard, ranks[rank-2][0], ranks[rank-1][0]}
			return NewMeld(MeldChi, false, tiles)
		}
		// middle (rank-1, rank+1)
		if len(ranks[rank-1]) > 0 && len(ranks[rank+1]) > 0 {
			tiles := Tiles{discard, ranks[rank-1][0], ranks[rank+1][0]}
			return NewMeld(MeldChi, false, tiles)
		}
		// right (rank+1, rank+2)
		if len(ranks[rank+1]) > 0 && len(ranks[rank+2]) > 0 {
			tiles := Tiles{discard, ranks[rank+1][0], ranks[rank+2][0]}
			return NewMeld(MeldChi, false, tiles)
		}
	}

	return nil
}
