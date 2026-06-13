package main

import (
	"cards/common"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

/*

  File:    dragdrop.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var BaseColor = color.NRGBA{20, 180, 20, 80}
var ColorRed = color.NRGBA{220, 50, 50, 30}
var ColorBlack = color.NRGBA{00, 00, 00, 255}

func createIconCanvas(card common.Card) *canvas.Text {
	icon := card.Icon
	suit := card.Suit
	rect := canvas.NewRectangle(BaseColor)
	var color = color.NRGBA{20, 180, 20, 80}
	switch suit {
	case common.Diamonds, common.Hearts:
		color = ColorRed
	case common.Clubs, common.Spades:
		color = ColorBlack
	}
	t := canvas.NewText(icon, color)
	t.TextStyle.Bold = true
	size := t.MinSize()
	rect.Resize(fyne.NewSize(size.Width, size.Height))

	return canvas.NewText(fmt.Sprintf(" %s ", icon), color)
}
