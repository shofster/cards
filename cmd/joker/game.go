package main

import (
	"fmt"
	"games/common"
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
var playerPointsData = binding.NewString()
var dealerPointsData = binding.NewString()

var betButton *widget.Button
var nextButton *widget.Button
var dealButton *widget.Button

type Game struct {
	deck           []string // 52 card codes (Icon)
	shuffle        []string // remaining shuffled cards
	street         int      // deal state
	dealerPoints   int
	playerPoints   int
	dealerDouble   bool
	playerDouble   bool
	dealerHand     [5]cardCell
	playerHand     [5]cardCell
	allowBetStreet [5]bool // optional for streets 1,2,3,4
	anteAmount     int     // from settings
	betAmount      int     // from settings
	pot            int
}

func initGame() {
	game = Game{}
	cardMap = common.CardMap()                   // icon vs Card
	game.deck = common.NewDeck()                 // 52 cards
	game.deck = common.AddJokerToDeck(game.deck) // plus 1 joker
	setStreetData(0)
	setPointsData(true, 0)
	setPointsData(false, 0)
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
		case 4:
			game.allowBetStreet[s] = boolGet(street4bet)
		default:
			game.allowBetStreet[s] = false
		}
	}
	// the common card back
	reverse = canvas.NewImageFromResource(cardMap["reverse"].Resource)
}

func setStreetData(street int) {
	game.street = street
	if street < 5 {
		_ = streetData.Set("  " + strconv.Itoa(game.street) + "  ")
	} else {
		_ = streetData.Set("  SHOWDOWN !  ")
	}
}

func setPointsData(dealer bool, points int) {
	if dealer {
		_ = dealerPointsData.Set(strconv.Itoa(points))
	} else {
		_ = playerPointsData.Set(strconv.Itoa(points))
	}
}

func (g *Game) setPotData(pot int) {
	g.pot = pot
	_ = potData.Set(strconv.Itoa(pot))
}

func (g *Game) deal() {
	setStreetData(0)
	_ = handState.Set("Dealing .........")
	g.NextStreet(false)
}

func (g *Game) NextStreet(bet bool) {

	//fmt.Printf("NextStreet %d: bet is %d\n", g.street, g.pot)

	switch g.street {

	//
	// street 0 → 1 : new deal, ante, deal first 2 cards
	//
	case 0:
		g.shuffle = g.deck
		common.Shuffle(game.shuffle)

		addToWallet(-g.anteAmount)
		g.setPotData(g.anteAmount)
		g.dealerPoints = 0
		g.playerPoints = 0
		g.dealerDouble = false
		g.playerDouble = false

		// show all 5 reverse images
		for h := 0; h < 5; h++ {
			setCardFace(true, h, g.dealerHand[h], reverse)
			setCardFace(false, h, g.playerHand[h], reverse)
		}

		// set street 1 condition
		// dealer & player -- 1 card obverse up
		idx := g.street
		g.dealNextCard(&g.dealerHand[idx])
		g.dealNextCard(&g.playerHand[idx])
		g.dealerPoints, g.playerPoints, g.dealerDouble, g.playerDouble =
			EvaluateNextCard(g.dealerPoints, g.playerPoints,
				g.dealerHand[idx].card.Rank, g.playerHand[idx].card.Rank,
				g.dealerDouble, g.playerDouble)

		setPointsData(true, g.dealerPoints)
		setPointsData(false, g.playerPoints)
		// show cards for street 1
		setCardFace(true, idx, g.dealerHand[idx], nil)
		setCardFace(false, idx, g.playerHand[idx], nil)

		_ = handState.Set("Dealing..........")
		g.street = 1
		setStreetData(1)

	//
	// streets 1, 2, 3, 4 → next street  (Optional bet)
	//
	case 1, 2, 3, 4:
		// If betting allowed on this street AND player pressed Bet
		if bet && g.allowBetStreet[g.street] {
			addToWallet(-g.betAmount)
			g.setPotData(g.pot + g.betAmount)
		}

		// deal next card index = street
		idx := g.street
		g.street++
		setStreetData(g.street)
		g.dealNextCard(&g.playerHand[idx])
		g.dealNextCard(&g.dealerHand[idx])
		setCardFace(true, idx, g.dealerHand[idx], nil)
		setCardFace(false, idx, g.playerHand[idx], nil)
		g.dealerPoints, g.playerPoints, g.dealerDouble, g.playerDouble =
			EvaluateNextCard(g.dealerPoints, g.playerPoints,
				g.dealerHand[idx].card.Rank, g.playerHand[idx].card.Rank,
				g.dealerDouble, g.playerDouble)

		setPointsData(true, g.dealerPoints)
		setPointsData(false, g.playerPoints)
		if g.street == 5 {
			setStreetData(g.street)
			g.showdown()
		}

	}

	g.updateBetButtons()
}

// pass or optionally bet for next street.
func (g *Game) updateBetButtons() {
	switch g.street {
	case 0, 5: // Deal, showdown
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

	result := g.playerPoints - g.dealerPoints
	addPoints(playerPoints, g.playerPoints)
	addPoints(dealerPoints, g.dealerPoints)

	switch {
	case result > 0:
		_ = handState.Set(fmt.Sprintf("Player wins %d", g.pot*2))
		addToWallet(g.pot * 2)
	case result < 0:
		// Dealer wins
		// Player already paid the pot; nothing returned
		_ = handState.Set(fmt.Sprintf("Player loses %d", g.pot))
	default:
		// Push: return pot to wallet
		_ = handState.Set(fmt.Sprintf("Push: %d returned to wallet", g.pot))
		addToWallet(g.pot)
	}
}
