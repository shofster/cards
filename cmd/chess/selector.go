package main

import (
	"time"

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
  Description: handle mouose clicks for source and destination selections.
*/

const (
	invalidSelect     = -1
	sourceSelect      = 0
	destinationSelect = 1
)

type PlaySelect struct {
	mode                           int   // progress of the selection
	selectTime                     int64 // time in MS when mouse button depresssed
	selectRow, selectCol           int   // current click location
	sourceRow, sourceCol           int   // final source location
	destinationRow, destinationCol int   // final destination
	badMoves                       []int
}

func NewPlaySelect() *PlaySelect {
	ps := &PlaySelect{}
	ps.badMoves = make([]int, 0)
	ps.reset()
	return ps
}

// resetable if the source is not valid
func (ps *PlaySelect) reset() {
	ps.mode = sourceSelect
	ps.selectTime = 0
	ps.selectRow = -1
	ps.selectCol = -1
	ps.sourceRow = -1
	ps.sourceCol = -1
	ps.destinationRow = -1
	ps.destinationCol = -1
	ps.resetBadMoves()

}

// reset the destination if not valid
func (ps *PlaySelect) resetDestination() {
	ps.mode = destinationSelect
	ps.selectTime = 0
	ps.selectRow = -1
	ps.selectCol = -1
	ps.destinationRow = -1
	ps.destinationCol = -1
}

// reset the down state for source or destination
func (ps *PlaySelect) resetDown() {
	ps.selectRow = -1
	ps.selectCol = -1
	ps.selectTime = 0
}

func (ps *PlaySelect) resetBadMoves() {
	for b := range ps.badMoves { // clear the bad ones
		dimSquare(ps.badMoves[b], false)
	}
	ps.badMoves = ps.badMoves[:0]
}

func (ps *PlaySelect) add(move int) {
	ps.badMoves = append(ps.badMoves, move)
	dimSquare(move, true)
}

func (ps *PlaySelect) onClick(down bool, ev *desktop.MouseEvent) int {
	if ev == nil {
		return invalidSelect
	}
	ok, row, col := mouse2RowCol(ev)
	if !ok {
		return invalidSelect
	}
	// start a down cycle for source or destination
	if down {
		ps.selectRow = row
		ps.selectCol = col
		ps.selectTime = time.Now().UnixMilli()
		return invalidSelect
	}
	// release time
	elapsed := time.Now().UnixMilli() - ps.selectTime
	if elapsed < 50 {
		return invalidSelect
	}
	// insure the up is in same as down
	if row != ps.selectRow || col != ps.selectCol {
		ps.resetDown()
		return invalidSelect
	}

	// have a completed down&up cycle (source or destination)
	switch ps.mode {
	case sourceSelect:
		ps.sourceRow = row
		ps.sourceCol = col
		ps.mode = destinationSelect
		ps.resetDown()
		return sourceSelect
	case destinationSelect:
		ps.destinationRow = row
		ps.destinationCol = col
		ps.mode = sourceSelect
		ps.resetDown()
		return destinationSelect
	}
	return invalidSelect
}

func mouse2RowCol(ev *desktop.MouseEvent) (ok bool, row, col int) {
	ix := int(ev.Position.X)
	iy := int(ev.Position.Y)
	if ix > 0 && ix < 8*SquareSize &&
		iy > 0 && iy < 8*SquareSize {
		col = ix / SquareSize
		row = iy / SquareSize
		ok = true
	}
	return
}

/*
func FindMove(id int) (*Move, bool) {
	for i := range legalMoves {
		if legalMoves[i].ID == id {
			return &legalMoves[i], true
		}
	}

	return nil, false
}

func IsLegalDestination(id int) bool {
	_, ok := FindMove(id)
	return ok
}

func TryMove(dstID int) (*Move, bool) {
	return FindMove(dstID)
}

*/
