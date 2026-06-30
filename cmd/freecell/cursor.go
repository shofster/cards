package main

import (
	"cards/common"
	"image/color"
	"strings"

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
moveable "cursor" - a Card.Icon as Red or Black based on Card.Suit
*/
type Cursor struct {
	cursorCard binding.String
	obj        canvas.Text
	layer      *fyne.Container
	show       bool
}

var ColorRed = color.NRGBA{255, 0, 0, 255}
var ColorBlack = color.NRGBA{00, 00, 00, 255}

func NewCursor() *Cursor {
	c := &Cursor{
		cursorCard: binding.NewString(),
		show:       false,
	}
	c.obj = canvas.Text{
		Alignment: 0,
		Color:     ColorBlack,
		Text:      "A♠",
		TextSize:  float32(Ctx.dragSize),
		TextStyle: fyne.TextStyle{
			Bold:      true,
			Monospace: false,
			Symbol:    false,
		},
		FontSource: nil,
	}
	c.cursorCard.AddListener(binding.NewDataListener(func() {
		c.obj.Text = common.StringGet(c.cursorCard)
		if len(c.obj.Text) > 0 {
			_, s := common.RankSuit(strings.Trim(c.obj.Text, " "))
			switch s {
			case common.Diamonds, common.Hearts:
				c.obj.Color = ColorRed
			default:
				c.obj.Color = ColorBlack
			}
			c.obj.Refresh()
		}
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
