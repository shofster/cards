package common

import (
	"math/rand/v2"
	"sort"
	"strconv"
	"time"
)

/*

  File:    cutDeck.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: provide a graphical cut for deal.
*/

func CutForDeal(players []string, lowest bool) []string {

	type cut struct {
		player string
		icon   string
		rank   int
		suit   int
	}

	// Parse rank from icon strings ("A♠", "10♦", "J♥", etc.)
	parseRank := func(icon string) int {
		// Suit is always the last rune (3 bytes in UTF-8)
		// Rank is everything before that
		r := icon[:len(icon)-3]
		switch r {
		case "A":
			if lowest {
				return 1
			}
			return 14
		case "K":
			return 13
		case "Q":
			return 12
		case "J":
			return 11
		default:
			n, _ := strconv.Atoi(r)
			return n
		}
	}
	// Suit parser: ♣ ♦ ♥ ♠
	parseSuit := func(icon string) int {
		s := icon[len(icon)-3:]
		switch s {
		case "♣":
			return 1
		case "♦":
			return 2
		case "♥":
			return 3
		case "♠":
			return 4
		}
		return 0
	}

	// Build and shuffle deck of Card ICONS
	deck := NewDeck() // returns []string
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	// Assign one cut card per player
	cuts := make([]cut, len(players))
	for i, p := range players {
		icon := deck[i]
		cuts[i] = cut{
			player: p,
			icon:   icon,
			rank:   parseRank(icon),
			suit:   parseSuit(icon),
		}
	}

	// Sort by rank, suit
	sort.Slice(cuts, func(i, j int) bool {
		if cuts[i].rank != cuts[j].rank {
			if lowest {
				return cuts[i].rank < cuts[j].rank
			}
			return cuts[i].rank > cuts[j].rank
		}
		if lowest {
			return cuts[i].suit < cuts[j].suit
		}
		return cuts[i].suit > cuts[j].suit
	})

	// Return ordered players
	ordered := make([]string, len(players))
	for i := range cuts {
		ordered[i] = cuts[i].player
	}

	return ordered
}
