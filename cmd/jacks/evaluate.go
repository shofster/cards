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
A flush or straight on the initial deal is not announced in real
 Jacks or Better.
Only a high pair (Jacks, Queens, Kings, Aces) is highlighted on
 the deal screen.
Everything else—flush, straight, full house, trips,
 two pair—is intentionally silent until after the draw.
This matches real casino Jacks or Better behavior:
the machine only alerts the player to a single high pair because
that’s the only deal‑phase hand that affects strategy.

🧩 Why trips (even high trips) must NOT announce
Trips (e.g., J‑J‑J‑x‑x) are a made hand,
  but casinos intentionally do not reveal that on the deal.
They want the player to focus on holds, not early wins.
Trips are only evaluated and announced after the draw.

Outcomes:
- A♠ A♦ 7♣ 4♥ 2♠ → announce (exactly one high pair)
- K♣ K♦ K♥ 9♠ 3♦ → no announcement (trips)
- Q♠ Q♦ J♣ J♦ 5♥ → no announcement (two pair)
- J♣ 8♦ 7♠ 6♥ 2♣ → no announcement (no pair)

*/

func DealAnnouncement(hand []common.Card) common.Rank {

	// count the number of same (HIGH) Rank
	rankCount := map[common.Rank]int{}
	for _, c := range hand {
		switch c.Rank {
		case common.Jack, common.Queen, common.King, common.Ace:
			rankCount[c.Rank]++
		}
	}
	// find first high pair
	var highPair common.Rank = -1
	for r, n := range rankCount {
		if n == 2 {
			highPair = r
			break
		}
	}
	return highPair
}

type HandRank int

const (
	NoWin HandRank = iota
	JacksOrBetter
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
	case JacksOrBetter:
		return "Jacks or Better"
	case NoWin:
		return "No Win"
	default:
		return "Unknown"
	}
}

var PayTables = map[string]map[HandRank][]int{
	"9/6": {
		RoyalFlush:    {250, 500, 750, 1000, 4000},
		StraightFlush: {50, 100, 150, 200, 250},
		FourOfAKind:   {25, 50, 75, 100, 125},
		FullHouse:     {9, 18, 27, 36, 45},
		Flush:         {6, 12, 18, 24, 30},
		Straight:      {4, 8, 12, 16, 20},
		ThreeOfAKind:  {3, 6, 9, 12, 15},
		TwoPair:       {2, 4, 6, 8, 10},
		JacksOrBetter: {1, 2, 3, 4, 5},
		NoWin:         {0, 0, 0, 0, 0},
	},

	"8/5": {
		RoyalFlush:    {250, 500, 750, 1000, 4000},
		StraightFlush: {50, 100, 150, 200, 250},
		FourOfAKind:   {25, 50, 75, 100, 125},
		FullHouse:     {8, 16, 24, 32, 40}, // 8/5 change
		Flush:         {5, 10, 15, 20, 25}, // 8/5 change
		Straight:      {4, 8, 12, 16, 20},
		ThreeOfAKind:  {3, 6, 9, 12, 15},
		TwoPair:       {2, 4, 6, 8, 10},
		JacksOrBetter: {1, 2, 3, 4, 5},
		NoWin:         {0, 0, 0, 0, 0},
	},
}

var payTables = []string{"9/6", "8/5"}

func Payout(table string, rank HandRank, bet int) int {
	if bet < 1 || bet > 5 {
		return 0
	}

	t, ok := PayTables[table]
	if !ok {
		return 0
	}
	payouts, ok := t[rank]
	if !ok {
		return 0
	}
	return payouts[bet-1]
}

func EvaluateFinal(hand []common.Card) HandRank {
	// count ranks
	rankCount := map[common.Rank]int{}
	for _, c := range hand {
		rankCount[c.Rank]++
	}

	flush := isFlush(hand)
	r := sortedRanks(hand)
	straight := isStraight(r)

	if straight && flush {
		if isRoyal(r) {
			return RoyalFlush
		}
		return StraightFlush
	}

	if hasN(rankCount, 4) {
		return FourOfAKind
	}
	if hasFullHouse(rankCount) {
		return FullHouse
	}
	if flush {
		return Flush
	}
	if straight {
		return Straight
	}
	if hasN(rankCount, 3) {
		return ThreeOfAKind
	}
	if hasTwoPair(rankCount) {
		return TwoPair
	}
	if hasHighPair(rankCount) {
		return JacksOrBetter
	}

	return NoWin
}

func sortedRanks(hand []common.Card) []int {
	r := make([]int, len(hand))
	for i, c := range hand {
		r[i] = int(c.Rank)
	}
	sort.Ints(r)
	return r
}

func hasN(rankCount map[common.Rank]int, n int) bool {
	for _, c := range rankCount {
		if c == n {
			return true
		}
	}
	return false
}

func hasTwoPair(rankCount map[common.Rank]int) bool {
	pairs := 0
	for _, c := range rankCount {
		if c == 2 {
			pairs++
		}
	}
	return pairs == 2
}

func hasFullHouse(rankCount map[common.Rank]int) bool {
	has3 := false
	has2 := false
	for _, c := range rankCount {
		if c == 3 {
			has3 = true
		} else if c == 2 {
			has2 = true
		}
	}
	return has3 && has2
}

func hasHighPair(rankCount map[common.Rank]int) bool {
	for r, c := range rankCount {
		if c == 2 && r > common.Ten { // Jack, Queen, King, Ace
			return true
		}
	}
	return false
}

func isFlush(hand []common.Card) bool {
	s := hand[0].Suit
	for _, c := range hand {
		if c.Suit != s {
			return false
		}
	}
	return true
}

func isStraight(r []int) bool {
	// A-2-3-4-5 wheel
	if r[0] == int(common.Two) &&
		r[1] == int(common.Three) &&
		r[2] == int(common.Four) &&
		r[3] == int(common.Five) &&
		r[4] == int(common.Ace) {
		return true
	}

	// normal straight
	for i := 0; i < 4; i++ {
		if r[i+1] != r[i]+1 {
			return false
		}
	}
	return true
}

func isRoyal(r []int) bool {
	return r[0] == int(common.Ten) &&
		r[1] == int(common.Jack) &&
		r[2] == int(common.Queen) &&
		r[3] == int(common.King) &&
		r[4] == int(common.Ace)
}
