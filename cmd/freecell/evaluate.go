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
*/

const ReasonNoLegal = 0

const HintMoveSafeFoundation = 10
const HintMoveAceFoundation = 11
const HintMoveTwoFoundation = 12

const HintMoveSafeCardColumn = 20
const HintFreeCardColumn = 21
const HintKingCardColumn = 22

const HintEnableSeq = 30
const HintExtendAltSeq = 31
const HintParkFreeCell = 32

const HintMoveLowValue = 40

var Hints = map[int]string{
	ReasonNoLegal:          "No available legal move",
	HintMoveSafeFoundation: "Safe foundation move",
	HintMoveAceFoundation:  "Move Ace dest foundation",
	HintMoveTwoFoundation:  "Move Two dest foundation",
	HintMoveSafeCardColumn: "Safe Card Column move",
	HintFreeCardColumn:     "Free a card column",
	HintKingCardColumn:     "Place King in empty card column",
	HintEnableSeq:          "Move enables sequence",
	HintExtendAltSeq:       "Extend descending alternating sequence",
	HintParkFreeCell:       "Park in free cell",
	HintMoveLowValue:       "Low-value move",
}

func isValidSequence(gs *GameState, src *CellItem, seqLen int) bool {
	loc, col, idx := decodeCellItemId(src.id)
	// only Cards columns can have sequences
	if loc != Cards {
		return seqLen == 1
	}

	cards := gs.CardCells[col]
	// does sequence fit in column
	if idx+seqLen > len(cards) {
		return false
	}
	// insure descending alternating sequence
	for i := idx; i < idx+seqLen-1; i++ {
		a := cards[i]
		b := cards[i+1]
		if common.IsSameColor(a.card.Suit, b.card.Suit) {
			return false
		}
		if a.card.Rank != b.card.Rank+1 {
			return false
		}
	}
	return true
}

func canMoveSequence(gs *GameState, seq []CellItemId, dest CellItemId) bool {
	if len(seq) == 0 || dest == 0 {
		return false
	}

	// destination parameters
	dl, dc, _ := decodeCellItemId(dest)

	// decode the first card in sequence
	src := seq[0]
	sl, sc, sx := decodeCellItemId(src)

	// the source CellItem
	var srcItem *CellItem
	switch sl {
	case Free:
		srcItem = gs.FreeCells[sc]
	case Foundation:
		srcItem = gs.FoundationCells[sc]
	case Cards:
		srcItem = gs.CardCells[sc][sx]
	}

	if srcItem == nil {
		return false
	}

	switch dl { // destination
	case Free:
		// single cards to free
		if len(seq) != 1 {
			return false
		}
		return canMoveToFree(gs, gs.FreeCells[dc])

	case Foundation:
		// single cards to foundation
		if len(seq) != 1 {
			return false
		}
		return canMoveToFoundation(gs, srcItem, dest)

	case Cards:
		// multi-card sequence to cards
		return canMoveSequenceToCards(gs, srcItem, dc, len(seq))
	}

	return false
}

// must be Empty
func canMoveToFree(_ *GameState, dest *CellItem) bool {
	return isEmptyCard(dest)
}

// Empty or same suit 1 rank greater
// ♣ = 0, ♦ = 1, ♥ = 2, ♠ = 3
func canMoveToFoundation(gs *GameState, src *CellItem, dest CellItemId) (r bool) {
	dc := common.Suit(dest.column())
	if src.card.Suit != dc {
		return false
	}
	top := gs.FoundationCells[dc]
	if isEmptyCard(top) {
		r = src.card.Rank == 1
	} else {
		r = src.card.Rank == top.card.Rank+1
	}
	if Ctx.debug {
		fmt.Printf("canMoveToFoundation: src= %s (%s), top=%s [%t]\n",
			src.id, src.card.Icon, top.card.Icon, r)
	}

	return
}

// opposite color, 1 rank less.
func canMoveToCards(gs *GameState, src *CellItem, column int) bool {

	if Ctx.debug {
		fmt.Printf("canMoveToCards: src= %s (%s) cardColumn %d\n",
			src.id, src.card.Icon, column)
	}

	col := gs.CardCells[column]
	if len(col) == 0 { // only Kings
		return src.card.Rank == common.King
	}
	top := col[len(col)-1]
	if common.IsSameColor(top.card.Suit, src.card.Suit) {
		return false
	}
	// one rank lower
	return src.card.Rank+1 == top.card.Rank
}

func canMoveSequenceToCards(gs *GameState, src *CellItem, column int, seqLen int) bool {
	col := gs.CardCells[column]
	if len(col) == 0 { // catch moving a King.. sequence - only the King is OK
		return seqLen == 1 && src.card.Rank == common.King
	}
	if !isValidSequence(gs, src, seqLen) {
		return false
	}
	return canMoveToCards(gs, src, column)
}

//  ---------------------------------------  //

func isEmptyCard(ci *CellItem) bool {
	return ci.card.Icon == EmptyCardIcon
}

// count of empty Free Cells.
func freeCellCount(gs *GameState) int {
	n := 0
	for _, ci := range gs.FreeCells {
		if isEmptyCard(ci) {
			n++
		}
	}
	return n
}

/*  MICROSOFT rule - NOT NEEDED / WANTED in this version

in func canMoveSequenceToCards

	// can only move if room for each card to be moved
	// individually (to Free or Cards)
	if seqLen > moveableCellCount(gs) {
		return false
	}


// count of places to move - dumb rule :-)
// can only move if room for each card to be moved
// individually (to Free or Cards)
func moveableCellCount(gs *GameState) int {
	free := freeCellCount(gs)
	emptyCardCols := emptyCardsCount(gs)
	return (free + 1) * (1 << emptyCardCols)
}

*/
