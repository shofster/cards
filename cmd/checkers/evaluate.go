package main

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

func IsPlayable(row, col int) bool {
	return (row+col)%2 == 1
}

func IsPromotionRow(color Suit, row int) bool {
	return (color == Black && row == 7) ||
		(color == Red && row == 0)
}

func IsCrowned(piece *CheckerSquare) bool {
	if piece.Rank != King { // not yet crowned
		switch piece.Suit {
		case Black:
			if piece.Row == 7 {
				piece.Rank = King
				piece.Icon = PieceIcon(Black, King)
				return true
			}
		case Red:
			if piece.Row == 0 {
				piece.Rank = King
				piece.Icon = PieceIcon(Red, King)
				return true
			}
		}
	}
	return false
}

var BlackDirs = []Delta{
	SouthEast,
	SouthWest,
}

var RedDirs = []Delta{
	NorthEast,
	NorthWest,
}

var KingDirs = []Delta{
	NorthEast,
	NorthWest,
	SouthEast,
	SouthWest,
}

// single move result - linear index.
type CheckerMove struct {
	SourceID      int
	DestinationID int
	Capture       bool
	CapturedID    int
	Captures      []CaptureMove
}

// subsequent Capture move
type CaptureMove struct {
	DestinationID int
	CapturedID    int
}

func LegalMoves(id int) []CheckerMove {
	piece := Board[id]
	if piece == nil {
		return nil
	}
	switch piece.Rank {
	case Pawn:
		return PawnMoves(Board, piece)
	case King:
		return KingMoves(Board, piece)
	}
	return nil
}

func PawnMoves(board [64]*CheckerSquare, piece *CheckerSquare) []CheckerMove {
	var moves []CheckerMove
	dirs := RedDirs
	if piece.Suit == Black {
		dirs = BlackDirs
	}
	for _, dir := range dirs {
		moves = append(moves,
			CheckDirection(board, piece, dir)...)
	}
	return moves
}

func KingMoves(board [64]*CheckerSquare, piece *CheckerSquare) []CheckerMove {
	var moves []CheckerMove
	for _, dir := range KingDirs {
		moves = append(moves,
			CheckDirection(board, piece, dir)...)
	}
	return moves
}

func CheckDirection(board [64]*CheckerSquare, piece *CheckerSquare, d Delta) []CheckerMove {
	var moves []CheckerMove
	sourceID := RowCol2ID(
		piece.Row,
		piece.Col,
	)

	// normal move
	row := piece.Row + d.Row
	col := piece.Col + d.Col
	if row >= 0 && row < 8 && col >= 0 && col < 8 {
		destinationID := RowCol2ID(row, col)
		if board[destinationID] == nil {
			moves = append(moves,
				CheckerMove{
					SourceID:      sourceID,
					DestinationID: destinationID,
					Capture:       false,
					CapturedID:    -1,
				},
			)
		}
	}

	// capture move
	enemyRow := piece.Row + d.Row
	enemyCol := piece.Col + d.Col
	jumpRow := piece.Row + 2*d.Row
	jumpCol := piece.Col + 2*d.Col
	if jumpRow < 0 || jumpRow > 7 || jumpCol < 0 || jumpCol > 7 {
		return moves // outside of board
	}

	enemyID := RowCol2ID(enemyRow, enemyCol)
	jumpID := RowCol2ID(jumpRow, jumpCol)
	enemy := board[enemyID]
	if enemy == nil || enemy.Suit == piece.Suit || board[jumpID] != nil {
		return moves // no candidates
	}
	// initial capture
	move := CheckerMove{
		SourceID:      sourceID,
		DestinationID: jumpID,
		Capture:       true,
		CapturedID:    enemyID,
	}
	// jumped to kings row
	if piece.Rank == Pawn && IsPromotionRow(piece.Suit, jumpRow) {
		moves = append(moves, move)
		return moves
	}

	// check for subsequent captures after this onr
	// temp setup the current state
	nextBoard := board
	nextPiece := *piece
	nextBoard[sourceID] = nil
	nextBoard[enemyID] = nil
	nextPiece.Row = jumpRow
	nextPiece.Col = jumpCol
	nextBoard[jumpID] = &nextPiece
	move.Captures = FindSubsequentCaptures(nextBoard, &nextPiece)
	moves = append(moves, move)

	return moves
}

func FindSubsequentCaptures(board [64]*CheckerSquare, piece *CheckerSquare) []CaptureMove {
	var dirs []Delta
	if piece.Rank == King {
		dirs = KingDirs
	} else {
		dirs = RedDirs
		if piece.Suit == Black {
			dirs = BlackDirs
		}
	}

	for _, d := range dirs {
		enemyRow := piece.Row + d.Row
		enemyCol := piece.Col + d.Col
		jumpRow := piece.Row + 2*d.Row
		jumpCol := piece.Col + 2*d.Col
		if jumpRow < 0 || jumpRow > 7 || jumpCol < 0 || jumpCol > 7 {
			continue // outside of board
		}
		enemyID := RowCol2ID(enemyRow, enemyCol)
		jumpID := RowCol2ID(jumpRow, jumpCol)
		enemy := board[enemyID]
		if enemy == nil || enemy.Suit == piece.Suit || board[jumpID] != nil {
			continue
		}
		capture := CaptureMove{
			DestinationID: jumpID,
			CapturedID:    enemyID,
		}
		// ends at kings row
		//
		if piece.Rank == Pawn && IsPromotionRow(piece.Suit, jumpRow) {
			return []CaptureMove{
				capture,
			}
		}

		// iterate for subsequent captures after this one
		// setup the current state
		nextBoard := board
		nextPiece := *piece
		currentID := RowCol2ID(piece.Row, piece.Col)
		nextBoard[currentID] = nil
		nextBoard[enemyID] = nil
		nextPiece.Row = jumpRow
		nextPiece.Col = jumpCol
		nextBoard[jumpID] = &nextPiece
		rest := FindSubsequentCaptures(nextBoard, &nextPiece)
		return append([]CaptureMove{capture}, rest...)
	}
	return nil
}
