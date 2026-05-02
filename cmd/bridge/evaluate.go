package main

import (
	"games/common"
	"strconv"
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

type Player int

const (
	South = iota
	West
	North
	East
)

func playerToPosition(player Player) string {
	switch player {
	case South:
		return "South"
	case West:
		return "West"
	case North:
		return "North"
	case East:
		return "East"
	}
	panic("playerToPosition: invalid Player")
}

func suitToSuit(suit string) common.Suit {
	switch suit {
	case "♣":
		return common.Clubs
	case "♦":
		return common.Diamonds
	case "♥":
		return common.Hearts
	case "♠":
		return common.Spades
	default:
		return common.NoTrump
	}
}

// determine the winnner of the played trick
func winnerOfTrick(cards [4]common.Card, trump common.Suit, leader Player) Player {

	// suit of the first card.
	leadSuit := cards[leader].Suit

	// collect trump cards played.
	trumpIndexes := []Player{}
	if trump != common.NoTrump {
		for p := Player(South); p <= Player(East); p++ {
			if cards[p].Suit == trump {
				trumpIndexes = append(trumpIndexes, p)
			}
		}
	}

	// highest trump card wins.
	if len(trumpIndexes) > 0 {
		winner := trumpIndexes[0]
		for _, p := range trumpIndexes[1:] {
			if cards[p].Rank > cards[winner].Rank {
				winner = p
			}
		}
		return winner
	}

	// highest card in the leadSuit wins.
	winner := leader
	for p := Player(South); p <= Player(East); p++ {
		if cards[p].Suit == leadSuit && cards[p].Rank > cards[winner].Rank {
			winner = p
		}
	}

	return winner
}

// dealScores returns:
//
//	deltaBid     = tricksWon - (6 + level)
//	deltaAbove   = NS above-the-line delta
//	deltaBelow   = NS below-the-line delta
//
// South is always declarer.
// winners = 13 strings: "N", "S", "E", "W"
func dealScores(winners []string, levelStr string, suitStr string) (deltaBid int, deltaAbove int, deltaBelow int) {

	level, _ := strconv.Atoi(levelStr)
	suit := suitToSuit(suitStr)

	// Count NS tricks
	ns := 0
	for _, w := range winners {
		if w == "N" || w == "S" {
			ns++
		}
	}

	// South always declarer → NS is declaring side
	tricksWon := ns

	// Required tricks
	required := 6 + level
	deltaBid = tricksWon - required // positive = made, negative = down

	// Trick points (below the line)
	trickPoints := 0
	switch suit {
	case common.NoTrump:
		trickPoints = 40 + (level-1)*30
	case common.Spades, common.Hearts:
		trickPoints = level * 30
	case common.Diamonds, common.Clubs:
		trickPoints = level * 20
	}

	// Overtrick values
	overValue := 0
	switch suit {
	case common.NoTrump, common.Spades, common.Hearts:
		overValue = 30
	default:
		overValue = 20
	}

	// Contract made
	if deltaBid >= 0 {
		// Below-the-line trick points
		deltaBelow = trickPoints

		// Overtricks
		deltaAbove += deltaBid * overValue

		// Part-score or game bonus
		if trickPoints >= 100 {
			deltaAbove += 300
		} else {
			deltaAbove += 50
		}

		// Slam bonuses
		if level == 6 {
			deltaAbove += 500
		}
		if level == 7 {
			deltaAbove += 1000
		}

		return deltaBid, deltaAbove, deltaBelow
	}

	// Contract down — NSAboveDelta is negative
	down := -deltaBid
	deltaAbove = -(down * 50)

	return deltaBid, deltaAbove, 0
}
