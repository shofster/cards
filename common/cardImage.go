package common

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    cardImage.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Widget wrapper for a dimmable, highlightable, tappable Card image.
*/

type CardImage struct {
	widget.BaseWidget

	Img    *canvas.Image
	Dim    *canvas.Rectangle
	HiLite *canvas.Rectangle // outer glow
	Mark   *canvas.Circle

	MinWidth  float32
	MinHeight float32

	ID          int
	Enabled     bool
	Highlighted bool
	OnTap       func(*CardImage)
}

var DefaultCardImageWidth float32 = 100
var DefaultCardImageHeight float32 = 140
var DefaultCardImageMarkSize float32 = 18

func NewCardImage(res fyne.Resource, id int) *CardImage {
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain

	// Outer glow
	hilite := canvas.NewRectangle(color.NRGBA{255, 215, 0, 220})
	hilite.TopLeftCornerRadius = 8
	hilite.TopRightCornerRadius = 8
	hilite.BottomLeftCornerRadius = 4
	hilite.BottomRightCornerRadius = 4
	hilite.Hide()

	// Dim overlay
	dim := canvas.NewRectangle(color.NRGBA{0, 0, 0, 60})
	dim.Hide()

	// the hint mark
	mark := canvas.NewCircle(color.NRGBA{0, 200, 0, 180})
	mark.Hide() // default hidden

	c := &CardImage{
		Img:       img,
		HiLite:    hilite,
		Dim:       dim,
		Mark:      mark,
		MinWidth:  DefaultCardImageWidth,
		MinHeight: DefaultCardImageHeight,
		ID:        id,
		Enabled:   true,
	}

	c.ExtendBaseWidget(c)
	return c
}

func (c *CardImage) CreateRenderer() fyne.WidgetRenderer {
	return &cardImageRenderer{c: c}
}

type cardImageRenderer struct {
	c *CardImage
}

func (r *cardImageRenderer) Layout(size fyne.Size) {

	// HiLite
	r.c.HiLite.Resize(size.Add(fyne.NewSize(8, 8)))
	r.c.HiLite.Move(fyne.NewPos(-4, -4))

	// Card image
	r.c.Img.Resize(size)

	// Dim overlay
	r.c.Dim.Resize(size)

	r.c.Mark.Resize(fyne.NewSize(DefaultCardImageMarkSize, DefaultCardImageMarkSize))
	r.c.Mark.Move(fyne.NewPos(3, size.Height-DefaultCardImageMarkSize-3))
}

func (r *cardImageRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.c.MinWidth, r.c.MinHeight)
}

func (r *cardImageRenderer) Refresh() {
	fyne.Do(func() {
		r.c.Img.Refresh()
		r.c.Dim.Refresh()
		r.c.HiLite.Refresh()
	})
}

func (r *cardImageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.c.HiLite,
		r.c.Img,
		r.c.Dim,
		r.c.Mark,
	}
}

func (r *cardImageRenderer) Destroy() {}

func (c *CardImage) Reset(res fyne.Resource) {
	c.Enabled = false
	c.OnTap = nil
	c.Highlighted = false
	c.HiLite.Hide()
	c.Dim.Hide()
	c.Replace(res)
}

func (c *CardImage) Replace(res fyne.Resource) {
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	c.Img = img
	c.Refresh()
}

func (c *CardImage) SetDim(on bool) {
	if on {
		c.Dim.Show()
	} else {
		c.Dim.Hide()
	}
	c.Refresh()
}

func (c *CardImage) SetHighlite(on bool) {
	c.Highlighted = on
	if on {
		c.HiLite.Show()
	} else {
		c.HiLite.Hide()
	}
	c.Refresh()
}

func (c *CardImage) SetDimmed(on bool) {
	if on {
		c.Dim.Show()
		c.Enabled = false
	} else {
		c.Dim.Hide()
		c.Enabled = true
	}
	c.Refresh()
}

func (c *CardImage) SetMark(on bool) {
	if on {
		c.Mark.Show()
	} else {
		c.Mark.Hide()
	}
	c.Refresh()
}

func (c *CardImage) Tapped(_ *fyne.PointEvent) {
	if c.Enabled && c.OnTap != nil {
		c.OnTap(c)
	}
}
