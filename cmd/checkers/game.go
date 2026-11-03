package main

import (
	"cards/common"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
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

var (
	// ChessPieces - fixed 8x8 squares with moveable pieces
	Board [64]*CheckerSquare
	// CardImages - fixed 8x8 slice and matching piece square
	CardImages       *fyne.Container
	cursor           *Cursor         // the cursor as a small image for a piece
	tracker          *common.Tracker // follow the mouse
	playSelect       *PlaySelect
	first            bool
	activePlayerSuit Suit
)

func initGame() {
	first = true
	activePlayerSuit = Red
	playerTurn.Importance = widget.HighImportance
	playSelect = NewPlaySelect()
	cursor = NewCursor()
	tracker = common.NewTracker()
	tracker.Event = func(ev *desktop.MouseEvent) {
		if ev != nil {
			FollowMouse(ev.Position)
		}
	}

	tracker.Click = func(down bool, ev *desktop.MouseEvent) {
		if first {
			resetPlayerPanels() // clear the captures
			first = false
			playerTurn.Refresh()
		}
		onClickEvent(activePlayerSuit, down, ev)
	}
	common.StringSet(cursor.cursorIcon, EmptyIcon)
}

func onClickEvent(suit Suit, down bool, ev *desktop.MouseEvent) {

	switch playSelect.onClick(down, ev) {

	case invalidSelect:
		return

	case sourceSelect:
		sourceId := RowCol2ID(playSelect.sourceRow, playSelect.sourceCol)
		// must be a square with active player's suit
		if Board[sourceId] == nil ||
			Board[sourceId].Suit != activePlayerSuit {
			playSelect.reset()
			return
		}
		borderImage(sourceId, true, common.BorderYellow)

	case destinationSelect:
		sourceId := RowCol2ID(
			playSelect.sourceRow,
			playSelect.sourceCol)
		destinationId := RowCol2ID(
			playSelect.destinationRow,
			playSelect.destinationCol)
		if sourceId == destinationId {
			playSelect.reset()
			borderImage(sourceId, false,
				common.BorderYellow)
			return
		}

		moves := LegalMoves(sourceId)
		// if this piece can capture,
		// discard all non-capturing moves
		hasCapture := false
		for _, move := range moves {
			if move.Capture {
				hasCapture = true
				break
			}
		}
		if hasCapture {
			var captureMoves []CheckerMove
			for _, move := range moves {
				if move.Capture {
					captureMoves = append(
						captureMoves, move)
				}
			}
			moves = captureMoves
		}
		valid := false
		var selectedMove CheckerMove
		for _, move := range moves {
			if move.DestinationID == destinationId {
				selectedMove = move
				valid = true
				break
			}
		}
		if !valid {
			playSelect.add(destinationId)
			playSelect.resetDestination()
			return
		}
		playSelect.reset()
		borderImage(sourceId, false, common.BorderYellow)
		applyMove(selectedMove)
		endTurn()
	}
}

func applyMove(move CheckerMove) {
	piece := Board[move.SourceID]
	applyCapture(piece, move.SourceID, move.DestinationID, move.CapturedID)

	currentID := move.DestinationID
	for _, capture := range move.Captures {
		applyCapture(piece, currentID, capture.DestinationID, capture.CapturedID)
		currentID = capture.DestinationID
	}
	if IsCrowned(piece) { // may end up in kings row - crown it
		updateSquareImage(currentID)
	}
	// next ...
	activePlayerSuit = otherPlayer(activePlayerSuit)
	updatePlayerIndicator()
}

func applyCapture(piece *CheckerSquare, sourceID int, destinationID int, capturedID int) {
	if capturedID >= 0 {
		captured := Board[capturedID]
		switch activePlayerSuit {
		case Red:
			player.addCapture(captured.Icon)
		case Black:
			bot.addCapture(captured.Icon)
		}
		Board[capturedID] = nil
		updateSquareImage(capturedID)
	}
	Board[destinationID] = piece
	Board[sourceID] = nil
	piece.Row, piece.Col = ID2RowCol(destinationID)
	updateSquareImage(sourceID)
	updateSquareImage(destinationID)
}

func endTurn() {
	if autoPlay && activePlayerSuit == Black {
		go func() {
			time.Sleep(time.Duration(Ctx.fadeTime) * time.Millisecond)
			fyne.Do(func() {
				CheckersMove()
			})
		}()
	}
}

func CheckersMove() {
	move, found := pickCaptureMove(activePlayerSuit)
	if !found {
		move, found = pickRandomMove(activePlayerSuit)
	}
	if found {
		applyMove(move)
	}
	endTurn()
}

func pickRandomMove(suit Suit) (CheckerMove, bool) {
	var allMoves []CheckerMove
	for id, piece := range Board {
		if piece == nil || piece.Suit != suit {
			continue
		}
		for _, m := range LegalMoves(id) {
			allMoves = append(allMoves, m)
		}
	}
	if len(allMoves) < 1 {
		return CheckerMove{}, false
	}
	return allMoves[rand.Intn(len(allMoves))], true
}

func pickCaptureMove(suit Suit) (CheckerMove, bool) {
	var captures []CheckerMove
	for id, piece := range Board {
		if piece == nil || piece.Suit != suit {
			continue
		}
		moves := LegalMoves(id)
		for _, m := range moves {
			if !m.Capture {
				continue
			}
			captures = append(captures, m)
		}
	}
	if len(captures) == 0 {
		return CheckerMove{}, false
	}
	return captures[rand.Intn(len(captures))], true
}
