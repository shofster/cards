package main

import (
	"cards/common"
)

/*

  File:    game.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: main game for FreeCell
*/

/*

	// Card relationships
	// CardItem(w/Card):CardImage(w/Resource):Geometry(position and size)
	//   related by index and/or CellItemId


GameState contains the CellItems for each of the cards:
   for Free and Foundation [4]
   for Cards [8][n]
      type CellItem struct {
		id         CellItemId
		card       common.Card
			type Card struct {
				rank     Rank
				suit     Suit
				icon     string
				resource fyne.Resource
				......  other
		....  other

--------------------

FreeCells[col] is always a single CellItem with:
    Location = Free
    card = emptyCard() when empty
        ... and card.Icon = "empty" for emptiness testing

  The CardImage stack (freeStack) for a FreeCell has only one object:
    When empty:
        ID = (Empty, 0, 0)
        resource = emptyCard().Icon
    When filled:
        ID = (Free, col, 0)
        resource = the real card.Icon


Foundation[col] is always a single CellItem with:
    Location = Foundation
    card = emptyCard() when empty
        ... and card.Icon = "empty" for emptiness testing

  The CardImage stack (foundationStack) for a Foundation has only one object:
    When empty:
        resource = emptyCard().Icon
    When filled:
        resource = the real card.Icon


Cards[col] is a []CellItem with:
    Location = Cards
    When empty:
        slice is nil or length 0
        CardImage stack has a single Empty card object
    When cards are added/removed:
        update the Cards[col] slice
        rebuild the CardImage stack from the slice

  The CardImage stack (cardStacks) for each Card Column has:
    When empty: 1 Object
        ID = (Empty, 0, 0)
        resource = emptyCard().Icon
    When filled: same len as the []*CellItem
        ID = (Cards, col, 0)
        resource = the real card.Icon


*/

func deal() {

	undos = nil
	undoButton.Disable()

	cardDeck = common.NewDeck() // []string
	common.Shuffle(cardDeck)
	var col = 0

	for col = 0; col < 4; col++ {
		gameState.FreeCells[col] = NewCellItem(Free, col, 0, emptyCard())
		gameState.FoundationCells[col] = NewCellItem(Foundation, col, 0, emptyCard())
	}
	resetFreeCellImages(&gameState, freeGrid)
	resetFoundationCellImages(&gameState, foundationGrid)

	// Build Cards and Image relationships
	dx := 0 // distribute the 52 cards
	for col = 0; col < 8; col++ {
		gameState.CardCells[col] = make([]*CellItem, 0)
		n := 7
		if col > 3 {
			n = 6
		}
		for ix := 0; ix < n; ix++ {
			icon := cardDeck[dx]
			rank, suit := common.RankSuit(icon)
			if rank == common.Ace {
				rank = common.Two - 1
			}
			card := common.Card{
				Rank:     rank,
				Suit:     suit,
				Icon:     icon,
				Resource: Ctx.cardMap[icon].Resource,
			}
			gameState.CardCells[col] = append(gameState.CardCells[col],
				NewCellItem(Cards, col, ix, card))
			dx++
		}
	}
	for col = 0; col < 8; col++ {
		resetCardColumnImages(col, &gameState)
	}
	tracker.Event = mouseMove  // for testing use trackMouse
	tracker.Click = mouseClick // for testing use nil
}

func pushUndo(mv Move) {
	undos = nil
	undos = append(undos, mv)
	undoButton.Enable()
}

func popUndo() {
	mv := undos[0]
	undos = nil

	// mv.src/mv.dest already set up for undo in performMove
	performMove(&gameState, mv.seq, mv.src, mv.dest, true)

	undoButton.Disable()
}

type GameState struct {
	FreeCells       [4]*CellItem   // 4 columns - 1 deep
	FoundationCells [4]*CellItem   // 4 columns - 1 deep
	CardCells       [8][]*CellItem // 8 columns - n deep, overlapped
}

type Move struct {
	src  CellItemId
	dest CellItemId
	seq  []CellItemId // len = 1 for single card
}

var (
	gameState GameState
	cardDeck  []string // 52 Card Icons (shuffled)
	geometry  Geometry // recomputed on every mouse movement
	undos     []Move   // a single Move for undo-ing
)

func getCellItem(gs *GameState, loc Location, col, idx int) *CellItem {
	switch loc {
	case Free:
		return gs.FreeCells[col]
	case Foundation:
		return gs.FoundationCells[col]
	case Cards:
		return gs.CardCells[col][idx]
	default:
		return nil
	}
}
