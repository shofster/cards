package main

import (
	"cards/common"
	"fmt"
	"log"
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

1. If I can score (15, 31, pair, run), do it.
2. Else play the lowest safe card.
3. If nothing is safe, dump the highest card.

*/

var (
	cardDeck []string       // current card Icons left in deck
	dealer   = -1           // pone(player) [0], dealer(ai)[1]
	hands    []common.Cards // pone(player) [0], dealer(ai)[1]
	crib     []common.Card
	lanes    []*common.CardStrip
	count    int
	starter  common.Card
	pegged   []common.Card
	keeps    [2][]common.Card
)

var (
	aiGO, playerGO bool
	bonusPoint     bool
	playerCard     *common.CardImage
	aiCard         *common.CardImage
)

/* -------------------------------------- */

func peggingAction(who int) {

	// Tap handler for selecting/unselecting pegging cards
	selectCard := func(cardImage *common.CardImage) {
		if playerCard != nil { // second tap - "un-choose"
			playerCard.SetDim(false) // Un-Dim, but leave Enabled
			playerCard = nil
			return
		}
		// Dim and choose this card
		playerCard = cardImage
		playerCard.SetDim(true)
	}
	// emulate the picking of a card
	selectAiCard := func() {
		aiCard = nil
		card, ok := choosePegCard(hands[AI], pegged, count)
		if !ok {
			panic("AI had no playable card but selectAiCard was called")
		}
		// Dim and choose this card
		aiCard = aiStack.Objects[card.Index].(*common.CardImage)
		aiCard.SetDim(true)
	}
	// see if a (PLAYER or AI) hand has any playable card
	canPlay := func(hand common.Cards) bool {
		for i := range hand {
			c := &hand[i]
			if !c.Selected && c.PegValue()+count <= 31 {
				return true
			}
		}
		return false
	}
	// eliminate cards that would exceed 31 from being picked
	limit := func(who int) {
		for i := range hands[PLAYER] {
			switch who {
			case PLAYER:
				c := &hands[PLAYER][i]
				if !c.Selected && c.PegValue()+count > 31 {
					ci := playerStack.Objects[c.Index].(*common.CardImage)
					ci.SetDimmed(true)
				}
			case AI:
				c := &hands[AI][i]
				if !c.Selected && c.PegValue()+count > 31 {
					ci := aiStack.Objects[c.Index].(*common.CardImage)
					ci.SetDimmed(true)
					hands[AI][ci.ID].Selected = true
				}
			}
		}
	}

	// play, peg, lane, and score the Card
	playCard := func(who int, card *common.Card) {
		card.Selected = true
		count += card.PegValue()
		pegged = append(pegged, *card)
		n := len(pegged) - 1
		switch who { // card is
		case PLAYER:
			lanes[PLAYER].SetCard(n, fmt.Sprintf("%3s", card.Icon))
			ci := playerStack.Objects[card.Index].(*common.CardImage)
			//ci.Replace(Ctx.cardMap["reverse"].Resource)
			ci.Enabled = false
		case AI:
			lanes[AI].SetCard(n, fmt.Sprintf("%3s", card.Icon))
			//ci := aiStack.Objects[card.Index].(*common.CardImage)
			//ci.Replace(Ctx.cardMap["reverse"].Resource)
		}
		common.IntSet(Ctx.handPoints, count)
	}

	//  Process Initialization, a Button press, or AI's first card  \\

	switch who {
	case -1: // prior to FIRST button press → set up pegging loop
		aiCard = nil
		playerCard = nil
		aiGO = false
		playerGO = false
		bonusPoint = false
		peggingButtonState(0)
		// enable tapping on all undiscarded cards
		for i := range playerStack.Objects {
			ci := playerStack.Objects[i].(*common.CardImage)
			if !hands[PLAYER][ci.ID].Selected {
				ci.Enabled = true
				ci.OnTap = selectCard
			}
		}
		return
	case PLAYER: // Button pressed
		if playerCard == nil { // but NO tapped card
			peggingButtonState(2) // re-assert "play a card"
			return
		}
		// play and score the card
		playCard(PLAYER, &hands[PLAYER][playerCard.ID])
		// played after AI was GO?
		if aiGO && !bonusPoint {
			scorePoints(PLAYER, 1, "bonus")
			bonusPoint = true
		}
		// award pegging points
		scorePegging(PLAYER, hands[PLAYER][playerCard.ID],
			pegged, common.IntGet(Ctx.handPoints))
		// eliminate cards from being picked
		limit(PLAYER)
		limit(AI)
		playerGO = !canPlay(hands[PLAYER])
		playerCard = nil
		if count == 31 {
			break
		}
		fallthrough // AI's turn
	case AI:
		// play until aiGO OR !playerGO
		for {
			aiGO = !canPlay(hands[AI])
			if aiGO {
				break
			}
			selectAiCard() // simulate a card select
			playCard(AI, &hands[AI][aiCard.ID])
			// played after PLAYER was GO?
			if playerGO && !bonusPoint {
				scorePoints(AI, 1, "bonus")
				bonusPoint = true
			}
			// award pegging points
			scorePegging(AI, hands[AI][aiCard.ID],
				pegged, common.IntGet(Ctx.handPoints))
			// eliminate cards from being picked
			limit(AI)
			limit(PLAYER)
			aiGO = !canPlay(hands[AI])

			// after AI play, PLAYER may not have a single card to play
			playerGO = !canPlay(hands[PLAYER])
			if playerGO { // disable all the unselected player's cards
				for i := range hands[PLAYER] {
					ci := playerStack.Objects[i].(*common.CardImage)
					if !hands[PLAYER][ci.ID].Selected {
						ci.SetDimmed(true)
					}
				}
			}
			// if PLAYER has a play, let it happen
			if !playerGO || aiGO {
				break
			}
		}
	}

	// transition to SCORE action
	if count == 31 || (aiGO && playerGO) {
		peggingButtonState(-1)
		duration = handScoreDuration
		scoreButton.Enable()
		return
	}

	return
}

/* -------------------------------------- */

func deal() {

	runSpiral = false

	count = 0
	pegged = nil
	duration = defaultDuration
	common.IntSet(Ctx.handPoints, 0)
	lanes[PLAYER].Reset()
	lanes[AI].Reset()
	cardPlayed(-1, ".x.") // clear all cards played

	if dealer == -1 { // low card is dealer
		d := common.CutForDeal([]string{Ctx.playerName, Ctx.aiName}, true)
		if d[0] == Ctx.playerName {
			dealer = PLAYER
		} else {
			dealer = AI
		}
	} else { // alternate on each deal
		dealer ^= 1
	}

	//dealer = PLAYER // testing

	if dealer == PLAYER {
		common.StringSet(Ctx.dealerName, Ctx.playerName)
	} else {
		common.StringSet(Ctx.dealerName, Ctx.aiName)
	}

	dealButton.Disable()
	scoreButton.Disable() // and the Score button

	// initialize the crib before discards
	crib = make([]common.Card, 0)
	ci := starterStack.Objects[0].(*common.CardImage)
	ci.Replace(Ctx.cardMap["reverse"].Resource)
	for i := range cribStack.Objects {
		ci := cribStack.Objects[i].(*common.CardImage)
		ci.Replace(Ctx.cardMap["reverse"].Resource)
	}

	// new shuffled deck of 52 cards
	cardDeck = common.NewDeck()
	common.Shuffle(cardDeck)

	hands, cardDeck = common.Deal(cardDeck, 2, 6) // player=hands[0], ai=hands[1]

	common.SortHandLow(hands[AI])
	for i := range hands[AI] { // save index for matching Image
		hands[AI][i].Index = i
		hands[AI][i].Selected = false
	}
	for i := range aiStack.Objects { // show(?) the AI hand
		ci := aiStack.Objects[i].(*common.CardImage)
		// Reset Disables tapping
		if Ctx.debug {
			ci.Reset(hands[AI][i].Resource)
		} else {
			ci.Reset(Ctx.cardMap["reverse"].Resource)
		}

	}

	common.SortHandLow(hands[PLAYER])
	for i := range hands[PLAYER] {
		hands[PLAYER][i].Index = i
		hands[PLAYER][i].Selected = false
	}
	for i := range 6 { // show the PLAYER hand
		ci := playerStack.Objects[i].(*common.CardImage)
		ci.Reset(hands[PLAYER][i].Resource)
	}
	playerStack.Refresh()

	discardAction(-1) // set up the tap handlers and Button state
	discardAction(AI)
}

/* -------------------------------------- */

var discards []*common.CardImage

func discardAction(who int) {

	// Tap handler for selecting/unselecting discard cards
	var selectCard = func(card *common.CardImage) {
		switch len(discards) {
		case 0:
			discards = append(discards, card)
			card.SetDim(true)
			// have 1 card. button = "Discard"
			discardButtonState(1)
			return
		case 1:
			// see if same as first
			if card.ID == discards[0].ID {
				discards[0].SetDim(false)
				discards = nil
				discardButtonState(0) // Discards?
				return
			}
			// have 2 now, need Button pressed
			discards = append(discards, card)
			card.SetDim(true)
			return
		}
		// have 2 cards
		// see if replace 1 of them
		if card.ID == discards[0].ID {
			discards[0].SetDim(false)
			discards = discards[1:]
			discardButtonState(1) // Discard
			return
		} else if card.ID == discards[1].ID {
			discards[1].SetDim(false)
			discards = discards[:1]
			discardButtonState(1) // Discard
			return
		}
		discardButtonState(2)
	}

	switch who {
	case -1: // prior to any button press → enter discard mode
		discardButtonState(0) // enable discard button
		// enable tapping on all 6 cards
		discards = make([]*common.CardImage, 0)
		for i := range playerStack.Objects {
			ci := playerStack.Objects[i].(*common.CardImage)
			ci.Reset(hands[PLAYER][i].Resource)
			ci.Enabled = true
			ci.OnTap = selectCard
		}
		return
	case AI: // choose AI discards
		var aiDiscard common.Cards
		_, aiDiscard = chooseDiscard(hands[AI], dealer == AI)
		crib = append(crib, aiDiscard...)
		// show "used" and unpickable
		hands[AI][aiDiscard[0].Index].Selected = true
		hands[AI][aiDiscard[1].Index].Selected = true
		aiStack.Objects[aiDiscard[0].Index].(*common.CardImage).Reset(Ctx.cardMap["reverse"].Resource)
		aiStack.Objects[aiDiscard[1].Index].(*common.CardImage).Reset(Ctx.cardMap["reverse"].Resource)
		aiStack.Refresh()
		return
	case PLAYER: // PLAYER clicked button - after selecting a Card
		// need 2 cards
		if len(discards) != 2 { // "wink"
			discardButtonState(2)
			return
		}
	}

	showCrib()    // show crib images - if debug
	drawStarter() // starter card draw and show

	// all discards in place, show
	for _, card := range discards {
		// disable and visually remove card from hand
		card.Reset(Ctx.cardMap["reverse"].Resource)
		// move card to crib
		crib = append(crib, hands[PLAYER][card.ID])
		hands[PLAYER][card.ID].Selected = true
	}

	// save the original hands for scoring
	keeps[PLAYER] = make([]common.Card, 0, 4)
	for _, card := range hands[PLAYER] {
		if !card.Selected {
			keeps[PLAYER] = append(keeps[PLAYER], card)
		}
	}
	keeps[AI] = make([]common.Card, 0, 4)
	for _, card := range hands[AI] {
		if !card.Selected {
			keeps[AI] = append(keeps[AI], card)
		}
	}

	discardButtonState(-1) // disable discard button
	peggingAction(-1)      // set up the pegging tap handlers and Button state
	// initiate pegging with AI, if dealer is PLAYER
	// otherwise, tap & Button
	if dealer == PLAYER {
		peggingAction(AI) // process a AI card play
	}
	// otherwise, taps & Button
}

func handToString(title string, who int, can bool) {
	log.Printf("%s: count %d, canPlay %t\n", title, common.IntGet(Ctx.handPoints), can)
	log.Printf("  ")
	for _, c := range hands[who] {
		log.Printf("%3s %t ", c.Icon, c.Selected)
	}
	log.Println()
}

// drawStarter sets the starter card, display starter in crib row, update cardDeck.
func drawStarter() {
	code := cardDeck[0]
	starter = Ctx.cardMap[code]
	cardDeck = cardDeck[1:]
	ci := starterStack.Objects[0].(*common.CardImage)
	ci.Replace(starter.Resource)
	if starter.Rank == common.Jack {
		scorePoints(dealer, 2, "heels")
	}
}

// show the crib
func showCrib() {
	if Ctx.debug { // show crib
		for i := range 4 {
			ci := cribStack.Objects[i].(*common.CardImage)
			ci.Reset(crib[ci.ID].Resource)
		}
	}
}

// show all hands
func showAll() {
	for i := range 4 {
		ci := cribStack.Objects[i].(*common.CardImage)
		ci.Reset(crib[ci.ID].Resource)
	}
	for i := range 6 {
		ci := playerStack.Objects[i].(*common.CardImage)
		ci.Reset(hands[PLAYER][ci.ID].Resource)
		ci = aiStack.Objects[i].(*common.CardImage)
		ci.Reset(hands[AI][ci.ID].Resource)
	}

}
