package common

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

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

func (c *Card) String() string {
	return c.Icon
}

// parse Rank and Suit from the Card's Icon
func (c *Card) RankSuit() (rank Rank, suit Suit) {
	r := []rune(c.Icon)
	rankStr := string(r[:len(r)-1])
	suitRune := r[len(r)-1]
	// parse rank
	switch rankStr {
	case "A":
		rank = Ace
	case "K":
		rank = King
	case "Q":
		rank = Queen
	case "J":
		rank = Jack
	default: //
		n, _ := strconv.Atoi(rankStr)
		rank = Rank(n)
	}
	// parse suit
	switch suitRune {
	case '♣':
		suit = Clubs
	case '♦':
		suit = Diamonds
	case '♥':
		suit = Hearts
	case '♠':
		suit = Spades
	}
	return
}

type Cards []Card

func (cs Cards) String() string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Icon
	}
	return strings.Join(out, " ")
}

func Contains(cards Cards, card Card) bool {
	for _, x := range cards {
		if x.Icon == card.Icon {
			return true
		}
	}
	return false
}

func (cards Cards) Remove(card Card) Cards {
	out := make(Cards, 0, len(cards))
	for _, x := range cards {
		if x.Icon != card.Icon {
			out = append(out, x)
		}
	}
	return out
}

// IsRun checks whether the given cards form a run of 3 (order doesn't matter).
func IsRun(cards Cards) bool {
	n := len(cards)
	if n < 3 {
		return false
	}

	ranks := make([]int, n)
	for i, c := range cards {
		ranks[i] = int(c.Rank)
	}

	sort.Ints(ranks)

	for i := 1; i < n; i++ {
		if ranks[i] != ranks[i-1]+1 {
			return false
		}
	}
	return true
}

// SubsetsOf returns all subsets of 'cards' of length 'k'.
// Order inside each subset is preserved.
func SubsetsOf(cards Cards, k int) []Cards {
	n := len(cards)
	if k > n || k <= 0 {
		return nil
	}

	// initial combination indices: [0,1,2,...,k-1]
	idx := make([]int, k)
	for i := 0; i < k; i++ {
		idx[i] = i
	}

	out := make([]Cards, 0, 32)

	for {
		// build subset
		subset := make(Cards, k)
		for i := 0; i < k; i++ {
			subset[i] = cards[idx[i]]
		}
		out = append(out, subset)

		// generate next combination
		i := k - 1
		for i >= 0 && idx[i] == i+n-k {
			i--
		}
		if i < 0 {
			break
		}

		idx[i]++
		for j := i + 1; j < k; j++ {
			idx[j] = idx[j-1] + 1
		}
	}

	return out
}
