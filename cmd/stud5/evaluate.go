package main

import (
	"cards/common"
	"sort"
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
A random 5‑card hand is overwhelmingly likely to be:
- High card — about 50% of all deals
- One pair — about 42%
- Two pair — about 4.75%
- Trips — about 2.1%
- Straight — about 0.39%
- Flush — about 0.20%
- Full house — about 0.14%
- Quads — about 0.024%
- Straight flush — about 0.0015%
So more than 90% of all hands are either high card or a single low pair.

- Dealer qualifies ~42% of the time
- Player wins antes on the other ~58%
- When the dealer qualifies, the dealer wins more often than the player
- Overall house edge is around 5%



*/

type HandRank int

type FinalRank struct {
	Rank   HandRank
	Kicker []int
}

const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (h HandRank) String() string {
	switch h {
	case RoyalFlush:
		return "Royal Flush"
	case StraightFlush:
		return "Straight Flush"
	case FourOfAKind:
		return "Four of a Kind"
	case FullHouse:
		return "Full House"
	case Flush:
		return "Flush"
	case Straight:
		return "Straight"
	case ThreeOfAKind:
		return "Three of a Kind"
	case TwoPair:
		return "Two Pair"
	case OnePair:
		return "One Pair"
	case HighCard:
		return "Ace-King"
	default:
		return "Unknown"
	}
}

func DealerQualifies(f FinalRank) bool {
	switch f.Rank {
	case OnePair, TwoPair, ThreeOfAKind, Straight, Flush, FullHouse,
		FourOfAKind, StraightFlush, RoyalFlush:
		return true
	case HighCard: // Ace-King
		return f.Kicker[0] == 14 && f.Kicker[1] == 13
	default:
		return false
	}
}

// > 0 player wins
// < 0 dealer wins
// 0 == tie (push)
func CompareFinal(player, dealer FinalRank) int {
	if player.Rank != dealer.Rank {
		return int(player.Rank) - int(dealer.Rank)
	}
	for i := 0; i < len(player.Kicker) && i < len(dealer.Kicker); i++ {
		if player.Kicker[i] != dealer.Kicker[i] {
			return player.Kicker[i] - dealer.Kicker[i]
		}
	}
	return 0
}

func EvaluateFinal(cards []common.Card) FinalRank {
	// Count ranks and suits
	rankCount := make(map[int]int)
	suitCount := make(map[int]int)
	ranks := make([]int, 0, 5)

	for _, c := range cards {
		r := int(c.Rank)
		s := int(c.Suit)
		rankCount[r]++
		suitCount[s]++
		ranks = append(ranks, r)
	}

	sort.Ints(ranks)

	// Straight detection
	isStraight := false
	if ranks[4]-ranks[0] == 4 &&
		rankCount[ranks[0]] == 1 &&
		rankCount[ranks[1]] == 1 &&
		rankCount[ranks[2]] == 1 &&
		rankCount[ranks[3]] == 1 &&
		rankCount[ranks[4]] == 1 {
		isStraight = true
	}
	// Ace-low straight (A 2 3 4 5) wheel
	if ranks[0] == 2 && ranks[1] == 3 && ranks[2] == 4 && ranks[3] == 5 && ranks[4] == 14 {
		isStraight = true
		ranks = []int{5, 4, 3, 2, 1}
	}

	// Detect flush
	isFlush := false
	for _, n := range suitCount {
		if n == 5 {
			isFlush = true
			break
		}
	}

	// Straight flush / royal flush
	if isStraight && isFlush {
		if ranks[4] == 14 && ranks[3] == 13 {
			return FinalRank{Rank: RoyalFlush, Kicker: ranks}
		}
		return FinalRank{Rank: StraightFlush, Kicker: ranks}
	}

	// Count same card(s)
	pairs := []int{}
	trips := []int{}
	quads := []int{}

	for r, n := range rankCount {
		switch n {
		case 4:
			quads = append(quads, r)
		case 3:
			trips = append(trips, r)
		case 2:
			pairs = append(pairs, r)
		}
	}

	// Four of a kind
	if len(quads) == 1 {
		kicker := []int{}
		for r, n := range rankCount {
			if n == 1 {
				kicker = append(kicker, r)
			}
		}
		return FinalRank{Rank: FourOfAKind, Kicker: append([]int{quads[0]}, kicker...)}
	}
	// Full house
	if len(trips) == 1 && len(pairs) == 1 {
		return FinalRank{Rank: FullHouse, Kicker: []int{trips[0], pairs[0]}}
	}
	// Flush
	if isFlush {
		return FinalRank{Rank: Flush, Kicker: reverseSort(ranks)}
	}
	// Straight
	if isStraight {
		return FinalRank{Rank: Straight, Kicker: ranks}
	}
	// Three of a kind
	if len(trips) == 1 {
		kickers := []int{}
		for r, n := range rankCount {
			if n == 1 {
				kickers = append(kickers, r)
			}
		}
		return FinalRank{Rank: ThreeOfAKind, Kicker: append([]int{trips[0]}, reverseSort(kickers)...)}
	}
	// Two pair
	if len(pairs) == 2 {
		sort.Ints(pairs)
		high := pairs[1]
		low := pairs[0]
		kicker := 0
		for r, n := range rankCount {
			if n == 1 {
				kicker = r
			}
		}
		return FinalRank{Rank: TwoPair, Kicker: []int{high, low, kicker}}
	}

	// One pair (any rank)
	if len(pairs) == 1 {
		pairRank := pairs[0]
		kickers := []int{}
		for r, n := range rankCount {
			if n == 1 {
				kickers = append(kickers, r)
			}
		}
		return FinalRank{Rank: OnePair, Kicker: append([]int{pairRank}, reverseSort(kickers)...)}
	}

	// High card
	return FinalRank{Rank: HighCard, Kicker: reverseSort(ranks)}
}

func reverseSort(a []int) []int {
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	return a
}
