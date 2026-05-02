package main

import (
	"games/common"
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

func RankValue(card common.Card) int {
	switch card.Rank {
	case common.Jack, common.Queen, common.King:
		return 10
	case common.Ace:
		return 1
	default:
		return int(card.Rank)
	}
}

func AceLow(r common.Rank) int {
	if r == common.Ace {
		return 1
	}
	return int(r)
}

func shouldTakeDiscard(discard common.Card, hand []common.Card) bool {

	// compute the hand with discard
	withDiscard := append(copyCards(hand), discard)
	current := computeMelds(hand)
	after := computeMelds(withDiscard)
	//fmt.Printf("ShouldTakeDiscard: (%s) current %d, afeter %d\n", discard.Icon,
	//	current.DeadwoodPoints, after.DeadwoodPoints)

	// discard if deadwood is reduced by threshold
	if current.DeadwoodPoints-after.DeadwoodPoints >= 4 {
		return true
	}

	// completes a meld (deadwood goes to 0 or big drop)
	if after.DeadwoodPoints == 0 {
		return true
	}

	// OR - creates near-meld
	if createsNearMeld(discard, hand) {
		return true
	}

	return false
}

// createsNearMeld: same rank present OR adjacent ranks same suit
func createsNearMeld(card common.Card, hand []common.Card) bool {
	// same-rank potential
	sameRank := 0
	for _, h := range hand {
		if h.Rank == card.Rank {
			sameRank++
		}
	}
	if sameRank >= 2 {
		return true
	}
	// run potential with Ace = 1
	rv := AceLow(card.Rank)
	for _, h := range hand {
		if h.Suit != card.Suit {
			continue
		}
		hv := AceLow(h.Rank)
		if hv == rv-1 || hv == rv+1 {
			return true
		}
	}
	return false
}

func chooseDiscard(hand []common.Card) common.Card {
	bestScore := 1_000_000
	var best common.Card

	for _, c := range hand {
		score := discardScore(c, hand)
		if score < bestScore {
			bestScore = score
			best = c
		}
	}
	return best
}

func discardScore(card common.Card, hand []common.Card) int {
	// weights
	const (
		wMeld     = 5
		wNear     = 3
		wDeadwood = 1
	)

	meldPot := meldPotential(card, hand)
	nearPot := nearMeldPotential(card, hand)
	dead := RankValue(card)

	return wMeld*meldPot + wNear*nearPot + wDeadwood*dead
}

func meldPotential(card common.Card, hand []common.Card) int {
	// crude: +3 if part of any run/set in best solution
	res := computeMelds(hand)
	for _, m := range res.Melds {
		for _, mc := range m.Cards {
			if sameCard(mc, card) {
				return 3
			}
		}
	}
	return 0
}

func nearMeldPotential(card common.Card, hand []common.Card) int {
	score := 0
	for _, h := range hand {
		if sameCard(h, card) {
			continue
		}
		// same rank (set potential)
		if h.Rank == card.Rank {
			score++
		}
		// same suit, close rank (run potential)
		if h.Suit == card.Suit {
			d := AceLow(h.Rank) - AceLow(card.Rank)
			if d == 1 || d == -1 {
				score++
			}
		}
	}
	return score
}

type Result struct {
	KnockerPoints int
	KnockeePoints int
	Diff          int
	Gin           bool
	KnockedOut    bool // KO'ed
}

func resolveKnock(knockerHand, knockeeHand []common.Card) Result {
	km := computeMelds(knockerHand)
	nm := computeMelds(knockeeHand)

	// GIN == knocker has 0 deadwood -- layoff not allowed
	if km.DeadwoodPoints == 0 {
		knockeePoints := sumPoints(nm.Deadwood)
		return scoreGin(knockeePoints)
	}

	// Normal knock: knockee may lay off onto knocker's melds
	_, remaining := autoLayoff(nm.Deadwood, km.Melds)
	knockeePoints := sumPoints(remaining)

	return scoreKnock(km.DeadwoodPoints, knockeePoints)
}

func sumPoints(cards []common.Card) int {
	total := 0
	for _, c := range cards {
		total += RankValue(c)
	}
	return total
}

// knocker wins by all knockee deadwood
func scoreGin(knockeePoints int) Result {
	return Result{
		KnockerPoints: 0,
		KnockeePoints: knockeePoints,
		Diff:          knockeePoints,
		Gin:           true,
		KnockedOut:    false,
	}
}

// winner may be knockee, if lay-offs applied
func scoreKnock(knockerPoints, knockeePoints int) Result {

	// undercut: knockee wins
	if knockeePoints <= knockerPoints {
		return Result{
			KnockerPoints: knockerPoints,
			KnockeePoints: knockeePoints,
			Diff:          knockeePoints - knockerPoints,
			Gin:           false,
			KnockedOut:    true,
		}
	}

	// not KO'ed : knocker wins
	return Result{
		KnockerPoints: knockerPoints,
		KnockeePoints: knockeePoints,
		Diff:          knockeePoints - knockerPoints,
		Gin:           false,
		KnockedOut:    false,
	}
}

func autoLayoff(deadwood []common.Card, knockerMelds []Meld) (played []common.Card, remaining []common.Card) {
	for _, c := range deadwood {
		if canPlayOnAnyMeld(c, knockerMelds) {
			played = append(played, c)
		} else {
			remaining = append(remaining, c)
		}
	}
	return
}

func canPlayOnAnyMeld(c common.Card, melds []Meld) bool {
	for _, m := range melds {
		switch m.Type {
		case MeldSet:
			if fitsSet(c, m.Cards) {
				return true
			}
		case MeldRun:
			if fitsRun(c, m.Cards) {
				return true
			}
		}
	}
	return false
}

func fitsSet(c common.Card, set []common.Card) bool {
	return c.Rank == set[0].Rank
}

func fitsRun(c common.Card, run []common.Card) bool {
	if c.Suit != run[0].Suit {
		return false
	}
	// low side extension
	if AceLow(c.Rank) == AceLow(run[0].Rank)-1 {
		return true
	}
	// upper side extension
	last := run[len(run)-1]
	if AceLow(c.Rank) == AceLow(last.Rank)+1 {
		return true
	}
	return false
}
