package main

import (
	"games/common"
	"log"

	"fyne.io/fyne/v2/canvas"
)

/*

  File:    game.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var cardMap = map[string]common.Card{}
var cardDeck []string

type Game struct {
	Stock           []string
	Discard         string
	Hands           [][]common.Card // [0]=player, [1]=dealer
	discardSelected int
	discard         common.Card
	started         bool
}

var game Game

func initGame() {
	cardMap = common.CardMap()
	cardDeck = common.NewDeck()
	game = Game{}
}

func deal() {

	game.started = false

	// shuffle the full deck
	game.Stock = cardDeck
	common.Shuffle(game.Stock)

	// deal 2 hands of 10 cards each
	var rest []string
	game.Hands, rest = common.Deal(game.Stock, 2, 10)
	common.SortHand(game.Hands[0])
	common.SortHand(game.Hands[1])

	// first discard is rest[0]
	game.Discard = rest[0]

	// stock is the remainder
	game.Stock = rest[1:]

	// show Melds - dealer in dev only
	_ = playerMeld.Set(computeMelds(game.Hands[0]).toString())
	if devMode {
		_ = dealerMeld.Set(computeMelds(game.Hands[1]).toString())
	}

	// set up the stock & discard piles
	stock.Img = canvas.NewImageFromResource(cardMap[game.Stock[0]].Resource)
	//stock.Img = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	discard.Img = canvas.NewImageFromResource(cardMap[game.Discard].Resource)
	stock.Refresh()
	discard.Refresh()
	stock.OnTap = replaceCard
	discard.OnTap = replaceCard

	resetHands() // display reverse images
	for i := 0; i < 10; i++ {
		// development -- show new hand
		if devMode {
			setDealerCard(i, game.Hands[1][i])
		}

		setPlayerCard(i, game.Hands[0][i])
		ci := playerHand.Objects[i].(*common.CardImage)
		ci.OnTap = onTap
	}

	// player begins by selecting a discard
	enablePlayerHand(true)
	playerKnockButton.Enable()
	dealerPlayButton.Disable()
	enableCenter(false)

}

/*

KNOCK WIN:      winner gets (deadwoodDiff)
UNDERCUT:       winner gets (abs(deadwoodDiff) + 25)
GIN:            winner gets (opponentDeadwood + 25)
LOSER:          always 0

RESULT:   tells you direction (win/loss)
SCORE:    tells winner points earned

*/

func processKnock(knocker int) {
	knockee := 1 - knocker

	var result Result
	var rawDiff int    // true scoring diff (knocker perspective)
	var resultDiff int // player-perspective diff (UI only)
	var layoffDelta int

	if !game.started { // first turn - draw
		playerResults := computeMelds(game.Hands[0])
		dealerResults := computeMelds(game.Hands[1])

		result.KnockerPoints = playerResults.DeadwoodPoints
		result.KnockeePoints = dealerResults.DeadwoodPoints
		result.Gin = false
		result.KnockedOut = false

		rawDiff = 0
		resultDiff = 0
	} else {
		// original deadwood for knockee
		preKnockee := computeMelds(game.Hands[knockee]).DeadwoodPoints
		// resolution
		result = resolveKnock(game.Hands[knocker], game.Hands[knockee])
		// scoring diff (knocker perspective)
		rawDiff = result.Diff
		// scoring diff (player perspective)
		resultDiff = rawDiff * (1 - 2*knocker)
		// Layoff delta
		postKnockee := result.KnockeePoints
		layoffDelta = preKnockee - postKnockee
	}

	// display winner position
	switch {
	case resultDiff == 0:
		_ = winnerData.Set("NOONE - IT's a DRAW")
	case resultDiff > 0:
		_ = winnerData.Set("PLAYER")
	case resultDiff < 0:
		_ = winnerData.Set("DEALER")
	}

	// BASE SCORING (diff-based, not deadwood-based)
	knockerPoints := 0
	knockeePoints := 0

	if rawDiff > 0 { // knocker wins
		knockerPoints = rawDiff
	} else if rawDiff < 0 { // knockee wins (undercut)
		knockeePoints = -rawDiff
	}

	// GIN BONUS
	if result.Gin {
		if rawDiff > 0 { // knocker gins
			knockerPoints += result.KnockeePoints + 25
		} else { // knockee gins
			knockeePoints += result.KnockerPoints + 25
		}
	}

	// UNDERCUT BONUS
	if result.KnockedOut && rawDiff < 0 {
		knockeePoints += 25
	}
	if result.KnockedOut && rawDiff > 0 {
		knockerPoints += 25
	}

	// PLAYER/DEALER points
	var playerPoints, dealerPoints int
	if knocker == 0 {
		playerPoints = knockerPoints
		dealerPoints = knockeePoints
	} else {
		playerPoints = knockeePoints
		dealerPoints = knockerPoints
	}

	// display the outcome
	var outcome string
	switch {
	case result.Gin:
		outcome = "GIN (+25)"
	case result.KnockedOut && rawDiff != 0:
		outcome = "UNDERCUT (+25)"
	case layoffDelta > 0:
		outcome = "KNOCK (with layoffs)"
	default:
		outcome = "KNOCK"
	}
	_ = outcomeData.Set(outcome)

	_ = playerData.Set(playerPoints)
	_ = dealerData.Set(dealerPoints)
	_ = handResult.Set(resultDiff)

	showDealerHand()
	if !devMode {
		_ = dealerMeld.Set(computeMelds(game.Hands[1]).toString())
	}

	resetPlayerHand()
	playerKnockButton.Disable()
	showScoreCenter()
}

// logic for the dealer to replace 1 card and check for KNOCK / GIN
func dealerPlay() {

	enablePlayerHand(false)
	playerKnockButton.Disable()
	defer dealerPlayButton.Disable()
	defer playerKnockButton.Enable()

	hand := game.Hands[1]

	// draw source ?
	var newCard common.Card
	if shouldTakeDiscard(cardMap[game.Discard], hand) {
		// take discard card
		newCard = cardMap[game.Discard]
		discard.SetDimmed(true)
	} else {
		if len(game.Stock) <= 1 {
			// Hand ends in a draw
			processKnock(-1)
			return
		}
		newCard = cardMap[game.Stock[0]]
		game.Stock = game.Stock[1:]
		// update stock UI
		stock.Replace(cardMap[game.Stock[0]].Resource)
	}

	// add to hand -- now have 11
	hand = append(hand, newCard)

	// choose one and remove discard
	discardCard := chooseDiscard(hand)
	log.Printf("Dealer discards %s\n", discardCard.Icon)
	newHand := make([]common.Card, 0, len(hand)-1)
	for _, c := range hand {
		if !sameCard(c, discardCard) {
			newHand = append(newHand, c)
		}
	}

	game.Hands[1] = newHand
	if devMode {
		showDealerHand()
	}

	melds := computeMelds(game.Hands[1])
	if devMode {
		_ = dealerMeld.Set(melds.toString())
	}

	// replace discard and image
	game.Discard = discardCard.Icon
	discard.Replace(discardCard.Resource)

	if melds.DeadwoodPoints == 0 {
		processKnock(1)
	} else if melds.DeadwoodPoints <= dealerKnock {
		processKnock(1)
	} else {
		enablePlayerHand(true)
	}
}

// called when player KNOCK Button pressed
func knock() {
	processKnock(0)
}

// stage a selected Card to replace
func onTap(ci *common.CardImage) {
	if ci.Highlighted {
		ci.SetHighlite(false)
		game.discardSelected = -1
		ci.SetDimmed(false)
		enableCenter(false)
		playerKnockButton.Enable()
	} else {
		if game.discardSelected != -1 {
			playerHand.Objects[game.discardSelected].(*common.CardImage).SetHighlite(false)
			playerHand.Objects[game.discardSelected].(*common.CardImage).SetDimmed(false)
		}
		game.discardSelected = ci.ID
		game.discard = game.Hands[0][ci.ID]
		ci.SetHighlite(true)
		ci.SetDimmed(true)
		enableCenter(true)
		playerKnockButton.Disable()
	}
}

// replace Player's staged card with either the Stock[0] or Discard
func replaceCard(ci *common.CardImage) {

	defer dealerPlayButton.Enable()
	defer func() {
		game.started = true
		enablePlayerHand(false)
	}()

	idx := game.discardSelected
	if idx == -1 {
		return
	}

	hi := playerHand.Objects[idx].(*common.CardImage)
	var newCard common.Card

	switch ci.ID {
	case 0: // Stock

		if len(game.Stock) <= 1 {
			// Hand ends in a draw
			processKnock(-1)
			return
		}
		// take top of stock
		newCard = cardMap[game.Stock[0]]
		// remove from stock
		game.Stock = game.Stock[1:]

		stock.Replace(cardMap[game.Stock[0]].Resource)

	case 1: // Discard
		// take top of discard
		newCard = cardMap[game.Discard]
	}

	// update player hand model
	game.Hands[0][idx] = newCard
	// update UI
	setPlayerCard(idx, newCard)

	game.discardSelected = -1
	game.Discard = game.discard.Icon
	discard.Replace(game.discard.Resource)

	_ = playerMeld.Set(computeMelds(game.Hands[0]).toString())

	// clean up the new CardImage
	hi.SetHighlite(false)
	hi.SetDimmed(false)
	enableCenter(false)

}

// Show the dealer's hand
var showDealerHand = func() {
	hand := game.Hands[1]
	for i := 0; i < 10; i++ {
		if i < len(hand) {
			setDealerCard(i, hand[i])
		} else {
			setDealerCard(i, cardMap["reverse"])
		}
	}
}

func enableCenter(enable bool) {
	stock.SetDimmed(!enable)
	discard.SetDimmed(!enable)
}

func enablePlayerHand(enable bool) {
	for i := 0; i < 10; i++ {
		ci := playerHand.Objects[i].(*common.CardImage)
		ci.OnTap = onTap
		ci.Enabled = enable
	}
}
