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

// discard is alway true in this game
func DetectMeld(hasPair bool, hand Tiles, discard *Tile) *Meld {
	// same rank melds
	var matches Tiles
	for _, h := range hand {
		if !h.Dummy && discard.Suit == h.Suit && discard.Rank == h.Rank {
			matches = append(matches, h)
		}
	}
	// 4 in hand can't be a meld, max = 3 (using discard)
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

	// Chi
	if discard.Suit == "Man" || discard.Suit == "Pin" || discard.Suit == "Sou" {
		rank := discard.Rank
		// same-suit tiles by rank
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

/*

// find the tile that has the least meld possibility
func (m *Meld) EstimateValue() int {
	switch m.kind {
	case MeldPair:
		return 1
	case MeldPon:
		return 5
	case MeldKan:
		return 12
	case MeldChi:
		return 3
	default:
		return 0
	}
}

// EvaluateTile gives a composite usefulness score for a tile.
// Higher score = better tile to KEEP.
// Lower score = worse tile to DISCARD.
func EvaluateTile(p *player, t Tile) int {
	if t.Dummy {
		return -999 // never keep Dummy
	}

	score := 0

	// 1. Meld potential (DetectMeld + EstimateValue)
	if meld := DetectMeld(p.hasPair, p.tiles, &t); meld != nil {
		score += meld.EstimateValue() // 1,3,5,12
	}

	// 2. Adjacency (rank ±1)
	score += adjacencyScore(p, t)

	// 3. Duplicate bonus (same tile)
	score += duplicateScore(p, t)

	// 4. Suit clustering (same suit neighbors)
	score += suitClusterScore(p, t)

	// 5. Pair protection (if no pair yet)
	if !p.hasPair {
		score += pairProtection(p, t)
	}

	// 6. Kan protection (3 copies in hand)
	score += kanProtection(p, t)

	return score
}
func adjacencyScore(p *player, t Tile) int {
	bonus := 0
	for _, x := range p.tiles {
		if x.Dummy || x.Suit != t.Suit {
			continue
		}
		if abs(x.Rank-t.Rank) == 1 {
			bonus += 2
		}
	}
	return bonus
}
func duplicateScore(p *player, t Tile) int {
	count := 0
	for _, x := range p.tiles {
		if x.Dummy {
			continue
		}
		if x.Suit == t.Suit && x.Rank == t.Rank {
			count++
		}
	}
	if count >= 2 {
		return 3 // strong duplicate
	}
	return 0
}
func suitClusterScore(p *player, t Tile) int {
	for _, x := range p.tiles {
		if x.Dummy {
			continue
		}
		if x.Suit == t.Suit {
			return 1
		}
	}
	return 0
}
func pairProtection(p *player, t Tile) int {
	count := 0
	for _, x := range p.tiles {
		if x.Dummy {
			continue
		}
		if x.Suit == t.Suit && x.Rank == t.Rank {
			count++
		}
	}
	if count == 2 {
		return 4 // protect the pair
	}
	return 0
}
func kanProtection(p *player, t Tile) int {
	count := 0
	for _, x := range p.tiles {
		if x.Dummy {
			continue
		}
		if x.Suit == t.Suit && x.Rank == t.Rank {
			count++
		}
	}
	if count == 3 {
		return 6 // protect potential Kan
	}
	return 0
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
*/
