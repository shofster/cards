package main

import (
	"cards/assets/checkers"
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

const ( // piece image names
	EmptyIcon      = "empty"
	BlackIcon      = "black"
	RedIcon        = "red"
	BlackCrownIcon = "blackCrown"
	RedCrownIcon   = "redCrown"
)

type Suit int
type Rank int

const (
	Black Suit = iota
	Red
)

const (
	Pawn Rank = iota
	King
)

func otherPlayer(suit Suit) Suit {
	if suit == Red {
		return Black
	}
	return Red
}

func SuitCode(s Suit) string {
	switch s {
	case Black:
		return checkers.BlackName
	case Red:
		return checkers.RedName
	}
	return fmt.Sprintf("%d", s)
}

func ID2RowCol(id int) (row, col int) {
	row = id / 8
	col = id - row*8
	return
}

func RowCol2ID(row, col int) (id int) {
	return row*8 + col
}

type CheckerSquare struct {
	Suit     Suit
	Rank     Rank   // Pawn, King
	Icon     string // EmptyIcon, BlackIcon, BlackCrownIcon, RedIcon, RedCrownIcon
	Row, Col int
}

func PieceIcon(suit Suit, rank Rank) string {
	switch suit {
	case Black:
		if rank == King {
			return BlackCrownIcon
		}
		return BlackIcon
	case Red:
		if rank == King {
			return RedCrownIcon
		}
		return RedIcon
	}
	return EmptyIcon
}

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
