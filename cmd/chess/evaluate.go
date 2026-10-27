package main

import "fmt"

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

// search entire Board for a threat.
func KingThreatened(side Suit) int {
	kingID := FindKing(side)
	if kingID < 0 {
		return -1
	}
	for id, sq := range Board {
		if sq == nil {
			continue
		}
		if sq.Suit == side { // same suit
			continue
		}
		moves := LegalMoves(id)
		for _, move := range moves {
			if move.ID == kingID {
				return kingID
			}
		}
	}
	return -1
}

// locate the King of this suit
func FindKing(side Suit) int {
	for id, sq := range Board {
		if sq == nil {
			continue
		}
		if sq.Suit == side &&
			sq.Rank == King {
			return id
		}
	}
	return -1
}

// single move result - linear index.
type Move struct {
	ID      int
	Capture bool
}

func (m Move) String() string {
	row, col := ID2RowCol(m.ID)
	if Board[m.ID] == nil {
		return fmt.Sprintf("id %d: row %d, col %d, empty", m.ID, row, col)
	}
	return fmt.Sprintf("id %d:row %d, col %d, %s - Capture %t",
		m.ID, row, col, Board[m.ID].Icon, m.Capture)
}

// list all the possible moves for this piece
func LegalMoves(id int) []Move {
	piece := Board[id]
	if piece == nil {
		return nil
	}
	switch piece.Rank {
	case Pawn:
		return PawnMoves(Board, piece)
	case Knight:
		return KnightMoves(Board, piece)
	case Bishop:
		return LgalPieceMoves(Board, piece, BishopDirs, 7)
	case Rook:
		return LgalPieceMoves(Board, piece, RookDirs, 7)
	case Queen:
		return LgalPieceMoves(Board, piece, QueenDirs, 7)
	case King:
		return LgalPieceMoves(Board, piece, QueenDirs, 1)
	}
	return nil
}

func LgalPieceMoves(board [64]*ChessSquare, piece *ChessSquare, dirs []Delta, maxSteps int) []Move {
	var moves []Move
	for _, d := range dirs {
		row := piece.Row
		col := piece.Col
		for step := 1; step <= maxSteps; step++ {
			row += d.Row
			col += d.Col
			if row < 0 || row > 7 ||
				col < 0 || col > 7 {
				break
			}
			id := RowCol2ID(row, col)
			target := board[id]
			if target == nil {
				moves = append(moves, Move{
					ID:      id,
					Capture: false,
				})
				continue
			}
			if target.Suit != piece.Suit {
				moves = append(moves, Move{
					ID:      id,
					Capture: true,
				})
			}
			break
		}
	}
	return moves
}

func PawnMoves(board [64]*ChessSquare, piece *ChessSquare) []Move {
	var moves []Move
	dir := 1
	startRow := 1
	if piece.Suit == White {
		dir = -1
		startRow = 6
	}
	//
	// one square forward
	//
	row := piece.Row + dir
	col := piece.Col
	if row >= 0 && row < 8 {
		id := RowCol2ID(row, col)
		if board[id] == nil {
			moves = append(moves, Move{
				ID:      id,
				Capture: false,
			})
			//
			// opening two-square move
			//
			if piece.Row == startRow {
				row2 := piece.Row + 2*dir
				id2 := RowCol2ID(row2, col)
				if board[id2] == nil {
					moves = append(moves, Move{
						ID:      id2,
						Capture: false,
					})
				}
			}
		}
	}
	//
	// capture left
	//
	row = piece.Row + dir
	col = piece.Col - 1
	if row >= 0 && row < 8 &&
		col >= 0 && col < 8 {
		id := RowCol2ID(row, col)
		target := board[id]
		if target != nil &&
			target.Suit != piece.Suit {
			moves = append(moves, Move{
				ID:      id,
				Capture: true,
			})
		}
	}
	//
	// capture right
	//
	row = piece.Row + dir
	col = piece.Col + 1
	if row >= 0 && row < 8 &&
		col >= 0 && col < 8 {
		id := RowCol2ID(row, col)
		target := board[id]
		if target != nil &&
			target.Suit != piece.Suit {
			moves = append(moves, Move{
				ID:      id,
				Capture: true,
			})
		}
	}
	return moves
}

func KnightMoves(board [64]*ChessSquare, piece *ChessSquare) []Move {
	var moves []Move
	for _, d := range KnightJumps {
		row := piece.Row + d.Row
		col := piece.Col + d.Col
		if row < 0 || row > 7 ||
			col < 0 || col > 7 {
			continue
		}
		id := RowCol2ID(row, col)
		target := board[id]
		if target == nil {
			moves = append(moves, Move{
				ID:      id,
				Capture: false,
			})
			continue
		}
		if target.Suit != piece.Suit {
			moves = append(moves, Move{
				ID:      id,
				Capture: true,
			})
		}
	}
	return moves
}
