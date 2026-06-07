package common

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    toast.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Provide a simple on window temporary message / notification.
*/

// yellow
var BackGroundColor = color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0x8f}
var InfoColor = color.NRGBA{B: 0x7f}
var GoColor = color.NRGBA{G: 0x7f}
var FailColor = color.NRGBA{R: 0x7f}

func Toast(can fyne.Canvas, pos fyne.Position, txt string, c color.NRGBA, d time.Duration) {

	rect := canvas.NewRectangle(BackGroundColor)
	txt = fmt.Sprintf("  %s  ", txt)

	sc := c     // start color
	sc.A = 0x80 //0x11
	ec := c     // end color
	ec.A = 0xff
	t := canvas.NewText(txt, sc)
	t.TextStyle.Bold = true
	a := canvas.NewColorRGBAAnimation(sc, ec, d,
		func(c color.Color) {
			t.Color = c
			canvas.Refresh(t)
		})
	size := t.MinSize()
	a.RepeatCount = 0
	a.AutoReverse = true // true
	a.Curve = fyne.AnimationLinear

	// set rectangle to be same size as text
	rect.Resize(fyne.NewSize(size.Width, size.Height))
	offset := theme.InputBorderSize() + theme.InnerPadding()
	// position the upper left corner
	if pos.X > offset && pos.Y > offset {
		pos.X -= offset
		pos.Y -= offset
	}

	go func() {
		pop := widget.NewPopUp(container.NewWithoutLayout(rect, t), can)
		pop.ShowAtPosition(pos)
		a.Start()
		pop.Show()
		fyne.Do(func() {
			time.Sleep(2 * d)
			pop.Hide()
		})
	}()
}
