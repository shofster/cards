package main

import (
	"cards/common"
	"log"
	"sort"
)

/*

  File:    pegging.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:  Adapted from CoPilot AI !
avoids giving up runs
avoids giving up 15s
avoids giving up pairs
aggressively scores when safe
plays defensively when behind
plays offensively when ahead
understands end‑of‑sequence traps
understands 31 traps
understands count‑control
understands card‑tracking (what ranks remain)
*/

func choosePegCard(hand common.Cards, played []common.Card, total int) (common.Card, bool) {
	legal := legalPegCards(hand, total)
	if len(legal) == 0 {
		return common.Card{}, false // "Go"
	}

	// 1. Score immediate points (15, 31, pair, run)
	if c, ok := pegImmediateScore(legal, played, total); ok {
		debugPegChoice(hand, legal, played, total, c)
		return c, true
	}

	// 2. Avoid giving opponent points (trap avoidance)
	safe := pegSafeCards(legal, played, total)
	if len(safe) > 0 {
		chosen := chooseLowest(safe)
		debugPegChoice(hand, legal, played, total, chosen)
		return chosen, true
	}

	// 3. No safe cards —> choose the least dangerous
	chosen := chooseLeastBad(legal, played, total, hand)
	debugPegChoice(hand, legal, played, total, chosen)
	return chosen, true
}

func legalPegCards(hand common.Cards, total int) []common.Card {
	var out []common.Card
	for _, c := range hand {
		if c.Selected {
			continue // ignore cards already played
		}
		if total+c.PegValue() <= 31 {
			out = append(out, c)
		}
	}
	return out
}

func chooseLowest(cards common.Cards) common.Card {
	best := cards[0]
	for _, c := range cards {
		if c.PegValue() < best.PegValue() {
			best = c
		}
	}
	return best
}

// pegImmediateScore - tries every legal card,
// computes pegPoints(c, played, total),
// and picks the card with the highest score
func pegImmediateScore(legal, played common.Cards, total int) (common.Card, bool) {
	best := common.Card{}
	bestPts := 0

	for _, c := range legal {
		pts := pegPoints(c, played, total)
		if pts > bestPts {
			bestPts = pts
			best = c
		}
	}

	if bestPts > 0 {
		return best, true
	}
	return common.Card{}, false
}

func pegPoints(c common.Card, played common.Cards, total int) int {
	pts := 0
	v := c.PegValue()

	// 15 or 31
	if total+v == 15 {
		pts += 2
	}
	if total+v == 31 {
		pts += 2
	}

	// Pair / pair royal / 4 of a kind
	if len(played) > 0 {
		last := played[len(played)-1]
		if last.Rank == c.Rank {
			pts += 2
			if len(played) >= 2 && played[len(played)-2].Rank == c.Rank {
				pts += 4
				if len(played) >= 3 && played[len(played)-3].Rank == c.Rank {
					pts += 6
				}
			}
		}
	}

	// Runs
	pts += runPoints(append(played, c))

	return pts
}

func pairPoints(card common.Card, pegged []common.Card) int {
	n := len(pegged)
	if n < 2 {
		return 0
	}
	// Count backwards from the card BEFORE the new one
	matches := 0
	for i := n - 2; i >= 0; i-- {
		if pegged[i].Rank == card.Rank {
			matches++
		} else {
			break
		}
	}
	switch matches {
	case 1:
		return 2 // pair
	case 2:
		return 6 // three of a kind
	case 3:
		return 12 // four of a kind
	}
	return 0
}

// runPoints - find the longest valid run in the sequence.
func runPoints(seq common.Cards) int {
	n := len(seq)
	if n < 3 {
		return 0
	}

	// Check backwards for longest run
	for size := n; size >= 3; size-- {
		sub := seq[n-size:]
		if common.IsRun(sub) {
			return size
		}
	}
	return 0
}

// pegSafeCards - any card that doesn’t set up opponent‑points.
func pegSafeCards(legal, played common.Cards, total int) common.Cards {
	safe := common.Cards{}
	for _, c := range legal {
		if !wouldGiveOpponentPoints(c, played, total) {
			safe = append(safe, c)
		}
	}
	return safe
}

func wouldGiveOpponentPoints(card common.Card, played common.Cards, total int) bool {
	newTotal := total + card.PegValue()

	// Avoid setting up 15 or 31
	if newTotal == 14 || newTotal == 30 {
		return true
	}

	// Avoid setting up pairs
	if len(played) > 0 && played[len(played)-1].Rank == card.Rank {
		return true
	}

	// Avoid setting up runs
	if createsRunOpportunity(card, played) {
		return true
	}

	return false
}

func createsRunOpportunity(card common.Card, played common.Cards) bool {
	// Opponent plays after we play card
	seq := append(played, card)
	n := len(seq)

	// Opponent can complete a run of 3, 4, or 5
	// So check windows of size 2,3,4 (opponent's card makes size+1)
	for window := 2; window <= 4; window++ {
		if n < window {
			continue
		}

		tail := seq[n-window:] // last N cards before opponent plays

		// Try all possible ranks the opponent could have
		for r := common.Rank(1); r <= common.Rank(13); r++ {
			test := append(tail, common.Card{Rank: r})
			if common.IsRun(test) {
				return true
			}
		}
	}

	return false
}

func chooseLeastBad(legal, played common.Cards, total int, hand common.Cards) common.Card {
	best := legal[0]
	bestDanger := 999

	for _, c := range legal {
		d := dangerScore(c, played, total, hand)
		if d < bestDanger {
			bestDanger = d
			best = c
		}
	}

	return best
}

// showPenalty - penalty if you break a potential pair.
func showPenalty(hand common.Cards, c common.Card) int {
	// Don’t break pairs
	count := 0
	for _, x := range hand {
		if x.Rank == c.Rank {
			count++
		}
	}
	if count >= 2 {
		return 2
	}

	// Don’t break 5-5-6-7, 6-7-8, 3-4-5-6 structures
	// (avoid breaking 4-card runs)
	ranks := []int{}
	for _, x := range hand {
		if x.Icon != c.Icon {
			ranks = append(ranks, int(x.Rank))
		}
	}
	sort.Ints(ranks)
	for i := 2; i < len(ranks); i++ {
		if ranks[i] == ranks[i-1]+1 && ranks[i-1] == ranks[i-2]+1 {
			return 1
		}
	}

	return 0
}

func positionalPenalty(card common.Card, total int) int {
	// Low cards early = good
	if total < 10 && card.PegValue() <= 3 {
		return 0
	}

	// High cards early = bad
	if total < 10 && card.PegValue() >= 10 {
		return 1
	}

	// High cards late = good (no penalty)
	if total >= 20 && card.PegValue() >= 10 {
		return 0
	}

	return 0
}

func dangerScore(card common.Card, played common.Cards, total int, hand common.Cards) int {
	return pairDanger(card, played) +
		runDanger(card, played) +
		fifteenDanger(card, total) +
		countControlDanger(card, total) +
		showPenalty(hand, card) +
		positionalPenalty(card, total)
}

func pairDanger(card common.Card, played common.Cards) int {
	if len(played) == 0 {
		return 0
	}
	last := played[len(played)-1]
	if last.Rank == card.Rank {
		return 3 // pairs are dangerous
	}
	return 0
}
func runDanger(card common.Card, played common.Cards) int {
	if createsRunOpportunity(card, played) {
		return 4 // runs are VERY dangerous
	}
	return 0
}
func fifteenDanger(card common.Card, total int) int {
	newTotal := total + card.PegValue()

	if newTotal == 14 {
		return 3
	} // opponent hits 15
	if newTotal == 30 {
		return 3
	} // opponent hits 31

	return 0
}
func countControlDanger(card common.Card, total int) int {
	newTotal := total + card.PegValue()

	// Avoid giving opponent easy 15s
	if newTotal == 5 || newTotal == 20 {
		return 2
	}

	// Avoid giving opponent easy 31s
	if newTotal == 21 {
		return 2
	}

	// Avoid playing 5 early
	if total <= 10 && card.PegValue() == 5 {
		return 2
	}

	return 0
}

/*
	Debugging
*/

//var PeggingDebug = true

func debugPegChoice(_, legal, played common.Cards, total int, chosen common.Card) {
	if !Ctx.debug {
		return
	}

	log.Println("PEGGING DEBUG:")
	log.Printf("  played: %s\n", played)
	log.Printf("  total: %d\n", total)
	log.Printf("  legal: %s\n", legal)

	/*	for _, c := range legal {
			pts := pegPoints(c, played, total)
			danger := dangerScore(c, played, total, hand)
			log.Printf("  card %-3s: points=%d danger=%d\n", c.Icon, pts, danger)
		}
	*/
	log.Printf("CHOSEN: %s\n\n", chosen.Icon)
}

/* from perplexity - less aggresive --
// dangerScore measures how much playing 'card' risks giving the opponent points
func dangerScore(card common.Card, played common.Cards, total int, hand common.Cards) int {
    return pairDanger(card, played) +
        runDanger(card, played) +
        fifteenDanger(card, total) +
        countControlDanger(card, total) +
        showPenalty(hand, card) +
        positionalPenalty(card, total)
}

// pairDanger: small penalty for playing into a possible pair
func pairDanger(card common.Card, played common.Cards) int {
    if len(played) == 0 {
        return 0
    }
    last := played[len(played)-1]
    if last.Rank == card.Rank {
        return 2 // milder than before, but still undesirable
    }
    return 0
}

// runDanger: reduced penalty for adjacency‑based run risks
func runDanger(card common.Card, played common.Cards) int {
    if createsRunOpportunity(card, played) {
        return 2 // was 4; now only a moderate risk
    }
    return 0
}

// fifteenDanger: big penalty for 14/30
func fifteenDanger(card common.Card, total int) int {
    newTotal := total + card.PegValue()

    if newTotal == 14 {
        return 4 // 15 trap is high penalty
    }
    if newTotal == 30 {
        return 4 // 31 trap is high penalty
    }
    return 0
}

// countControlDanger: medium penalties for 5/20/21 and early 5s
func countControlDanger(card common.Card, total int) int {
    newTotal := total + card.PegValue()

    if newTotal == 5 || newTotal == 20 {
        return 2 // easy 15 for opponent
    }
    if newTotal == 21 {
        return 2 // easy 31 for opponent
    }
    if total <= 10 && card.PegValue() == 5 {
        return 2 // avoid 5‑early unless you know it’s safe
    }
    return 0
}

// createsRunOpportunity is now only concerned with "nearby" ranks that can easily extend runs
func createsRunOpportunity(card common.Card, played common.Cards) bool {
    if len(played) == 0 {
        return false
    }
    seq := append(played, card)
    n := len(seq)

    // Only check if the last 2–3 cards create a “near‑run” window
    // (i.e., 2‑cards forming a partial run or 2‑gap that a neighbor could fill).
    // Instead of trying all 13 ranks, only check rank‑neighbors.

    // Let’s look at windows of size 2 and 3.
    for window := 2; window <= 3 && window < n; window++ {
        tail := seq[n-window:]

        // Collect relative offsets and check if they’re close to a run.
        // Instead of testing all 13 ranks, we only test ranks adjacent to existing ones.
        for _, base := range tail {
            for _, offset := range []int{-1, 0, 1} {
                r := int(base.Rank) + offset
                if r < 1 || r > 13 {
                    continue
                }
                test := append(tail, common.Card{Rank: common.Rank(r)})
                if common.IsRun(test) {
                    return true
                }
            }
        }
    }

    return false
}

*/
