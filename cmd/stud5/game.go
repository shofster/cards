package main

import (
	"cards/common"
	"fmt"
	"strconv"
	"time"

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
	dealerHand     [5]cardCell
	playerHand     [5]cardCell
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
	if street < 5 {
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
	setStreetData(0)
	setPayoutData("Dealing .........")
	g.NextStreet(false)
}

func (g *Game) NextStreet(bet bool) {

	//fmt.Printf("NextStreet %d: bet is %d, wallet is %d\n",
	//	g.street, g.pot, intGet(wallet))

	switch g.street {

	//
	// street 0 → 1 : new deal, ante, deal first 2 cards
	//
	case 0:
		g.shuffle = g.deck
		common.Shuffle(game.shuffle)

		addToWallet(-g.anteAmount)
		setPotData(g.anteAmount)

		// show all 5 reverse images
		for h := 0; h < 5; h++ {
			setCardFace(true, h, g.dealerHand[h], reverse)
			setCardFace(false, h, g.playerHand[h], reverse)
		}

		// set street 1 condition
		// dealer -- 2 cards, hole card down
		g.dealNextCard(&g.dealerHand[0])
		setCardFace(true, 0, g.dealerHand[0], reverse) // hidden
		g.dealNextCard(&g.dealerHand[1])
		setCardFace(true, 1, g.dealerHand[1], reverse)

		// player -- 2 cards, hole card down
		g.dealNextCard(&g.playerHand[0])
		setCardFace(false, 0, g.playerHand[0], reverse) // hidden
		g.dealNextCard(&g.playerHand[1])
		setCardFace(false, 1, g.playerHand[1], reverse)

		//time.Sleep(300 * time.Millisecond)

		// show cards for street 1 - dealer hole card hidden
		setCardFace(true, 1, g.dealerHand[1], nil)
		setCardFace(false, 0, g.playerHand[0], nil)
		setCardFace(false, 1, g.playerHand[1], nil)

		setStreetData(1)
		_ = handState.Set("Dealing..........")

	//
	// streets 1, 2, 3 → next street  (Optional bet)
	//
	case 1, 2, 3:
		// If betting allowed on this street AND player pressed Bet
		if bet && g.allowBetStreet[g.street] {
			addToWallet(-g.betAmount)
			setPotData(g.pot + g.betAmount)
			//	fmt.Printf("street %d: bet is %d, wallet is %d\n",
			//		g.street, g.pot, intGet(wallet))
			//} else {
			//	fmt.Printf("street %d: NO BET\n", g.street)
		}

		// deal next card index = street
		setStreetData(g.street + 1)
		idx := g.street
		g.dealNextCard(&g.playerHand[idx])
		g.dealNextCard(&g.dealerHand[idx])
		setCardFace(true, idx, g.dealerHand[idx], nil)
		setCardFace(false, idx, g.playerHand[idx], nil)

	case 4: // last card bet
		if bet {
			addToWallet(-g.betAmount)
			setPotData(g.pot + g.betAmount)
			//	fmt.Printf("street %d: bet is %d, wallet is %d\n",
			//		g.street, g.pot, intGet(wallet))
			//} else {
			//	fmt.Printf("street %d: NO BET\n", g.street)
		}
		setStreetData(g.street + 1)
		setCardFace(true, 0, g.dealerHand[0], nil)
		time.Sleep(time.Millisecond * 500)
		g.showdown()

	}

	g.updateBetButtons()
}

// pass or optionally bet for next street.
func (g *Game) updateBetButtons() {
	switch g.street {
	case 0:
		allowUpdateBets(false)
		dealButton.Enable()
		nextButton.Disable()
		betButton.Disable()
	case 1, 2, 3:
		dealButton.Disable()
		nextButton.Enable()
		if game.allowBetStreet[g.street] {
			betButton.Enable()
		} else {
			betButton.Disable()
		}
	case 4:
		allowUpdateBets(true)
		nextButton.Enable()
		betButton.Enable()
	case 5:
		allowUpdateBets(false)
		dealButton.Enable()
		nextButton.Disable()
		betButton.Disable()
	default:
		_ = handState.Set(fmt.Sprintf("Street %d .........", g.street))
		dealButton.Disable()
		if g.street >= 1 && g.street <= 3 && g.allowBetStreet[g.street] {
			betButton.Enable()
			nextButton.Enable()
		} else {
			betButton.Disable()
			nextButton.Enable()
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

	playerHand := make([]common.Card, 5)
	for i := 0; i < 5; i++ {
		playerHand[i] = g.playerHand[i].card
	}
	dealerHand := make([]common.Card, 5)
	for i := 0; i < 5; i++ {
		dealerHand[i] = g.dealerHand[i].card
	}
	playerRank := EvaluateFinal(playerHand)
	dealerRank := EvaluateFinal(dealerHand)

	// Dealer must have Ace-King to qualify
	dq, _ := qualify.Get()
	if dq && !DealerQualifies(dealerRank) {
		// win original ante × 2
		addToWallet(g.anteAmount * 2)
		// return bets
		addToWallet(g.pot - g.anteAmount)
		_ = handState.Set("Dealer unqualified, Ante wins, other bets returned")
		return
	}

	result := CompareFinal(playerRank, dealerRank)

	switch {
	case result > 0:

		// Player wins full pot × 2

		n := g.pot * 2
		_ = handState.Set(fmt.Sprintf("Player wins %d", n))
		addToWallet(n)
	case result < 0:
		// Dealer wins
		// Player already paid the pot; nothing returned
		_ = handState.Set(fmt.Sprintf("Dealer wins %d", g.pot))
	default:
		// Push: return pot to wallet
		_ = handState.Set(fmt.Sprintf("Push: %d returned to wallet", g.pot))
		addToWallet(g.pot)
	}

}
