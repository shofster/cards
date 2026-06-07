package common

import "cards/assets"

/*

  File:    CardMap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

func CardMap() map[string]Card {

	cards := make(map[string]Card)

	cards["reverse"] = Card{
		Rank:     Reverse,
		Icon:     "reverse",
		Resource: assets.ResourceByCode["reverse"],
	}
	cards["joker"] = Card{
		Rank:     Joker,
		Icon:     "joker",
		Resource: assets.ResourceByCode["joker"],
	}
	for s := Suit(Clubs); s <= Spades; s = s + SuitInc {
		for r := Rank(Two); r <= Ace; r = r + RankInc {
			icon := Icon(r, s)
			resource, ok := assets.ResourceByCode[icon]
			if !ok {
				panic("CardMap is missing " + icon)
			}
			c := Card{
				Rank:     r,
				Suit:     s,
				Icon:     icon,
				Resource: resource,
			}
			cards[c.Icon] = c
		}
	}

	return cards
}
