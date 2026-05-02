package common

import "sort"

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

/*

Kicker[] encodes the exact tie‑break signature for the hand.
The slice is always ordered high → low and its meaning depends on the hand type:

  HighCard:     [card1, card2, card3, card4, card5]
  OnePair:      [pairRank, kicker1, kicker2, kicker3]
  TwoPair:      [highPairRank, lowPairRank, kicker]
  Trips:        [tripRank, kicker1, kicker2]
  Straight:     [highCard]  // 5-high straight encoded as 5
  Flush:        [card1, card2, card3, card4, card5]
  FullHouse:    [tripRank, pairRank]
  Quads:        [quadRank, kicker]
  StraightFlush:[highCard]

CompareFinal uses comparison of this slice to break ties.

These games use only 1 deck shuffled @ each deal.
So, duplicate trips and quads aren't possible.
It's only for future (?) enhancments.

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
		return "High Card"
	default:
		return "Unknown"
	}
}

func Best5of7(cards []Card) FinalRank {
	best := FinalRank{} // zero = worst hand

	var five [5]Card

	// choose i < j < k < l < m
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			for k := j + 1; k < 5; k++ {
				for l := k + 1; l < 6; l++ {
					for m := l + 1; m < 7; m++ {

						five[0] = cards[i]
						five[1] = cards[j]
						five[2] = cards[k]
						five[3] = cards[l]
						five[4] = cards[m]

						rank := evaluateFinal(five[:])

						if rank.betterThan(best) {
							best = rank
						}
					}
				}
			}
		}
	}

	return best
}

// > 0 player wins
// < 0 dealer wins
// 0 == tie (push)
func CompareFinal(player, dealer FinalRank) int {
	// Compare category first
	if player.Rank != dealer.Rank {
		return int(player.Rank) - int(dealer.Rank)
	}

	// Compare kickers lexicographically
	lp := len(player.Kicker)
	ld := len(dealer.Kicker)
	n := lp
	if ld < n {
		n = ld
	}

	for i := 0; i < n; i++ {
		if player.Kicker[i] != dealer.Kicker[i] {
			return player.Kicker[i] - dealer.Kicker[i]
		}
	}

	// If all compared kickers equal, longer kicker list wins
	if lp != ld {
		return lp - ld
	}

	return 0
}

func evaluateFinal(cards []Card) FinalRank {
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

	// Normal straight (all ranks unique and span 4)
	if ranks[4]-ranks[0] == 4 && len(rankCount) == 5 {
		isStraight = true
	}

	// Ace-low straight (A 2 3 4 5)
	if ranks[0] == 2 && ranks[1] == 3 && ranks[2] == 4 && ranks[3] == 5 && ranks[4] == 14 {
		isStraight = true
		ranks = []int{5, 4, 3, 2, 1}
	}

	// Flush detection
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

	// Count pairs/trips/quads
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

	// One pair
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

func (a FinalRank) betterThan(b FinalRank) bool {
	// Compare rank category first
	if a.Rank != b.Rank {
		return a.Rank > b.Rank
	}

	// Compare kickers lexicographically
	la := len(a.Kicker)
	lb := len(b.Kicker)
	n := la
	if lb < n {
		n = lb
	}

	for i := 0; i < n; i++ {
		if a.Kicker[i] != b.Kicker[i] {
			return a.Kicker[i] > b.Kicker[i]
		}
	}

	// If all compared kickers equal, longer kicker list wins
	if la != lb {
		return la > lb
	}

	// Completely equal
	return false
}
