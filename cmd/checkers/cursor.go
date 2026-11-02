package main

import (
	"cards/assets/checkers"
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
		moveable "cursor" - a scaled piece Image.
*/

type Cursor struct {
	cursorIcon binding.String
	obj        *canvas.Image
	layer      *fyne.Container
	show       bool
}

func NewCursor() *Cursor {
	c := &Cursor{
		cursorIcon: binding.NewString(),
		show:       false,
	}
	c.obj = canvas.NewImageFromResource(
		checkers.CheckerResource[EmptyIcon])

	c.obj.Resize(fyne.NewSize(
		float32(Ctx.cursorSize),
		float32(Ctx.cursorSize),
	))
	c.obj.FillMode = canvas.ImageFillContain
	c.cursorIcon.AddListener(binding.NewDataListener(func() {
		icon := common.StringGet(c.cursorIcon)
		c.obj.Resource = checkers.CheckerResource[icon]
		c.obj.Refresh()
	}))
	c.layer = container.NewWithoutLayout(c.obj)
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
	icon = EmptyIcon
	ix := int(mouse.X)
	iy := int(mouse.Y)
	if ix > 0 && ix < 8*Ctx.squareSize &&
		iy > 0 && iy < 8*Ctx.squareSize {
		col = ix / Ctx.squareSize
		row = iy / Ctx.squareSize
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
