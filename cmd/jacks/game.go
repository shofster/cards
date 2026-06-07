package main

import (
	"cards/common"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    game.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/

type cardCell struct {
	image  *fyne.Container
	action *widget.Check
	card   common.Card
}

var cardMap = map[string]common.Card{}
var cardCells = make([]cardCell, 5)
var dealButton *widget.Button
var drawButton *widget.Button

type Game struct {
	Deck       []string // 52 card codes (Icon)
	Shuffle    []string // remaining shuffled cards
	betAmount  *widget.Select
	payoutType *widget.Select
}

var game Game

func initGame() {
	cardMap = common.CardMap()   // icon vs Card
	game.Deck = common.NewDeck() // list of 52 Card icons
}

func startGame() {

	for h := 0; h < 5; h++ { // h is scoped to loop - since 1.22
		cardCells[h].action.SetText("keep")
		cardCells[h].action.SetChecked(true)
	}
}

// fill Shuffle from Deck
// set up 5 cards in Hand
// put remaining back to Shuffle.
func deal() {

	drawButton.Enable()
	game.betAmount.Disable()
	game.payoutType.Disable()

	game.Shuffle = game.Deck
	common.Shuffle(game.Shuffle)
	ix := 0
	hand := make([]common.Card, 5)
	for h := 0; h < 5; h++ {
		hand[h] = cardMap[game.Shuffle[ix+h]]
		cardCells[h].card = hand[h]
		cardCells[h].card.Selected = true
		setCardFace(cardCells[h].image, cardCells[h].card.Icon, true)
		cardCells[h].action.OnChanged = func(keep bool) {
			setCardFace(cardCells[h].image, cardCells[h].card.Icon, keep)
			cardCells[h].card.Selected = keep
			if !keep {
				cardCells[h].action.SetText(" " + cardCells[h].card.Icon)
			} else {
				cardCells[h].action.SetText("keep")
			}
		}
		cardCells[h].action.SetChecked(true)
	}

	// eliminate the first 5 cards from being dealt
	game.Shuffle = game.Shuffle[5:]
	// deduct bet amount from wallet
	wager, _ := bet.Get()
	betAmount, _ := strconv.Atoi(wager)
	addBetAmountToWallet(-betAmount)

	// first pair > 10
	switch DealAnnouncement(hand) {
	case common.Jack:
		_ = handResult.Set("Pair of Jacks")
	case common.Queen:
		_ = handResult.Set("Pair of Queens")
	case common.King:
		_ = handResult.Set("Pair of Kings")
	case common.Ace:
		_ = handResult.Set("Pair of Aces")
	default:
		_ = handResult.Set("NO Pair of Jacks or Better")
	}

}

func draw() {

	drawButton.Disable()
	game.betAmount.Enable()
	game.payoutType.Enable()
	// printHand("beore draw")

	// draw replacement cards
	hand := make([]common.Card, 5)
	ix := 0
	for h := 0; h < 5; h++ {
		cardCells[h].action.SetText("held")
		hand[h] = cardCells[h].card
		if !cardCells[h].card.Selected {
			hand[h] = cardMap[game.Shuffle[ix]]
			cardCells[h].card = hand[h]
			cardCells[h].card.Icon = hand[h].Icon
			setCardFace(cardCells[h].image, cardCells[h].card.Icon, true)
			cardCells[h].action.SetText(cardCells[h].card.Icon)
			cardCells[h].action.SetChecked(false)
			ix++
		}
	}

	// printHand("draw")

	result := EvaluateFinal(hand)
	table, _ := payout.Get()
	wager, _ := bet.Get()
	betAmount, _ := strconv.Atoi(wager)
	won := Payout(table, result, betAmount)
	addBetAmountToWallet(won)
	s := fmt.Sprintf("%s - Payout is %d", result, won)
	_ = handResult.Set(s)
}
