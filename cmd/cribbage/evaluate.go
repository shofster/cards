package main

import (
	"cards/common"
	"log"
)

/*

  File:    evaluate.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: CoPilot AI player implementation.

As dealer, AI happily throws 5s, pairs, suited cards, run cards into the crib.
As non‑dealer, AI avoids feeding those like poison.
AI keeps strong 4‑card cores (15s, pairs, runs).
AI avoids “trap” hands that look good but have low expected value.

*/

func scorePegging(who int, card common.Card, pegged []common.Card, total int) int {
	points := 0
	reason := ""
	switch total {
	case 15:
		if Ctx.debug {
			log.Printf("scorePegging-15: who %d total %d scored 2\n",
				who, total)
		}
		points += 2
		reason = "15"
	case 31:
		if Ctx.debug {
			log.Printf("scorePegging-31: who %d total %d scored 2\n",
				who, total)
		}
		points += 2
		reason = "31"
	}
	pp := pairPoints(card, pegged)
	if pp > 0 { // 2 or more
		points += pp
		if Ctx.debug {
			log.Printf("scorePegging-pairs: who %d total %d scored %d\n",
				who, total, pp)
		}
		reason = "pair"
	}
	pp = runPoints(pegged)
	if pp > 0 { // 3 or more
		points += pp
		if Ctx.debug {
			log.Printf("scorePegging-run: who %d total %d scored %d\n",
				who, total, pp)
		}
		reason = "run"
	}
	// apply points to "who"
	if points > 0 {
		scorePoints(who, points, reason)
	}
	return points
}

// chooseDiscard selects the best 4‑card keep and 2‑card discard
// using expected hand value + crib bias.
// Input: 6 cards
// Output: keep (4 cards), discard (2 cards)
func chooseDiscard(cards common.Cards, isDealer bool) (keep common.Cards, discard common.Cards) {
	if len(cards) != 6 {
		panic("chooseDiscard: needs exactly 6 cards")
	}

	bestScore := -1e9 // effectively -infinity

	// Enumerate all 15 ways to choose 2 cards to discard
	for i := 0; i < 6; i++ {
		for j := i + 1; j < 6; j++ {

			// Build keep (4 cards)
			keepTmp := make(common.Cards, 0, 4)
			for k := 0; k < 6; k++ {
				if k != i && k != j {
					keepTmp = append(keepTmp, cards[k])
				}
			}

			// Build discard (2 cards)
			discardTmp := common.Cards{cards[i], cards[j]}

			// Expected hand value
			ev := expectedHandScore(keepTmp)

			// Crib bias
			bias := cribBias(discardTmp, isDealer)

			total := ev + float64(bias)

			if total > bestScore {
				bestScore = total
				keep = append([]common.Card{}, keepTmp...)
				discard = append([]common.Card{}, discardTmp...)
			}
		}
	}

	return keep, discard
}

// expectedHandScore returns the average score of a 4‑card hand
// over all 46 possible starter cards.
// keep must contain exactly 4 cards.
//
//	used to determine discards after deal
func expectedHandScore(keep common.Cards) float64 {
	total := 0
	count := 0

	for _, code := range cardDeck {

		starter := Ctx.cardMap[code]
		// skip cards already in the hand
		if common.Contains(keep, starter) {
			continue
		}

		// score the 5‑card hand
		total += scoreHand(keep, starter)
		count++
	}

	return float64(total) / float64(count)
}

// scoreHand scores a 4 card hand (PLAYER, AI, or CRIB) + starter.
func scoreHand(hand common.Cards, starter common.Card) int {
	cards := append([]common.Card{}, hand...)
	cards = append(cards, starter)

	score := 0

	score15 := score15s(cards)
	scorePair := scorePairs(cards)
	scoreRun := scoreRuns(cards)
	scoreFl := scoreFlush(hand, starter)
	scoreNob := scoreNobs(hand, starter)
	score += score15
	score += scorePair
	score += scoreRun
	score += scoreFl
	score += scoreNob

	return score
}

func isLow(c common.Card) bool {
	return c.Rank == 2 || c.Rank == 3 || c.Rank == 4
}

func isTenValue(c common.Card) bool {
	return c.Rank >= 10 // 10, J, Q, K
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func score15s(cards common.Cards) int {
	n := len(cards)
	score := 0

	// All subsets except empty
	for mask := 1; mask < (1 << n); mask++ {
		sum := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sum += cribValue(cards[i]) // A=1, 2=2, ... 10/J/Q/K=10
			}
		}
		if sum == 15 {
			score += 2
		}
	}
	return score
}

func cribValue(c common.Card) int {
	if c.Rank >= 10 {
		return 10
	}
	return int(c.Rank)
}

func scorePairs(cards common.Cards) int {
	score := 0
	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if cards[i].Rank == cards[j].Rank {
				score += 2
			}
		}
	}
	return score
}

func scoreRuns(cards common.Cards) int {
	n := len(cards)
	if n < 3 {
		return 0
	}

	// Try all subsets of size 5,4,3
	for size := 5; size >= 3; size-- {
		count := 0
		for _, subset := range common.SubsetsOf(cards, size) {
			if common.IsRun(subset) {
				count++
			}
		}
		if count > 0 {
			return count * size
		}
	}
	return 0
}

func scoreFlush(hand common.Cards, starter common.Card) int {
	// hand = 4 cards only
	suit := hand[0].Suit
	for _, c := range hand[1:] {
		if c.Suit != suit {
			return 0
		}
	}
	// 4‑card flush
	score := 4
	// 5‑card flush?
	if starter.Suit == suit {
		score = 5
	}
	return score
}

func scoreNobs(hand common.Cards, starter common.Card) int {
	for _, c := range hand {
		if c.Rank == common.Jack && c.Suit == starter.Suit {
			return 1
		}
	}
	return 0
}

// cribBias returns a score for the two discarded cards.
// Positive = good for AI's crib (isDealer = true)
// Negative = good for opponent's crib (isDealer = false)
func cribBias(discard common.Cards, isDealer bool) int {
	score := 0

	c1, c2 := discard[0], discard[1]

	// 5s are strong, but less dominant
	if c1.Rank == 5 {
		score += 3
	}
	if c2.Rank == 5 {
		score += 3
	}

	// Pair
	if c1.Rank == c2.Rank {
		score += 3
	}

	// Suited cards
	if c1.Suit == c2.Suit {
		score += 2
	}

	// Run-connected (difference of 1)
	if abs(int(c1.Rank)-int(c2.Rank)) == 1 {
		score += 2
	}

	// Low cards (2,3,4)
	if isLow(c1) {
		score += 2
	}
	if isLow(c2) {
		score += 2
	}

	// Tens and face cards (good crib filler)
	if isTenValue(c1) {
		score += 1
	}
	if isTenValue(c2) {
		score += 1
	}

	// If crib belongs to opponent, invert the score
	if !isDealer {
		score = -score
	}

	return score
}
