package common

import (
	"math/rand/v2"
	"time"
)

/*

  File:    deck.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: manage a deck of Cards, shuffle, and deal.
*/

func NewDeck() []string {
	deck := make([]string, 0, 52)
	for _, suit := range []Suit{Clubs, Diamonds,
		Hearts, Spades} {
		for _, rank := range []Rank{Two, Three,
			Four, Five, Six, Seven,
			Eight, Nine, Ten, Jack,
			Queen, King, Ace} {
			code := RankCode(rank) + SuitCode(suit)
			deck = append(deck, code)
		}
	}
	return deck
}

func AddJokerToDeck(deck []string) []string {
	joker := make([]string, len(deck))
	copy(joker, deck)
	code := RankCode(Joker)
	joker = append(joker, code)
	return joker
}

func Shuffle(deck []string) {
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

// Deal nCards to nHands, return remaining deck.
func Deal(deck []string, nHands, nCards int) ([]Cards, []string) {
	hands := make([]Cards, nHands)
	for i := range hands {
		hands[i] = make([]Card, 0)
	}
	idx := 0
	for c := 0; c < nCards; c++ {
		for h := 0; h < nHands; h++ {
			code := deck[idx]
			hands[h] = append(hands[h], CardMap()[code])
			idx++
			if idx >= len(deck) {
				return hands, nil
			}
		}
	}
	return hands, deck[idx:]

}
