package main

import (
	"cards/common"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

/*

  File:    selector.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Follow the mouse - select source, destination - move cards.
*/

type InteractionPhase int

const (
	MouseDown = 1
	Dragging  = 2
)

type Interaction struct {
	Phase       InteractionPhase
	SourceCell  CellItemId   // where mouseDown happened
	HoverCell   CellItemId   // cell currently under mouse
	InvalidCell CellItemId   // a hovered cell that is invalid for destination
	Seq         []CellItemId // selected sequence (1 or more cards)
	DragCard    *common.Card // card shown next to cursor
}

var (
	mouseCell   CellItemId
	interaction Interaction
)

func mouseClick(down bool, ev *desktop.MouseEvent) {
	if ev == nil {
		return
	}
	if down {
		common.StringSet(clickDir, "\u2193")
		mouseDown(ev.AbsolutePosition)
	} else {
		common.StringSet(clickDir, "\u2191")
		mouseUp()
	}
}

// recompute the UI geometry
// and place the mouse pointer into a valid UI element
func mousePosition(mp fyne.Position) {
	common.IntSet(mouseX, int(mp.X))
	common.IntSet(mouseY, int(mp.Y))
	buildGeometry()
	mouseCell = findCellItem(mp)
}

// called on every mouse movement.
func mouseMove(ev *desktop.MouseEvent) {
	if ev == nil {
		return
	}

	// for same cell test
	oldCell := mouseCell

	// update geometry + mouseCell and cursor
	mousePosition(ev.AbsolutePosition)

	// start drag?
	if interaction.Phase == MouseDown {
		interaction.Phase = Dragging
	}
	if interaction.Phase != Dragging {
		return
	}

	cursor.Show(true)
	cursor.Move(ev.AbsolutePosition)

	// ignore if haven't moved outside of last Cell
	if mouseCell == oldCell {
		return
	}

	// cleal old highlite
	if interaction.HoverCell != 0 {
		if ci := cardImageFromCellId(interaction.HoverCell); ci != nil {
			ci.SetHighlite(false)
		}
	}
	if interaction.InvalidCell != 0 { // and invalid
		if ci := cardImageFromCellId(interaction.InvalidCell); ci != nil {
			ci.SetDim(false)
		}
		interaction.InvalidCell = 0
	}

	// update hover cell
	interaction.HoverCell = mouseCell
	// ignore source cell
	if mouseCell == interaction.SourceCell {
		return
	}
	// nothing under mouse
	if mouseCell == 0 {
		return
	}

	// new highlight/dim
	if canMoveSequence(&gameState, interaction.Seq, mouseCell) {
		if ci := cardImageFromCellId(mouseCell); ci != nil {
			ci.SetHighlite(true)
		}
	} else {
		interaction.InvalidCell = mouseCell
		if ci := cardImageFromCellId(mouseCell); ci != nil {
			ci.SetDim(true)
		}
	}

	// follow the mouse with the card Icon
	cursor.Move(ev.AbsolutePosition)
}

func mouseDown(mp fyne.Position) {
	// updates geometry + mouseCell
	mousePosition(mp)

	// nothing under mouse → reset interaction
	if mouseCell == 0 {
		interaction = Interaction{}
		cursor.Show(false)
		return
	}

	var invalidDown = func() {
		interaction.InvalidCell = mouseCell
		if ci := cardImageFromCellId(mouseCell); ci != nil {
			ci.SetDim(true)
		}
		return
	}

	// might be invalid - let me count the ways
	card := cardFromCellId(mouseCell)

	if Ctx.debug {
		if card == nil { // empty column
			fmt.Printf("mouseDown: card is nil\n")
		} else { // empty Free or Foundation
			fmt.Printf("mouseDown: cell %s %s \n", mouseCell, card.Icon)
		}
	}
	if card == nil { // Cards
		invalidDown()
		return
	} else {
		icon := card.Icon // may still be Empty
		if mouseCell.location() == Empty || icon == EmptyCardIcon ||
			mouseCell.location() == Foundation {
			invalidDown()
			return
		}
		common.StringSet(cursor.cursorCard, fmt.Sprintf(".  %s", icon))
		common.StringSet(mouseIcon, fmt.Sprintf("%3s", icon))
	}

	cursor.Move(mp)
	cursor.Show(true)
	interaction = Interaction{}

	// build the movable sequence starting on this cell
	seq, ok := findMovableSequence(&gameState, mouseCell)
	if !ok {
		interaction = Interaction{}
		return
	}
	// Begin mouse interaction
	interaction.Phase = MouseDown
	interaction.SourceCell = mouseCell
	interaction.HoverCell = mouseCell
	interaction.Seq = seq
	interaction.DragCard = cardFromCellId(seq[0])
	dimSequence(seq)
}

func mouseUp() {
	// hide cursor
	cursor.Show(false)

	// clear dim on the last invalid cell
	if interaction.InvalidCell != 0 {
		if ci := cardImageFromCellId(interaction.InvalidCell); ci != nil {
			ci.SetDim(false)
		}
		interaction.InvalidCell = 0
	}
	// clear highlight on hover cell
	if interaction.HoverCell != 0 {
		if ci := cardImageFromCellId(interaction.HoverCell); ci != nil {
			ci.SetHighlite(false)
		}
		interaction.HoverCell = 0
	}
	// clear sequence highlight
	clearSequenceHighlight(interaction.Seq)
	// clear previous hover highlight
	//   except for the source cell
	if interaction.HoverCell != 0 && interaction.HoverCell != interaction.SourceCell {
		if ci := cardImageFromCellId(interaction.HoverCell); ci != nil {
			ci.SetHighlite(false)
		}
	}
	// clear highlight on the cell currently under the mouse
	if mouseCell != 0 {
		if ci := cardImageFromCellId(mouseCell); ci != nil {
			ci.SetHighlite(false)
		}
	}

	// not dragging
	if interaction.Phase != Dragging {
		interaction = Interaction{}
		return
	}

	// move using the last known mouseCell
	if canMoveSequence(&gameState, interaction.Seq, mouseCell) {
		performMove(&gameState, interaction.Seq, interaction.SourceCell, mouseCell, false)
	}

	// initial state
	interaction = Interaction{}
}

// move workhorse
func performMove(gs *GameState, seq []CellItemId, src CellItemId, dest CellItemId,
	isRedo bool) {

	if len(seq) == 0 {
		return
	}

	sl, sc, si := decodeCellItemId(src)
	dl, dc, _ := decodeCellItemId(dest)

	s := getCellItem(gs, sl, sc, si)
	if Ctx.debug {
		log.Printf("performMove: src= %s (%s), dest=%s\n",
			src, s.card.Icon, dest)
	}

	var moving []*CellItem

	switch sl {
	case Free:
		moving = []*CellItem{gs.FreeCells[sc]}
		gs.FreeCells[sc] = NewCellItem(Free, sc, 0, emptyCard())
	case Foundation:
		moving = []*CellItem{gs.FoundationCells[sc]}
		gs.FoundationCells[sc] = NewCellItem(Foundation, sc, 0, emptyCard())
	case Cards:
		col := gs.CardCells[sc]

		// extract moving cards
		moving = make([]*CellItem, len(seq))
		for i := range seq {
			_, _, idx := decodeCellItemId(seq[i])
			moving[i] = col[idx]
		}

		// build new column without moved cards
		newCol := make([]*CellItem, 0, len(col)-len(seq))
		for i := range col {
			keep := true
			for _, ci := range seq {
				_, _, idx := decodeCellItemId(ci)
				if i == idx {
					keep = false
					break
				}
			}
			if keep {
				newCol = append(newCol, col[i])
			}
		}

		for i := range newCol {
			newCol[i].id = NewCellItemId(Cards, sc, i)
		}

		gs.CardCells[sc] = newCol
	default:
	}

	switch dl {
	case Free:
		// build a fresh CellItem for the Free cell
		card := moving[0].card
		gs.FreeCells[dc] = NewCellItem(Free, dc, 0, card)
		updateFreeCell(dc, gs)
	case Foundation:
		// build a fresh CellItem for the Foundation
		card := moving[0].card
		gs.FoundationCells[dc] = NewCellItem(Foundation, dc, 0, card)
		updateFoundationCell(dc, gs)

	case Cards:
		base := len(gs.CardCells[dc])
		for i, ci := range moving {
			ci.id = NewCellItemId(Cards, dc, base+i)
		}
		gs.CardCells[dc] = append(gs.CardCells[dc], moving...)
	default:
	}

	// ? - add current move to undos
	if !isRedo {
		// skip undo for moves INTO Foundation
		sl, _, _ := decodeCellItemId(dest)
		if sl == Foundation {
			// disable button
			undoButton.Disable()
		} else {
			// normal undo recording
			undoSeq := make([]CellItemId, len(moving))
			for i, ci := range moving {
				undoSeq[i] = ci.id
			}
			pushUndo(Move{
				src:  dest,
				dest: src,
				seq:  undoSeq,
			})
		}
	}

	common.IntSet(playCount, common.IntGet(playCount)+1)
	moveHint(gs) // hint for next move

	// bulk re-create of the Cell images
	resetFreeCellImages(gs, freeGrid)
	resetFoundationCellImages(gs, foundationGrid)
	if src.location() == Cards && dest.location() == Cards && src.column() == dest.column() {
		resetCardColumnImages(src.column(), gs)
	} else {
		if src.location() == Cards {
			resetCardColumnImages(src.column(), gs)
		}
		if dest.location() == Cards {
			resetCardColumnImages(dest.column(), gs)
		}
	}
	interaction = Interaction{}
}

func cardImageFromCellId(ci CellItemId) *common.CardImage {
	if ci != 0 {
		l, c, x := decodeCellItemId(ci)
		switch l {
		case Free:
			return freeGrid.Objects[c].(*common.CardImage)
		case Foundation:
			return foundationGrid.Objects[c].(*common.CardImage)
		case Cards:
			return cardStacks[c].Objects[x].(*common.CardImage)
		default:
		}
	}
	return nil
}

func cardFromCellId(ci CellItemId) *common.Card {
	l, c, x := decodeCellItemId(ci)
	if ci != 0 {
		switch l { // find the CellItem
		case Free: // CellItems
			return &gameState.FreeCells[c].card
		case Foundation:
			return &gameState.FoundationCells[c].card
		case Cards: // may be place holder
			if len(gameState.CardCells[c]) > 0 {
				return &gameState.CardCells[c][x].card
			}
		default:
		}
	}
	return nil
}

func findMovableSequence(gs *GameState, start CellItemId) ([]CellItemId, bool) {
	l, col, idx := decodeCellItemId(start)
	if l != Cards {
		return []CellItemId{start}, true
	}

	cells := gs.CardCells[col]
	seq := []CellItemId{start}
	valid := true

	for i := idx + 1; i < len(cells); i++ {
		above := &cells[i-1].card
		below := &cells[i].card

		// append so know the full sequence
		seq = append(seq, cells[i].id)

		// insure alternating color
		if common.IsSameColor(above.Suit, below.Suit) {
			valid = false
			break
		}

		// and descending rank
		if above.Rank != below.Rank+1 {
			valid = false
			break
		}
	}

	return seq, valid
}

func dimSequence(seq []CellItemId) {
	for _, ci := range seq {
		ci := cardImageFromCellId(ci)
		if ci != nil {
			ci.SetHighlite(true)
		}
	}
}
func clearSequenceHighlight(seq []CellItemId) {
	for _, ci := range seq {
		ci := cardImageFromCellId(ci)
		if ci != nil {
			ci.SetHighlite(false)
		}
	}
}

// return the CellItem that the mouse is "inside"
func findCellItem(mouse fyne.Position) CellItemId {
	for col := 0; col < 4; col++ {
		if isInside(geometry.FreeGrid[col], mouse) {
			return NewCellItemId(Free, col, 0)
		}
		if isInside(geometry.FoundationGrid[col], mouse) {
			return NewCellItemId(Foundation, col, 0)
		}
	}
	for col := 0; col < 8; col++ {
		ix := cardAt(geometry.CardGrids[col], mouse)
		if ix >= 0 {
			return NewCellItemId(Cards, col, ix)
		}
	}
	return NewCellItemId(Empty, 0, 0)
}

// tracking echo - used for original testing only
func trackMouse(ev *desktop.MouseEvent) {
	if ev == nil {
		common.StringSet(mouseIcon, "   ")
		return
	}
	mousePosition(ev.AbsolutePosition)
	common.StringSet(mouseIcon, "")
	if mouseCell == 0 {
		cursor.Show(false)
	} else {
		card := cardFromCellId(mouseCell)
		if card != nil {
			icon := card.Icon
			common.StringSet(cursor.cursorCard, fmt.Sprintf(".  %s", icon))
			common.StringSet(mouseIcon, fmt.Sprintf("%3s", icon))
			cursor.Show(true)
		}
	}
	cursor.Move(ev.AbsolutePosition)
}

/*

to add keyboard - someday

1. Arrow keys move a “focus rectangle”
Up/down moves within a column
Left/right moves between columns
Home jumps to top of column
End jumps to bottom of column
2. Space or Enter selects the card/sequence
This is the keyboard equivalent of mouse‑down.
You’d dim the sequence exactly the same way.
3. Arrow keys move the “target highlight”
This is the keyboard equivalent of hover.
4. Enter confirms the move
This is the keyboard equivalent of mouse‑up.
5. Escape cancels the drag

Tab and Shift+Tab to move around board from Free<->Foundation<->Cards
SPACE or ENTER
Selects the focused card(s)
*/
