package main

import (
	"cards/common"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
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
	card  common.Card
	image *canvas.Image
}

var cardMap = map[string]common.Card{}
var game Game
var reverse *canvas.Image

var potData = binding.NewString()
var streetData = binding.NewString()
var payoutData = binding.NewString()

var betButton *widget.Button
var nextButton *widget.Button
var dealButton *widget.Button

type Game struct {
	deck           []string // 52 card codes (Icon)
	shuffle        []string // remaining shuffled cards
	street         int      // deal state
	pot            int
	cards          [9]cardCell
	allowBetStreet [5]bool // optional for streets 1,2,3
	anteAmount     int     // from settings
	betAmount      int     // from settings
}

func initGame() {
	game = Game{}
	cardMap = common.CardMap()   // icon vs Card
	game.deck = common.NewDeck() // 52 cards
	game.shuffle = make([]string, len(game.deck))
	setPotData(0)
	setStreetData(0)
	setPayoutData("")
	dealButton.OnTapped = func() { game.deal() }
	betButton.OnTapped = func() { game.NextStreet(true) }
	nextButton.OnTapped = func() { game.NextStreet(false) }
	for s := 0; s < 5; s++ {
		switch s {
		case 1:
			game.allowBetStreet[s] = boolGet(street1bet)
		case 2:
			game.allowBetStreet[s] = boolGet(street2bet)
		case 3:
			game.allowBetStreet[s] = boolGet(street3bet)
		default:
			game.allowBetStreet[s] = false
		}
	}
	// the common card back
	reverse = canvas.NewImageFromResource(cardMap["reverse"].Resource)
}

func setStreetData(street int) {
	game.street = street
	if street < 4 {
		_ = streetData.Set("  " + strconv.Itoa(game.street) + "  ")
	} else {
		_ = streetData.Set("  SHOWDOWN !  ")
	}

}

func setPotData(pot int) {
	game.pot = pot
	_ = potData.Set(strconv.Itoa(game.pot))
}

func setPayoutData(po string) {
	_ = payoutData.Set(po)
}

func (g *Game) deal() {
_:
	setStreetData(0)
	//  reverse all 9 images
	resetCardFaces()
	_ = handState.Set("Dealing .........")
	g.NextStreet(false)
}

func (g *Game) NextStreet(bet bool) {

	//fmt.Printf("NextStreet %d: bet is %d\n", g.street, g.pot)

	switch g.street {

	//
	// street 0 → 1 : new deal, ante, deal first 2 cards, flop cards
	//
	case 0:
		g.shuffle = g.deck
		common.Shuffle(game.shuffle)

		addToWallet(-g.anteAmount)
		setPotData(g.anteAmount)

		for i := 0; i < 9; i++ {
			g.dealNextCard(&g.cards[i])
		}
		for i := 2; i <= 6; i++ {
			setCardFace(i, g.cards[i].image)
		}
		_ = handState.Set("Dealing..........")
		g.street = 1
		setStreetData(1)

	//
	// streets 1, 2, 3 → next street  (Optional bet)
	//
	case 1, 2:
		// If betting allowed on this street AND player pressed Bet
		if bet && g.allowBetStreet[g.street] {
			addToWallet(-g.betAmount)
			setPotData(g.pot + g.betAmount)
		}

		g.street++
		setStreetData(g.street)

		setCardFace(g.street+5, g.cards[g.street+5].image)

	case 3:
		if bet && g.allowBetStreet[g.street] {
			addToWallet(-g.betAmount)
			setPotData(g.pot + g.betAmount)
		}
		g.street++
		setStreetData(g.street)

		if g.street == 4 {
			setCardFace(0, g.cards[0].image)
			setCardFace(1, g.cards[1].image)
			setStreetData(g.street)
			g.showdown()
		}

	}

	g.updateBetButtons()
}

// pass or optionally bet for next street.
func (g *Game) updateBetButtons() {
	switch g.street {
	case 0, 4: // Deal, showdown
		allowUpdateBets(true)
		dealButton.Enable()
		nextButton.Disable()
		betButton.Disable()
	default:
		allowUpdateBets(false)
		dealButton.Disable()
		nextButton.Enable()
		if game.allowBetStreet[g.street] {
			betButton.Enable()
		} else {
			betButton.Disable()
		}
	}
}

// set the next cardCell properties
func (g *Game) dealNextCard(cell *cardCell) {
	icon := g.shuffle[0]
	cell.card = cardMap[icon]
	cell.image = canvas.NewImageFromResource(cell.card.Resource)
	g.shuffle = g.shuffle[1:]
}

func (g *Game) showdown() {

	cards := make([]common.Card, 9)
	for i := 0; i < 9; i++ {
		cards[i] = g.cards[i].card
	}
	winner, result := EvaluateDeal(cards)

	switch {
	case winner > 0:
		_ = handState.Set(fmt.Sprintf("%s - Player wins %d ", result, g.pot*2))
		addToWallet(g.pot * 2)
	case winner < 0:
		// Dealer wins
		// Player already paid the pot; nothing returned
		_ = handState.Set(fmt.Sprintf("%s - Player loses %d", result, g.pot))
	default:
		// Push: return pot to wallet
		_ = handState.Set(fmt.Sprintf("%s - Push: %d returned to wallet", result, g.pot))
		addToWallet(g.pot)
	}

}
