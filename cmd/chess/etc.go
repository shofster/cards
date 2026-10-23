package main

import (
	"fmt"
)

/*

  File:    etc.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

type Suit int
type Rank int

const (
	Black Suit = iota
	White
)

func otherPlayer(suit Suit) Suit {
	if suit == White {
		return Black
	}
	return White
}

func SuitCode(s Suit) string {
	switch s {
	case Black:
		return "GOLD"
	case White:
		return "SILVER"
	}
	return fmt.Sprintf("%d", s)
}

const (
	Pawn Rank = iota
	Rook
	Knight
	Bishop
	Queen
	King
)

type ChessSquare struct {
	Suit     Suit
	Rank     Rank
	Icon     string
	Row, Col int
}

func ID2RowCol(id int) (row, col int) {
	row = id / 8
	col = id - row*8
	return
}

func RowCol2ID(row, col int) (id int) {
	return row*8 + col
}

func RankCode(r Rank) string {
	switch r {
	case Pawn:
		return "pawn"
	case Rook:
		return "rook"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Queen:
		return "queen"
	case King:
		return "king"
	}
	return fmt.Sprintf("%d", r)
}

var blackHonorsRow = [8]string{"gR", "gHl", "gB", "gQ", "gK", "gB", "gHr", "gR"}
var blackPawn = "gP"
var whiteHonorsRow = [8]string{"sR", "sHl", "sB", "sQ", "sK", "sB", "sHr", "sR"}
var whitePawn = "sP"
var honorsRank = [8]Rank{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

type Delta struct {
	Row int
	Col int
}

var North = Delta{-1, 0}
var South = Delta{1, 0}
var East = Delta{0, 1}
var West = Delta{0, -1}
var NorthEast = Delta{-1, 1}
var NorthWest = Delta{-1, -1}
var SouthEast = Delta{1, 1}
var SouthWest = Delta{1, -1}

var RookDirs = []Delta{
	North,
	South,
	East,
	West,
}
var BishopDirs = []Delta{
	NorthEast,
	NorthWest,
	SouthEast,
	SouthWest,
}
var QueenDirs = []Delta{
	North,
	South,
	East,
	West,
	NorthEast,
	NorthWest,
	SouthEast,
	SouthWest,
}
var KnightJumps = []Delta{
	{-2, -1},
	{-2, 1},
	{-1, -2},
	{-1, 2},
	{1, -2},
	{1, 2},
	{2, -1},
	{2, 1},
}
