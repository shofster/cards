package common

import "sort"

/*

  File:    sorter.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: sort a "hand" of Cards by suit ♣♦♥♠ and descending rank.
    NoTrump and non 52 are last.
*/

func suitOrder(s Suit) int {
	switch s {
	case Clubs:
		return 0
	case Diamonds:
		return 1
	case Hearts:
		return 2
	case Spades:
		return 3
	case NoTrump:
		return 4
	default:
		return 99
	}
}

// Alternate order to make red,black,res,black
func altSuitOrder(s Suit) int {
	switch s {
	case Clubs:
		return 1
	case Diamonds:
		return 0
	case Hearts:
		return 2
	case Spades:
		return 3
	case NoTrump:
		return 4
	default:
		return 99
	}
}

func rankOrder(r Rank) int {
	switch r {
	case Ace:
		return 0
	case King:
		return 1
	case Queen:
		return 2
	case Jack:
		return 3
	case Ten:
		return 4
	case Nine:
		return 5
	case Eight:
		return 6
	case Seven:
		return 7
	case Six:
		return 8
	case Five:
		return 9
	case Four:
		return 10
	case Three:
		return 11
	case Two:
		return 12
	case Joker:
		return 98
	case Reverse:
		return 99
	default:
		return 100
	}
}

func SortHand(cards []Card) {

	sort.Slice(cards, func(i, j int) bool {
		ci, cj := cards[i], cards[j]
		// Suit
		si, sj := suitOrder(ci.Suit), suitOrder(cj.Suit)
		if si != sj {
			return si < sj
		}
		// Suit by Rank
		return rankOrder(ci.Rank) < rankOrder(cj.Rank)
	})
}

func AltSortHand(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		ci, cj := cards[i], cards[j]
		// Suit
		si, sj := altSuitOrder(ci.Suit), altSuitOrder(cj.Suit)
		if si != sj {
			return si < sj
		}
		// Suit by Rank
		return rankOrder(ci.Rank) < rankOrder(cj.Rank)
	})
}
