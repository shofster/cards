package main

import (
	"cards/assets/checkers"
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

// initial Board images
func newCheckerBoard() *fyne.Container {

	// background
	board := canvas.NewImageFromResource(resourceBoardPng)
	board.FillMode = canvas.ImageFillContain
	board.Resize(fyne.NewSize(800, 800))
	board.Move(fyne.NewPos(0, 0))

	squareSize := Ctx.squareSize
	imageSize := float32(squareSize - 2)
	common.DefaultCardImageWidth = imageSize
	common.DefaultCardImageHeight = imageSize

	resetBoardPieces()
	CardImages = container.NewWithoutLayout()
	resetBoardImages(squareSize, imageSize)

	// overlay the tracker grid on top of board
	tracker.Resize(fyne.NewSize(float32(squareSize*8), float32(squareSize*8)))
	tracker.Move(fyne.NewPos(0, 0))
	cursor.layer.Resize(fyne.NewSize(float32(squareSize*8), float32(squareSize*8)))
	cursor.layer.Move(fyne.NewPos(0, 0))
	return container.NewWithoutLayout(board, CardImages, tracker, cursor.layer)

}

func resetBoardImages(squareSize int, imageSize float32) {
	for row := range 8 {
		for col := range 8 {
			id := RowCol2ID(row, col)
			res := checkers.CheckerResource[EmptyIcon]
			if Board[id] != nil {
				res = checkers.CheckerResource[Board[id].Icon]
			}
			ci := common.NewCardImage(res, id)
			// position inside the board’s 1px border
			x := float32(col*squareSize + 1)
			y := float32(row*squareSize + 1)
			ci.Move(fyne.NewPos(x, y))
			ci.Resize(fyne.NewSize(imageSize, imageSize))
			CardImages.Add(ci)
		}
	}
}

func resetBoardPieces() {
	// clear board
	for row := range 8 {
		for col := range 8 {
			Board[RowCol2ID(row, col)] = nil
		}
	}
	// ("black") squares only
	for row := range 8 {
		for col := range 8 {
			if !IsPlayable(row, col) {
				continue
			}
			id := RowCol2ID(row, col)
			switch {
			// Black pieces
			case row < 3:
				Board[id] = &CheckerSquare{
					Suit: Black,
					Rank: Pawn,
					Icon: PieceIcon(Black, Pawn),
					Row:  row,
					Col:  col,
				}
			// Red pieces
			case row > 4:
				Board[id] = &CheckerSquare{
					Suit: Red,
					Rank: Pawn,
					Icon: PieceIcon(Red, Pawn),
					Row:  row,
					Col:  col,
				}
			}
		}
	}
}
