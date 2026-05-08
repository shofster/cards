package main

import (
	"fmt"
	"games/common"
	"time"

	"fyne.io/fyne/v2/widget"
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
var hands [][]common.Card
var trickStates *common.CardStrip

var southManager *HandManager
var northManager *HandManager
var eastManager *HandManager
var westManager *HandManager
var managers [4]*HandManager

var commits [4]*widget.Button

type gameType struct {
	leadPlayer  Player      // initially South, then the trick winner
	leadSuit    common.Suit // initially NoSuit, set after first lead
	player      Player      // the current player : leadPLayer, next, next, .. back to leader
	trickNumber int         // trick # completed
}

var g gameType

func initGame() {
	g = gameType{}
	cardMap = common.CardMap() // icon vs Card
	cardDeck = common.NewDeck()
}

func deal() {

	common.Shuffle(cardDeck)

	// 4 hands of 13 cards each
	hands, _ = common.Deal(cardDeck, 4, 13)
	common.SortHand(hands[South])
	common.SortHand(hands[North])
	if common.BoolGet(showEast) {
		common.SortHand(hands[East])
	}
	if common.BoolGet(showWest) {
		common.SortHand(hands[West])
	}

	southManager = NewSouthHandManager((hands[South]))
	southManager.Deal()
	managers[0] = southManager
	commits[0] = southCommit
	westManager = NewWestHandManager((hands[West]))
	westManager.Deal()
	managers[1] = westManager
	commits[1] = westCommit
	northManager = NewNorthHandManager((hands[North]))
	northManager.Deal()
	managers[2] = northManager
	commits[2] = northCommit
	eastManager = NewEastHandManager((hands[East]))
	eastManager.Deal()
	managers[3] = eastManager
	commits[3] = eastCommit

	southCommit.OnTapped = func() { commitTrick() }
	westCommit.OnTapped = func() { commitTrick() }
	northCommit.OnTapped = func() { commitTrick() }
	eastCommit.OnTapped = func() { commitTrick() }

	trickStates.Empty = "*"
	trickStates.Reset()
	scoreCard.Reset()
	g.trickNumber = 0
}

func play() {
	common.StringSet(dealResult, "  Playing...")
	g.leadPlayer = South
	g.player = South
	g.leadSuit = common.NoSuit
	common.IntSet(contractsBid, common.IntGet(contractsBid)+1)
	trickPlayer()
}

func trickPlayer() {

	common.StringSet(nextAction, fmt.Sprintf("%s - it's your turn",
		playerToPosition(g.player)))

	// the position manager enables its commit when a card is selected
	for i := 0; i < 4; i++ {
		commits[i].Disable()
	}

	managers[g.player].Play(g.player, g.leadPlayer,
		suitToSuit(common.StringGet(contractSuit)),
		g.leadSuit)
}

func commitTrick() {
	// Capture the player at the moment of play
	p := g.player

	// Set lead suit when the leader plays the first card
	if p == g.leadPlayer {
		g.leadSuit = managers[p].cards[managers[p].trickSelection].Suit
	}

	// Mark card as used in this deal
	managers[p].CommitSelection()

	// Scorecard update
	scoreCard.SetCard(
		g.trickNumber,
		playerToPosition(p)[:1],
		managers[p].cards[managers[p].trickSelection].Icon,
	)

	// Inline deterministic next-player logic
	startNextPlayer := func(p Player) {
		// If this player is the last in the trick, finish it
		if int(p) == (int(g.leadPlayer)+3)%4 {
			finishTrick()
			return
		}

		// Otherwise advance to next player
		g.player = Player((int(p) + 1) % 4)
		trickPlayer()
	}

	// East/West delay path
	if p == East || p == West {
		go func(p Player) {
			time.Sleep(1500 * time.Millisecond)
			startNextPlayer(p)
		}(p)
		return
	}

	// Immediate path
	startNextPlayer(p)
}

func finishTrick() {

	// build cards in order of position led
	trickCards := buildTrickCards(g.leadPlayer)
	relativeWinner := winnerOfTrick(
		trickCards,
		suitToSuit(common.StringGet(contractSuit)),
		0, // leader is [0] in trickCards
	)
	winner := (g.leadPlayer + relativeWinner) % 4

	trickStates.SetCard(g.trickNumber, playerToPosition(winner)[:1])
	// fmt.Sprintf("%d", winner))

	//done with the trick - disable ALL commits
	for i := 0; i < 4; i++ {
		commits[i].Disable()
		managers[i].trickSelection = -1
	}

	g.trickNumber++
	if g.trickNumber == 13 {
		winners := make([]string, 13)
		for i := 0; i < 13; i++ {
			v, _ := trickStates.States.GetValue(i)
			winners[i] = v
		}
		bid, above, below := dealScores(
			winners,
			common.StringGet(contractLevel),
			common.StringGet(contractSuit))
		s := "  Bid Made"
		if bid < 0 {
			s = fmt.Sprintf("  Lost Bid by %d", bid)
		} else if bid > 0 {
			common.IntSet(contractsWon, common.IntGet(contractsWon)+1)
			s += fmt.Sprintf(", + %d", bid)
		}
		common.StringSet(dealResult, s)
		common.IntSet(aboveLine, common.IntGet(aboveLine)+above)
		common.IntSet(belowLine, common.IntGet(belowLine)+below)
		southManager.RevealAll()
		westManager.RevealAll()
		northManager.RevealAll()
		eastManager.RevealAll()
		resetPlayedCards()
		return
	}

	var startNextTrick = func() {
		// winner becomes leader of next trick
		g.leadPlayer = winner
		g.player = winner
		g.leadSuit = common.NoSuit
		resetPlayedCards() // reset played cards
		for i := 0; i < 4; i++ {
			commits[i].Disable()
			managers[i].trickSelection = -1
		}

		trickPlayer()
	}
	if winner == East || winner == West {
		go func() {
			time.Sleep(20 * time.Millisecond)
			// Continue the Deal flow after the delay
			startNextTrick()
		}()
		return
	}
	startNextTrick()

}

// build the played cards in order played
func buildTrickCards(leader Player) [4]common.Card {
	var out [4]common.Card
	l := leader
	for i := 0; i < 4; i++ {
		out[i] = managers[l].cards[managers[l].trickSelection]
		l = (l + 1) % 4
	}
	return out
}
