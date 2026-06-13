package main

import "cards/common"

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
	HintMoveAceFoundation:  "Move Ace to foundation",
	HintMoveTwoFoundation:  "Move Two to foundation",
	HintMoveSafeCardColumn: "Safe Card Column move",
	HintFreeCardColumn:     "Free a card column",
	HintKingCardColumn:     "Place King in empty card column",
	HintEnableSeq:          "Move enables sequence",
	HintExtendAltSeq:       "Extend descending alternating sequence",
	HintParkFreeCell:       "Park in free cell",
	HintMoveLowValue:       "Low-value move",
}

// Hint represents a suggested move with a reason.
type Hint struct {
	Move   Move
	Reason int
	Score  int
}

// find the best :-) available hint based on game state.
func HintMove(gs *GameState) *Hint {
	moves := GenerateLegalMoves(gs)
	if len(moves) == 0 {
		return &Hint{
			Reason: ReasonNoLegal,
			Score:  0,
		}
	}
	var best *Hint
	for _, mv := range moves {
		score, reason := scoreMove(gs, mv)
		h := &Hint{
			Move:   mv,
			Reason: reason,
			Score:  score,
		}
		if best == nil || h.Score > best.Score {
			best = h
		}
	}
	return best
}

// create a score and reason.
func scoreMove(state *GameState, move Move) (int, int) {
	c := move.Card
	// Aces to Foundation - best move
	if c.Rank == common.Ace && move.To == Foundation {
		return 100, HintMoveAceFoundation
	}
	// Twos → Foundation (if Ace already there)
	if c.Rank == common.Two && move.To == Foundation && state.FoundationHas(freeValue(common.Ace), c.Suit) {
		return 90, HintMoveTwoFoundation
	}
	// Safe foundation move
	if move.To == Foundation && isSafeFoundationMove(state, c) {
		return 80, HintMoveSafeFoundation
	}
	// Free a column
	if move.To == EmptyColumn && move.From == CardColumn {
		return 75, HintFreeCardColumn
	}
	// continue a descending alternating run
	if move.To == CardColumn && isGoodSequenceExtension(state, move) {
		return 60, HintExtendAltSeq
	}
	// enqble a sequence
	if enablesSequence(state, move) {
		return 50, HintEnableSeq
	}
	// King to empty column
	if c.Rank == common.King && move.To == EmptyColumn {
		return 40, HintKingCardColumn
	}
	// any legal move
	if move.To == CardColumn {
		return 20, HintMoveSafeCardColumn
	}
	// park in a free cell - if not last one
	if move.To == FreeCell && state.FreeCellCount() > 1 {
		return 10, HintParkFreeCell
	}
	// any
	return 1, HintMoveLowValue
}

// safe if no lower card of same suit is buried
func isSafeFoundationMove(gs *GameState, c common.Card) bool {
	return !gs.IsLowerCardBuried(c)
}

// check alternating color + rank-1
func isGoodSequenceExtension(gs *GameState, mv Move) bool {
	return gs.ValidSequenceExtension(mv)
}

// exposes a needed card?
func enablesSequence(gs *GameState, mv Move) bool {
	return gs.ExposesUsefulCard(mv)
}

// total count of emptty Free Cells
func (gs *GameState) FreeCellCount() int {
	n := 0
	for _, c := range gs.FreeCells {
		if c == nil {
			n++
		}
	}
	return n
}

// check if a lower rank card is buried
func (gs *GameState) IsLowerCardBuried(c common.Card) bool {
	target := freeValue(c.Rank) - 1
	if target < 1 {
		return false
	}
	for _, col := range gs.CardColumns {
		for _, card := range col {
			if card.Suit == c.Suit && freeValue(card.Rank) == target {
				return true
			}
		}
	}
	return false
}

// chack if a sequence extension
func (gs *GameState) ValidSequenceExtension(mv Move) bool {
	if mv.To != CardColumn {
		return false
	}
	dest := gs.CardColumns[mv.ToIndex]
	if len(dest) == 0 {
		return false
	}
	top := dest[len(dest)-1]
	// alternating colors
	if sameColor(top.Suit, mv.Card.Suit) {
		return false
	}
	// descending by 1
	return mv.Card.Rank+1 == top.Rank
}

// check if would expose a card
func (gs *GameState) ExposesUsefulCard(mv Move) bool {
	if mv.From != CardColumn {
		return false
	}
	col := gs.CardColumns[mv.FromIndex]
	if len(col) < 2 {
		return false
	}
	under := col[len(col)-2]
	// useful if the next foundation card
	return freeValue(under.Rank) == gs.Foundation[under.Suit]+1
}

// check if same color ("♣" and "♠" - black, "♦" and "♥" - red)
func sameColor(a, b common.Suit) bool {
	return (a == common.Hearts || a == common.Diamonds) ==
		(b == common.Hearts || b == common.Diamonds)
}
func (gs *GameState) FoundationHas(value int, suit common.Suit) bool {
	return gs.Foundation[suit] >= value
}
