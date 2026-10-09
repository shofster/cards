package main

import (
	"fmt"
	"math/rand"
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

// ------  possible discard evaluation  ------- \\

func sameRank(a, b *Tile) bool { return a.Rank == b.Rank }
func sameSuit(a, b *Tile) bool { return a.Suit == b.Suit }
func isAdj(a, b *Tile) bool    { return abs(a.Rank-b.Rank) == 1 }
func isNear(a, b *Tile) bool   { return abs(a.Rank-b.Rank) == 2 }
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func (t *Tile) IsHonor() bool {
	return t.Suit == WindSuit || t.Suit == DragonSuit
}
func (t *Tile) HasNoFriends(hand Tiles) bool {
	for _, h := range hand {
		if h.ID == t.ID {
			continue
		}
		if sameRank(t, h) {
			return false
		}
		if isAdj(t, h) {
			return false
		}
		if isNear(t, h) {
			return false
		}
		if sameSuit(t, h) {
			return false
		}
	}
	return true
}

func rateDiscard(candidate *Tile, hand Tiles, take *Tile) int {
	score := 0
	// how mush the candidate "likes" the new brother
	if sameRank(candidate, take) {
		score += 5
	}
	if isAdj(candidate, take) {
		score += 3
	}
	if isNear(candidate, take) {
		score += 1
	}
	if sameSuit(candidate, take) {
		score += 1
	}
	// count if the candidate is a 2'nd cousin
	for _, h := range hand {
		if h.ID == candidate.ID {
			continue
		}
		if sameRank(candidate, h) {
			score += 4
		}
		if isAdj(candidate, h) {
			score += 2
		}
		if isNear(candidate, h) {
			score += 1
		}
		if sameSuit(candidate, h) {
			score += 1
		}
	}
	// Penalize isolated honors (no meld potential)
	if candidate.IsHonor() && candidate.HasNoFriends(hand) {
		score -= 2
	}
	return score
}

func BestDiscard(hand Tiles, take *Tile) int {
	lowest := 999999
	var candidates []int
	for _, candidate := range hand {
		if !candidate.Dummy {
			score := rateDiscard(candidate, hand, take)
			if score < lowest {
				lowest = score
				candidates = []int{candidate.ID}
			} else if score == lowest {
				candidates = append(candidates, candidate.ID)
			}
		}
	}
	if len(candidates) == 0 {
		return -1
	}
	return candidates[rand.Intn(len(candidates))]
}
