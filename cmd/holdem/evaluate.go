package main

import (
	"cards/common"
	"fmt"
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
A random 5‑card hand is overwhelmingly likely to be:
- High card — about 50% of all deals
- One pair — about 42%
- Two pair — about 4.75%
- Trips — about 2.1%
- Straight — about 0.39%
- Flush — about 0.20%
- Full house — about 0.14%
- Quads — about 0.024%
- Straight flush — about 0.0015%
So more than 90% of all hands are either high card or a single low pair.

*/

func EvaluateDeal(cards []common.Card) (result int, hands string) {

	dealerHand := make([]common.Card, 0)
	playerHand := make([]common.Card, 0)

	for idx := 0; idx < 9; idx++ {
		switch idx {
		case 0, 1:
			dealerHand = append(dealerHand, cards[idx])
		case 2, 3:
			playerHand = append(playerHand, cards[idx])
		default:
			dealerHand = append(dealerHand, cards[idx])
			playerHand = append(playerHand, cards[idx])
		}
	}

	dealerRank := common.Best5of7(dealerHand)
	playerRank := common.Best5of7(playerHand)

	result = common.CompareFinal(playerRank, dealerRank)
	hands = fmt.Sprintf("Result: player %s, dealer %s",
		playerRank.Rank.String(), dealerRank.Rank.String())

	//fmt.Printf("Results %s: player %s, dealer %s\n",
	//	winner(result), playerRank.Rank.String(), dealerRank.Rank.String())
	return
}

func winner(comp int) string {
	if comp > 0 {
		return "Player Wins - "
	} else if comp < 0 {
		return "Player Loses - "
	}
	return "PUSH"
}
