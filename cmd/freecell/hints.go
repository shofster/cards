package main

import (
	"cards/common"
	"fmt"
)

/*

  File:    hints.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: generate a "best" move - similar to evaluate functions.
*/

var HintMap = map[int]string{
	200: "Extend Column / Free Cell",
	100: "Ace -> Foundation        ",
	90:  "2 -> Foundation          ",
	80:  "Card -> Foundation       ",
	75:  "Empty a Card Column      ",
	60:  "Extend a Card Column     ",
	50:  "Expose Card              ",
	40:  "King → Empty Column      ",
	20:  "Card -> Card Column      ",
	10:  "Card -> Free             ",
	1:   "Any Legal Move           ",
	2:   "No Legal Moves           ",
	0:   "Hints are OFF            ", // UI message when Hints are OFF
}

// Hint represents a Move with a score (reason)
type Hint struct {
	move  Move
	score int
}

// find the "best" :-) available move based on game state
func moveHint(gs *GameState) {
	if !common.BoolGet(hintOnOff) {
		return
	}
	moves := generateLegalMoves(gs)
	if len(moves) == 0 {
		common.StringSet(hintReason, HintMap[2])
		return
	}
	var best Hint
	for _, mv := range moves {
		score := scoreMove(gs, mv)
		if score > 0 {
			if Ctx.debug {
				fmt.Printf(" - score %d\n", score)
			}
			h := Hint{
				move:  mv,
				score: score,
			}
			if h.score > best.score {
				best = h
			}

		}
	}

	// either extend a Cards column OR Free (if any open)
	if best.score == 60 && freeCellCount(gs) > 0 {
		common.StringSet(hintReason, HintMap[200])
	} else {
		common.StringSet(hintReason, HintMap[best.score])
	}
}

func generateLegalMoves(gs *GameState) []Move {
	moves := []Move{}

	// Cards → anywhere
	for colIx, col := range gs.CardCells {
		if len(col) == 0 {
			continue
		}
		// top card
		top := col[len(col)-1]
		src := top.id
		// try all destinations
		for _, dest := range gs.getDestnations() {
			if canMoveSequence(gs, []CellItemId{src}, dest) {
				moves = append(moves, Move{
					src:  src,
					dest: dest,
					seq:  []CellItemId{src},
				})
			}
		}
		// sequences
		seq := gs.findSequence(colIx)
		if len(seq) > 1 {
			src := seq[0] // the top card of a moveable
			for _, dest := range gs.getDestnations() {
				if canMoveSequence(gs, seq, dest) {
					moves = append(moves, Move{
						src:  src,
						dest: dest,
						seq:  append([]CellItemId(nil), seq...),
					})
				}
			}
		}
	}

	// Free  → Cards or Foundation
	for _, cell := range gs.FreeCells {
		if cell != nil {
			src := cell.id
			for _, dest := range gs.getDestnations() {
				if canMoveSequence(gs, []CellItemId{src}, dest) {
					moves = append(moves, Move{
						src:  src,
						dest: dest,
						seq:  []CellItemId{src},
					})
				}
			}
		}
	}
	// Foundation → Cards
	for _, cell := range gs.FoundationCells {
		if cell != nil {
			src := cell.id
			for _, dest := range gs.getDestnations() {
				if canMoveSequence(gs, []CellItemId{src}, dest) {
					moves = append(moves, Move{
						src:  src,
						dest: dest,
						seq:  []CellItemId{src},
					})
				}
			}
		}
	}

	return moves
}

func scoreMove(gs *GameState, mv Move) int {

	// test using the top card of the "sequence" - may be length=1
	top := mv.seq[len(mv.seq)-1]
	item := getCellItem(gs, top.location(), top.column(), top.index())
	card := item.card
	if card.Icon == EmptyCardIcon {
		return 0
	}
	toLoc := mv.dest.location()
	fromLoc := mv.src.location()
	if Ctx.debug {
		fmt.Printf("scoreMove: card %s - from %s, to %s", card.Icon, mv.src, mv.dest)
	}

	// ace → foundation
	if card.Rank == 1 && toLoc == Foundation {
		return 100
	}
	// two → foundation
	if card.Rank == 2 && toLoc == Foundation && gs.foundationHasRank(card.Suit, 1) {
		return 90
	}
	// safe foundation move
	if toLoc == Foundation && gs.isSafeFoundationMove(card) {
		return 80
	}
	// expose useful card
	if gs.canExposeUsefulCard(mv) {
		return 50
	}
	// extend Cards sequence
	if toLoc == Cards {
		// destination column
		col := gs.CardCells[mv.dest.column()]
		// empty column → only Kings
		if len(col) == 0 {
			if card.Rank == common.King {
				return 40
			}
			return 0
		}
		// top card of destination
		d := col[len(col)-1].card
		// opposite color and descending rank
		if isRed(card.Suit) == isRed(d.Suit) {
			return 0
		}
		if card.Rank+1 != d.Rank {
			return 0
		}
		// score 60 if the source card is part of a real chain
		srcCol := gs.CardCells[mv.src.column()]
		srcIdx := mv.src.index()
		// check if the card is the top of the maximal chain
		isChainTop := true
		if srcIdx > 0 {
			prev := srcCol[srcIdx-1].card
			if isRed(prev.Suit) != isRed(card.Suit) && prev.Rank == card.Rank+1 {
				// there is a card above that continues the chain
				isChainTop = false
			}
		}

		if isChainTop {
			return 60
		}

		return 20

		return 20
		//return 60
	}
	// would free a column
	if fromLoc == Cards && len(gs.CardCells[mv.src.column()]) == 1 {
		return 75
	}
	// King → empty column
	if card.Rank == 13 && toLoc == Empty {
		return 40 // King-> empty column
	}
	// move to free cell
	if toLoc == Free && freeCellCount(gs) > 1 {
		return 10
	}

	return 1 // ant legal Move
}

func (gs *GameState) foundationHasRank(suit common.Suit, rank int) bool {
	top := gs.FoundationCells[suit]
	if top == nil {
		return false
	}
	return int(top.card.Rank) >= rank
}

func (gs *GameState) isSafeFoundationMove(card common.Card) bool {
	suit := card.Suit
	rank := int(card.Rank)

	// Ace and Two are always safe
	if rank <= 2 {
		return true
	}

	// opposite-color suits
	var opp1, opp2 common.Suit
	if isRed(suit) {
		opp1, opp2 = common.Spades, common.Clubs
	} else {
		opp1, opp2 = common.Hearts, common.Diamonds
	}
	// and rank >= rank-1
	r1 := gs.foundationRank(opp1)
	r2 := gs.foundationRank(opp2)
	return r1 >= rank-1 && r2 >= rank-1
}

func (gs *GameState) foundationRank(suit common.Suit) int {
	top := gs.FoundationCells[suit]
	if top == nil {
		return 0
	}
	return int(top.card.Rank)
}

func isRed(s common.Suit) bool {
	return s == common.Hearts || s == common.Diamonds
}

func (gs *GameState) canExposeUsefulCard(mv Move) bool {
	// only Cards expose anything
	if mv.src.location() != Cards {
		return false
	}

	col := gs.CardCells[mv.src.column()]
	ix := mv.src.index()

	// if  top card, nothing to exposed
	if !isBuried(col, ix) {
		return false
	}

	// the card that would be exposed
	exposedItem := col[ix+1]
	exposed := exposedItem.card

	// move to Foundation?
	dest := NewCellItemId(Foundation, int(exposed.Suit), 0)
	if canMoveToFoundation(gs, exposedItem, dest) {
		return true
	}
	// any Cards column?
	for c := 0; c < 8; c++ {
		dest := NewCellItemId(Cards, c, 0)
		if canMoveToCards(gs, exposedItem, dest.column()) {
			return true
		}
	}
	return false
}

func (gs *GameState) getDestnations() []CellItemId {
	dests := []CellItemId{}

	// Free
	for i, fc := range gs.FreeCells {
		if isEmptyCard(fc) {
			dests = append(dests, NewCellItemId(Free, i, 0))
		}
	}
	// Foundation
	for i := range gs.FoundationCells {
		dests = append(dests, NewCellItemId(Foundation, i, 0))
	}
	// Cards column
	for colIx, col := range gs.CardCells {
		if len(col) == 0 {
			// empty
			dests = append(dests, NewCellItemId(Cards, colIx, 0))
		} else {
			// or top
			topIx := len(col) - 1
			dests = append(dests, NewCellItemId(Cards, colIx, topIx))
		}
	}
	return dests
}

func isBuried(col []*CellItem, ix int) bool {
	return ix < len(col)-1
}

func (gs *GameState) findSequence(colIx int) []CellItemId {
	col := gs.CardCells[colIx]
	if len(col) == 0 {
		return nil
	}

	// start at the top card
	seq := []CellItemId{col[len(col)-1].id}

	// walk downward thru the sequence
	for i := len(col) - 2; i >= 0; i-- {
		s := seq[len(seq)-1]
		item := getCellItem(gs, s.location(), s.column(), s.index())
		a := item.card
		b := col[i]
		// alternating colors
		if common.IsSameColor(a.Suit, b.card.Suit) {
			break
		}
		// descending rank
		if a.Rank != b.card.Rank+1 {
			break
		}
		seq = append(seq, b.id)
	}
	return seq
}
