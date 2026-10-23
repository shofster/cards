package main

import (
	"cards/assets/chess"
	"cards/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

/*

  File:    board.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: the 8x8 square board. with pieces placed.
*/

var SquareSize int = 100
var ImageSize float32 = 98
var EmptyCell = "empty"

func NewChessBoard() *fyne.Container {

	// background
	board := canvas.NewImageFromResource(resourceBoardPng)
	board.FillMode = canvas.ImageFillContain
	board.Resize(fyne.NewSize(800, 800))
	board.Move(fyne.NewPos(0, 0))

	ImageSize = float32(SquareSize - 2)
	common.DefaultCardImageWidth = ImageSize
	common.DefaultCardImageHeight = ImageSize

	CardImages = container.NewWithoutLayout()
	for row := range 8 {
		for col := range 8 {
			id := RowCol2ID(row, col)
			// transparent empty square
			ci := common.NewCardImage(chess.ChessResource[EmptyCell], id)
			// position inside the board’s 1px border
			x := float32(col*SquareSize + 1)
			y := float32(row*SquareSize + 1)
			ci.Move(fyne.NewPos(x, y))
			ci.Resize(fyne.NewSize(ImageSize, ImageSize))
			CardImages.Add(ci)
		}
	}
	// set up piece codes
	var layout [8][8]string
	layout[0] = blackHonorsRow
	for col := range 8 {
		layout[1][col] = "gP"
	}
	for col := range 8 {
		layout[6][col] = "sP"
	}
	layout[7] = whiteHonorsRow
	for row := range 8 {
		for col := range 8 {
			key := layout[row][col]
			if key != EmptyCell { // the image widget
				ix := row*8 + col
				ci := CardImages.Objects[ix].(*common.CardImage)
				ci.ReplaceImage(chess.ChessResource[key])
			}
		}
	}

	// overlay grid on top of board
	tracker.Resize(fyne.NewSize(float32(SquareSize*8), float32(SquareSize*8)))
	tracker.Move(fyne.NewPos(0, 0))
	cursor.layer.Resize(fyne.NewSize(float32(SquareSize*8), float32(SquareSize*8)))
	cursor.layer.Move(fyne.NewPos(0, 0))
	return container.NewWithoutLayout(board, CardImages, tracker, cursor.layer)

}

func buildBoard() {
	for row := range 8 { // reset if new game
		for col := range 8 {
			Board[RowCol2ID(row, col)] = nil
		}
	}
	// mtaches NewChessBoard's images
	for col := range 8 {
		Board[RowCol2ID(0, col)] = &ChessSquare{
			Suit: Black,
			Rank: honorsRank[col],
			Icon: blackHonorsRow[col],
			Row:  0,
			Col:  col,
		}
		Board[RowCol2ID(1, col)] = &ChessSquare{
			Suit: Black,
			Rank: Pawn,
			Icon: "gP",
			Row:  1,
			Col:  col,
		}
		Board[RowCol2ID(6, col)] = &ChessSquare{
			Suit: White,
			Rank: Pawn,
			Icon: "sP",
			Row:  6,
			Col:  col,
		}
		Board[RowCol2ID(7, col)] = &ChessSquare{
			Suit: White,
			Rank: honorsRank[col],
			Icon: whiteHonorsRow[col],
			Row:  7,
			Col:  col,
		}
	}
}
