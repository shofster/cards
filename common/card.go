package common

import (
	"fmt"

	"fyne.io/fyne/v2"
)

/*

  File:    Card.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

type Suit int

const (
	Clubs = iota
	Diamonds
	Hearts
	Spades
	NoTrump
	NoSuit = -1
)
const SuitInc Suit = 1

type Rank int

const (
	Reverse Rank = iota
	Joker
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)
const RankInc Rank = 1

type Card struct {
	Rank     Rank
	Suit     Suit
	Icon     string
	Resource fyne.Resource
	Selected bool
	Index    int
}

func RankCode(r Rank) string {
	switch r {
	case Joker:
		return "joker"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return fmt.Sprintf("%d", r)
	}
}

func SuitCode(s Suit) string {
	switch s {
	case Clubs:
		return "♣"
	case Diamonds:
		return "♦"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	}
	return "?"
}

func Icon(r Rank, s Suit) string {

	return RankCode(r) + SuitCode(s)
}

// AceHighOrdinal maps Rank (Ace=14) into Ace-high (default) ordering.
func AceHighOrdinal(r Rank) int {
	return int(r)
}

// AceLowOrdinal maps Rank (Ace=14) into Ace-low ordering.
// i.e.-- Ace → 1, Two → 2, ..., King → 13.
func AceLowOrdinal(r Rank) int {
	switch r {
	case Ace:
		return 1
	case Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten,
		Jack, Queen, King:
		return int(r) // 2 thru 13
	case Reverse:
		return -1
	}
	return 0
}
