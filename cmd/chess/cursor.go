package main

import (
	"cards/assets/chess"
	"cards/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

/*

  File:    cursor.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/
/*
  Description:
		moveable "cursor" - a chess piece as the UNICODE glyph.
*/

type Cursor struct {
	cursorIcon binding.String
	obj        canvas.Text
	layer      *fyne.Container
	show       bool
}

func NewCursor() *Cursor {
	c := &Cursor{
		cursorIcon: binding.NewString(),
		show:       false,
	}
	c.obj = canvas.Text{
		Alignment: 0,
		Text:      "♚", // the UNICODE character to track position
		TextSize:  float32(Ctx.dragSize),
		TextStyle: fyne.TextStyle{
			Bold:      true,
			Monospace: false,
			Symbol:    false,
		},
		FontSource: nil,
	}
	c.cursorIcon.AddListener(binding.NewDataListener(func() {
		c.obj.Text = chess.ChessUnicode[common.StringGet(c.cursorIcon)]
		c.obj.Refresh()
	}))
	c.layer = container.NewWithoutLayout(&c.obj)
	return c
}

func (c *Cursor) Show(show bool) {
	if show {
		c.obj.Show()
		c.obj.Refresh()
	} else {
		c.obj.Hide()
	}
}

func (c *Cursor) Move(position fyne.Position) {
	c.obj.Move(position)
}

func FollowMouse(mouse fyne.Position) (row, col int, icon string) {
	icon = EmptyCell
	ix := int(mouse.X)
	iy := int(mouse.Y)
	if ix > 0 && ix < 8*SquareSize &&
		iy > 0 && iy < 8*SquareSize {
		col = ix / SquareSize
		row = iy / SquareSize
		if Board[RowCol2ID(row, col)] != nil {
			icon = Board[RowCol2ID(row, col)].Icon
			cursor.Move(fyne.NewPos(
				mouse.X-10,
				mouse.Y-10))
		}
	}
	common.StringSet(cursor.cursorIcon, icon)
	return
}
