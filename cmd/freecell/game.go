package main

import (
	"cards/common"

	"fyne.io/fyne/v2"
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

type Location int

const (
	EmptyColumn Location = iota
	CardColumn
	FreeCell
	Foundation
)

// a combined int for use as CardImage.ID
type CellItemId int

func cellItemId(loc Location, ix int) CellItemId {
	return CellItemId(int(loc<<8) | ix)
}
func cellItemLocation(id CellItemId) Location {
	return Location(id >> 8)
}
func cellItemIndex(id CellItemId) int {
	return int(id & 0xff)
}

type CellItem struct {
	id CellItemId
	common.Card
	buried bool
}

type GameState struct {
	CardColumns [8]common.Cards // 8 columns
	FreeCells   [4]*common.Card // nil = empty
	Foundation  [4]int          // highest per suit (0=empty, Ace -> King)
}

var (
	undos []Move
)

func deal() {
	undos = nil
	undoButton.Disable()
}
func undo() {

}

func freeValue(rank common.Rank) int {
	switch rank {
	case common.Ace:
		return 1
	case common.Jack, common.Queen, common.King:
		return 10
	default:
		return int(rank)
	}
}

func PushUndo(mv Move) {
	undos = append(undos, mv)
	undoButton.Enable()
}
func PopUndo() (Move, bool) {
	if !CanUndo() {
		return Move{}, false
	}
	n := len(undos) - 1
	mv := undos[n]
	undos = undos[:n]
	if len(undos) < 1 {
		undoButton.Disable()
	}
	return mv, true
}
func CanUndo() bool {
	return len(undos) > 0
}

func (gs *GameState) IsWin() bool {
	for _, v := range gs.Foundation {
		if v != 10 {
			return false
		}
	}
	return true
}

type BoardGeometry struct {
	FreeCellRect   [4]fyne.Position
	FoundationRect [4]fyne.Position
	ColumnX        [8]float32
	ColumnY        float32
}

/*
func HitTest(x, y float32) (region RegionType, col int, card int) {
	slotW := Ctx.cardSizeWidth
	slotH := Ctx.cardSizeHeight
	offset := Ctx.cardSizeHeight * 0.22

	// --- FreeCells / Foundations ---
	if y >= topRowY && y <= topRowY+slotH {
		idx := int(x / slotW)
		switch {
		case idx >= 0 && idx < 4:
			return RegionFreeCell, idx, -1
		case idx >= 4 && idx < 8:
			return RegionFoundation, idx - 4, -1
		}
	}

	// --- Tableau ---
	if y >= tableauY {
		col = int(x / slotW)
		if col < 0 || col >= 8 {
			return RegionNone, -1, -1
		}

		relY := y - tableauY
		card = int(relY / offset)

		if card >= len(cardStacks[col].Objects) {
			card = len(cardStacks[col].Objects) - 1
		}

		return RegionColumn, col, card
	}

	return RegionNone, -1, -1
}
*/
/*

func() {

	slotWidth  := Ctx.cardSizeWidth
	slotHeight := Ctx.cardSizeHeight

	// Y must be inside the top row
	if y >= topRowY && y <= topRowY+slotHeight {
		index := int(x / slotWidth)
		if index >= 0 && index < 4 {
			return RegionFreeCell, index
		}
		if index >= 4 && index < 8 {
			return RegionFoundation, index-4
		}
	}

	if y < tableauY {
		return RegionNone
	}
	col := int(x / Ctx.cardSizeWidth)
	if col < 0 || col >= 8 {
		return RegionNone
	}
	relY := y - tableauY
	index := int(relY / offset)

	if index >= len(columnCards) {
		index = len(columnCards) - 1
	}



}

*/
