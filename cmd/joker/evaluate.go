package main

import (
	"cards/common"
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

func EvaluateNextCard(
	dealerPoints, playerPoints int, // current points in hand
	dealerRank, playerRank common.Rank, // Rank of next card
	wasDealerDouble, wasPlayerDouble bool) ( // if previous was King
	dealerResult, playerResult int, dealerDouble, playerDouble bool) {

	dealerResult = dealerPoints
	playerResult = playerPoints
	dealer := 0
	player := 0

	switch dealerRank {
	case common.Seven:
		dealer = 4
		playerResult -= 3
	case common.Ace:
		dealer = 1
	case common.Jack:
		playerResult -= 3
	case common.Queen:
		dealer = 12
	case common.King:
		dealerDouble = true
	case common.Joker:
		dealer = dealerPoints
	default:
		dealer = int(dealerRank)
	}
	if wasDealerDouble {
		dealerResult += dealer * 2
	} else {
		dealerResult += dealer
	}

	switch playerRank {
	case common.Seven:
		player = 4
		dealerResult -= 3
	case common.Ace:
		player = 1
	case common.Jack:
		dealerResult -= 3
	case common.Queen:
		player = 12
	case common.King:
		playerDouble = true
	case common.Joker:
		player = playerPoints
	default:
		player = int(playerRank)
	}
	if wasPlayerDouble {
		playerResult += player * 2
	} else {
		playerResult += player
	}
	//s := fmt.Sprintf("eval: dealer=%d,player=%d, result: dealer=%d,player=%d",
	//	dealer, player, dealerResult, playerResult)
	//fmt.Println(s)

	return
}
