package main

import (
	"cards/assets/chess"
	"cards/common"
	"math/rand"

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
	Board [64]*ChessSquare
	// CardImages - fixed 8x8 slice and matching piece square
	CardImages       *fyne.Container
	cursor           *Cursor         // the cursor as a UNICODE char for a piece
	tracker          *common.Tracker // follow the mouse
	playSelect       *PlaySelect
	first            bool
	activePlayerSuit Suit
	//badMoves      []int
	blackKingWarn int
	whiteKingWarn int
)

func initGame() {
	buildBoard()
	first = true
	blackKingWarn = -1
	whiteKingWarn = -1
	activePlayerSuit = White
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
			resetPlayerPanels()
			first = false
			playerTurn.Refresh()
		}
		onClickEvent(activePlayerSuit, down, ev)
	}
	common.StringSet(cursor.cursorIcon, "empty")
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
		sourceId := RowCol2ID(playSelect.sourceRow, playSelect.sourceCol)
		destinationId := RowCol2ID(playSelect.destinationRow, playSelect.destinationCol)
		if sourceId == destinationId { // clicked on same square
			playSelect.reset()
			borderImage(sourceId, false, common.BorderYellow)
			return
		}
		// validate the move.
		// dim the image if invalid and reset the Destination
		// leave the dims until get a good one, then hide them
		moves := LegalMoves(sourceId)
		valid := false
		capture := false
		// any matches with the destination?
		for _, move := range moves {
			if move.ID == destinationId {
				capture = move.Capture
				valid = true
				break
			}
		}
		if !valid { // save the invalid and mark it
			playSelect.add(destinationId)
			playSelect.resetDestination()
			return
		}

		// reset for next move
		playSelect.reset()
		borderImage(sourceId, false, common.BorderYellow)

		// process the valid move
		applyMove(sourceId, destinationId, capture)
		endTurn()

	}
}

func applyMove(sourceId, destinationId int, capture bool) {
	if capture {
		switch activePlayerSuit {
		case White:
			player.addCapture(chess.ChessUnicode[Board[destinationId].Icon])
		case Black:
			bot.addCapture(chess.ChessUnicode[Board[destinationId].Icon])
		}
	}
	swapBoardPieceImages(sourceId, destinationId)
	Board[destinationId] = Board[sourceId]
	Board[sourceId] = nil
	Board[destinationId].Row, Board[destinationId].Col = ID2RowCol(destinationId)
	activePlayerSuit = otherPlayer(activePlayerSuit)
	updatePlayerIndicator()
}

func endTurn() {
	kingWarn()
	if autoPlay && activePlayerSuit == Black {
		ChessterMove()
	}
}

func ChessterMove() {
	move, found := pickCaptureMove(activePlayerSuit)
	if !found {
		move = pickRandomMove(activePlayerSuit)
	}
	applyMove(move.SourceID, move.DestinationID, move.Capture)
	endTurn()
}

func pickRandomMove(suit Suit) ChessMove {
	var allMoves []ChessMove
	for id, piece := range Board {
		if piece == nil || piece.Suit != suit {
			continue
		}
		for _, m := range LegalMoves(id) {
			allMoves = append(allMoves, ChessMove{
				SourceID:      id,
				DestinationID: m.ID,
				Capture:       m.Capture,
			})
		}
	}
	if len(allMoves) < 1 {
		panic("no moves")
	}
	return allMoves[rand.Intn(len(allMoves))]
}

func pickCaptureMove(suit Suit) (ChessMove, bool) {
	var captures []ChessMove
	for id, piece := range Board {
		if piece == nil || piece.Suit != suit {
			continue
		}
		moves := LegalMoves(id)
		for _, m := range moves {
			if !m.Capture {
				continue
			}
			captures = append(captures, ChessMove{
				SourceID:      id,
				DestinationID: m.ID,
				Capture:       true,
			})
		}
	}
	if len(captures) == 0 {
		return ChessMove{}, false
	}
	return captures[rand.Intn(len(captures))], true
}
