package main

import (
	"cards/common"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    ui.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var (
	dealButton      *widget.Button
	undoButton      *widget.Button
	freeStack       *fyne.Container
	foundationStack *fyne.Container
	playCount       binding.Int
	cardStacks      [8]*fyne.Container
)

func initUI() *fyne.Container {
	return container.NewStack(makeGreenPanel(initGameUI()))
}

func makeGreenPanel(content *fyne.Container) *fyne.Container {
	// Soft felt green — color.NRGBA{0, 160, 0, 255}
	// Bright emerald — color.NRGBA{20, 180, 20, 255}
	// Muted mint — color.NRGBA{60, 200, 120, 255}
	// Casino - color.NRGBA{0, 128, 0, 255}
	bg := canvas.NewRectangle(color.NRGBA{20, 180, 20, 50})
	return container.NewStack(bg, content)
	//return content
}

// build the widgets for the actual game.
func initGameUI() *fyne.Container {

	dealButton = widget.NewButton("DEAL", func() {
		deal()
	})
	dealButton.Disable() // TODO
	undoButton = widget.NewButtonWithIcon("UNDO", theme.CancelIcon(), func() {
		undo()
	})
	undoButton.Disable()

	freeStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.cardOffset) * 1.8,
		Vertical: false},
		buildHandImages(4, true)...)
	freeCellGrid := container.NewGridWithColumns(4, freeStack.Objects...)
	foundationStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.cardOffset) * 1.8,
		Vertical: false},
		buildHandImages(4, true)...)
	foundationGrid := container.NewGridWithColumns(4, foundationStack.Objects...)

	playCount = binding.NewInt()
	playText := widget.NewLabel("0")
	playCount.AddListener(binding.NewDataListener(func() {
		playText.SetText(fmt.Sprintf("%4d",
			common.IntGet(playCount)))
	}))
	playGrid := container.NewVBox(
		container.NewCenter(widget.NewLabel("    PLAYS    ")),
		container.NewCenter(playText),
		container.NewHBox(widget.NewLabel("\u2190 Free Cells"), layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(), widget.NewLabel("Foundation \u2192")))
	freeFoundationRow := container.NewHBox(freeCellGrid, playGrid, foundationGrid)

	cardsRow := container.NewGridWithColumns(8)
	//nAdd := 6
	for i := range 8 {
		cardStacks[i] = container.New(&OverlapLayout{
			Offset:   Ctx.cardSizeHeight * .22,
			Vertical: true},
			buildHandImages(1, true)...)
		cardsRow.Add(cardStacks[i])
		//cardStacks[i].Objects = append(cardStacks[i].Objects,
		//	buildHandImages(nAdd, false)...)
		//if i == 3 {
		//	nAdd = 5
		//}
	}

	//board := widget.NewLabel("hello")
	l := cellItemId(CardColumn, 4)
	loc := cellItemLocation(l)
	ix := cellItemIndex(l)
	fmt.Printf("ID %8x, Location %d, index %d\n", l, loc, ix)

	controlBar := container.NewHBox()
	controlBar.Add(dealButton)
	controlBar.Add(undoButton)

	return container.NewBorder(freeFoundationRow, controlBar, nil, nil, cardsRow)
}

func buildHandImages(n int, reverse bool) []fyne.CanvasObject {
	common.DefaultCardImageWidth = Ctx.cardSizeWidth
	common.DefaultCardImageHeight = Ctx.cardSizeHeight
	objs := make([]fyne.CanvasObject, 0)
	for i := 0; i < n; i++ {
		if reverse {
			c := common.NewCardImage(Ctx.cardMap["reverse"].Resource, i)
			objs = append(objs, c)
			c.Refresh()
		} else {
			c := common.NewCardImage(Ctx.cardMap["joker"].Resource, i)
			objs = append(objs, c)
			c.Refresh()
		}
	}
	return objs
}

// ------------------------------------------------------------------
//
// Layout for a horizontal or vertical "hands" - Card images.
//
// ------------------------------------------------------------------

type OverlapLayout struct {
	Offset   float32
	Vertical bool // true = top down, false = left right
}

func (l *OverlapLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n == 0 {
		return
	}

	// Fixed card size — never shrink
	cardW := Ctx.cardSizeWidth
	cardH := Ctx.cardSizeHeight

	// Centering anchor
	var x, y float32
	if l.Vertical {
		x = (size.Width - cardW) / 2
		y = 0
	} else {
		x = 0
		y = (size.Height - cardH) / 2
	}

	// Position each card with fixed overlap
	for i, c := range objects {
		c.Resize(fyne.NewSize(cardW, cardH))

		if l.Vertical {
			c.Move(fyne.NewPos(x, float32(i)*l.Offset))
		} else {
			c.Move(fyne.NewPos(float32(i)*l.Offset, y))
		}
	}
}

func (l *OverlapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// Only enough space for ONE card.
	// The stack grows beyond this using overlap.
	return fyne.NewSize(Ctx.cardSizeWidth, Ctx.cardSizeHeight)
}

/*
func (l *OverlapLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n == 0 {
		return
	}

	// East/West players (vertical hands)
	maxVerticalCardHeight := Ctx.cardSizeHeight
	maxVerticalCardWidth := Ctx.cardSizeWidth

	// North/South players (horizontal hands)
	maxHorizontalCardHeight := Ctx.cardSizeHeight
	maxHorizontalCardWidth := Ctx.cardSizeWidth

	var cardWidth, cardHeight float32

	if l.Vertical {
		// Vertical -- Height is the limiting dimension
		usableHeight := size.Height - float32(n-1)*l.Offset
		if usableHeight < 40 {
			usableHeight = 40
		}

		cardHeight = usableHeight
		cardWidth = cardHeight * (5.0 / 7.0)

		// Clamp to vertical maximums
		if cardHeight > maxVerticalCardHeight {
			cardHeight = maxVerticalCardHeight
			cardWidth = cardHeight * (5.0 / 7.0)
		}
		if cardWidth > maxVerticalCardWidth {
			cardWidth = maxVerticalCardWidth
			cardHeight = cardWidth * (7.0 / 5.0)
		}

	} else {
		// horizontal -- Width is the limiting dimension
		usableWidth := size.Width - float32(n-1)*l.Offset
		if usableWidth < 40 {
			usableWidth = 40
		}

		cardWidth = usableWidth
		cardHeight = cardWidth * (7.0 / 5.0)

		// Clamp to horizontal maximums
		if cardWidth > maxHorizontalCardWidth {
			cardWidth = maxHorizontalCardWidth
			cardHeight = cardWidth * (7.0 / 5.0)
		}
		if cardHeight > maxHorizontalCardHeight {
			cardHeight = maxHorizontalCardHeight
			cardWidth = cardHeight * (5.0 / 7.0)
		}
	}

	// anchor positions
	var x, y float32

	// Centering --
	if l.Vertical {
		// Center horizontally
		x = (size.Width - cardWidth) / 2
		y = 0
	} else {
		// Center vertically
		x = 0
		y = (size.Height - cardHeight) / 2
	}

	// Card size and positions.
	for i, c := range objects {
		c.Resize(fyne.NewSize(cardWidth, cardHeight))

		if l.Vertical {
			c.Move(fyne.NewPos(x, float32(i)*l.Offset))
		} else {
			c.Move(fyne.NewPos(float32(i)*l.Offset, y))
		}
	}

}

func (l *OverlapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	n := len(objects)
	if n == 0 {
		return fyne.NewSize(0, 0)
	}

	// East/West (vertical hands)
	maxVerticalCardHeight := Ctx.cardSizeHeight
	maxVerticalCardWidth := Ctx.cardSizeWidth

	// North/South (horizontal hands)
	maxHorizontalCardHeight := Ctx.cardSizeHeight
	maxHorizontalCardWidth := Ctx.cardSizeWidth

	if l.Vertical {
		// East/West
		h := maxVerticalCardHeight + float32(n-1)*l.Offset
		return fyne.NewSize(maxVerticalCardWidth, h)
	}
	if l.Vertical {
		return fyne.NewSize(maxVerticalCardWidth, maxVerticalCardHeight)
	}

	// North/South
	w := maxHorizontalCardWidth + float32(n-1)*l.Offset
	return fyne.NewSize(w, maxHorizontalCardHeight)
}

*/
