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
	dealButton     *widget.Button
	undoButton     *widget.Button
	cardStacks     [8]*fyne.Container
	freeGrid       *fyne.Container
	playGrid       *fyne.Container
	foundationGrid *fyne.Container
	cardsGrid      *fyne.Container
	playCount      binding.Int
	hintGroup      *widget.RadioGroup
	hintOnOff      binding.Bool
	hintReason     binding.String
	tracker        *common.Tracker
	mouseX, mouseY binding.Int
	cursor         *Cursor
	mouseIcon      binding.String
	clickDir       binding.String
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

	common.DefaultCardImageWidth = Ctx.cardSizeWidth
	common.DefaultCardImageHeight = Ctx.cardSizeHeight

	dealButton = widget.NewButton("DEAL", func() {
		deal()
		hintGroup.SetSelected("off")
		common.IntSet(playCount, 0)
	})
	dealButton.Enable()
	undoButton = widget.NewButtonWithIcon("UNDO", theme.CancelIcon(), func() {
		popUndo()
	})
	undoButton.Disable()

	var col int
	for col = 0; col < 4; col++ {
		gameState.FreeCells[col] = NewCellItem(Free, col, 0, emptyCard())
		gameState.FoundationCells[col] = NewCellItem(Foundation, col, 0, aceCard(col))
	}

	freeGrid = container.NewGridWithColumns(4)
	resetFreeCellImages(&gameState, freeGrid)
	foundationGrid = container.NewGridWithColumns(4)
	resetFoundationCellImages(&gameState, foundationGrid)

	Ctx.cardSizeWidth *= 1.2 // smaller for Card Columns (and picker)
	Ctx.cardSizeHeight *= 1.2
	cardsGrid = container.NewGridWithColumns(8)
	for col := 0; col < 8; col++ {
		gameState.CardCells[col] = make([]*CellItem, 0)
		cardStacks[col] = container.New(&OverlapLayout{
			Offset:   float32(Ctx.cardOffset) * 1,
			Vertical: true})
		resetCardColumnImages(col, &gameState)
		cardsGrid.Objects = append(cardsGrid.Objects, cardStacks[col])
	}

	playCount = binding.NewInt()
	playText := widget.NewLabel("0")
	playCount.AddListener(binding.NewDataListener(func() {
		playText.SetText(fmt.Sprintf("%4d",
			common.IntGet(playCount)))
	}))
	playGrid = container.NewVBox(
		container.NewCenter(widget.NewLabel("    PLAYS    ")),
		container.NewCenter(playText),
		container.NewHBox(widget.NewLabel("\u2190  Free Cells"), layout.NewSpacer()),
		container.NewHBox(layout.NewSpacer(), widget.NewLabel("Foundation \u2192")))
	freeFoundationRow := container.NewHBox(freeGrid, playGrid, foundationGrid)

	xmouse := widget.NewLabel("")
	ymouse := widget.NewLabel("")
	mouseX = binding.NewInt()
	mouseY = binding.NewInt()
	mouseX.AddListener(binding.NewDataListener(func() {
		xmouse.SetText(fmt.Sprintf("X:%-4d",
			common.IntGet(mouseX)))
	}))
	mouseY.AddListener(binding.NewDataListener(func() {
		ymouse.SetText(fmt.Sprintf("Y:%-4d",
			common.IntGet(mouseY)))
	}))
	mouseIcon = binding.NewString()
	clickDir = binding.NewString()
	mouseInfo := container.NewHBox(xmouse, ymouse,
		widget.NewLabelWithData(mouseIcon), widget.NewLabelWithData(clickDir))
	common.StringSet(mouseIcon, "xxx")
	common.StringSet(clickDir, "\u2191") // up

	hintOnOff = binding.NewBool()
	hintGroup = widget.NewRadioGroup([]string{"on", "off"}, func(s string) {
		switch s { // set the moveHint trigger
		case "on":
			common.BoolSet(hintOnOff, true)
			common.StringSet(hintReason, HintMap[1])
			moveHint(&gameState)
		case "off":
			common.BoolSet(hintOnOff, false)
			common.StringSet(hintReason, HintMap[0])
		}
	})
	hintGroup.Horizontal = true
	hintGroup.Selected = "off"
	hintReason = binding.NewString()
	common.StringSet(hintReason, HintMap[0])
	hintText := widget.NewLabelWithData(hintReason)
	hintInfo := container.NewHBox(widget.NewLabel("Hints?"), hintGroup, hintText)
	//hintOnOff.AddListener(binding.NewDataListener(func() {
	//	if common.BoolGet(hintOnOff) {
	//		hintx.Selected = "on"
	//		common.StringSet(hintReason, HintMap[1])
	//	} else {
	//		hintx.Selected = "off"
	//		common.StringSet(hintReason, HintMap[0])
	//	}
	//}))

	controlBar := container.NewHBox()
	controlBar.Add(dealButton)
	controlBar.Add(undoButton)
	controlBar.Add(layout.NewSpacer())
	controlBar.Add(hintInfo)
	controlBar.Add(mouseInfo)

	cursor = NewCursor()
	tracker = common.NewTracker()
	board := container.NewStack(
		container.NewBorder(freeFoundationRow,
			nil, nil, nil, cardsGrid),
		tracker,
		cursor.layer,
	)
	overlay := container.NewBorder(nil, controlBar, nil, nil, board)
	tracker.Event = trackMouse
	tracker.Click = nil

	return overlay
}

func resetFreeCellImages(gs *GameState, grid *fyne.Container) {
	grid.Objects = nil
	for i := 0; i < 4; i++ {
		ci := common.NewCardImage(gs.FreeCells[i].card.Resource, gs.FreeCells[i].card.Index)
		grid.Add(ci)
	}
	grid.Refresh()
}

func resetFoundationCellImages(gs *GameState, grid *fyne.Container) {
	grid.Objects = nil
	for i := 0; i < 4; i++ {
		ci := common.NewCardImage(gs.FoundationCells[i].card.Resource, gs.FreeCells[i].card.Index)
		grid.Add(ci)
	}
	grid.Refresh()
}

func resetCardColumnImages(col int, gs *GameState) {
	stack := cardStacks[col]
	cells := gs.CardCells[col]
	stack.Objects = stack.Objects[:0]
	if len(cells) == 0 {
		emptyID := NewCellItemId(Empty, col, 0)
		ci := common.NewCardImage(emptyCard().Resource, int(emptyID))
		stack.Objects = append(stack.Objects, ci)
		stack.Refresh()
		return
	}
	for _, cell := range cells {
		ci := common.NewCardImage(cell.card.Resource, int(cell.id))
		stack.Objects = append(stack.Objects, ci)
	}
	stack.Refresh()
}

// replace the CardImage of a free cell
func updateFreeCell(col int, gs *GameState) {
	cell := gs.FreeCells[col]
	ci := common.NewCardImage(cell.card.Resource, int(cell.id))
	freeGrid.Objects[col] = ci
	freeGrid.Refresh()
}

// replace the CardImage of a foundation cell
func updateFoundationCell(col int, gs *GameState) {
	cell := gs.FoundationCells[col]
	ci := common.NewCardImage(cell.card.Resource, int(cell.id))
	foundationGrid.Objects[col] = ci
	foundationGrid.Refresh()
}

// ------------------------------------------------------------------
//
// CardImage Layout for horizontal or vertical "hands".
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

	// Ffixed card size
	cardWidth := Ctx.cardSizeWidth
	cardHeight := Ctx.cardSizeHeight

	// Centering anchor
	var x, y float32
	if l.Vertical {
		x = (size.Width - cardWidth) / 2
		y = 0
	} else {
		x = 0
		y = (size.Height - cardHeight) / 2
	}

	// position each card with fixed overlap
	for i, c := range objects {
		c.Resize(fyne.NewSize(cardWidth, cardHeight))
		if l.Vertical {
			c.Move(fyne.NewPos(x, float32(i)*l.Offset))
		} else {
			c.Move(fyne.NewPos(float32(i)*l.Offset, y))
		}
	}
}

func (l *OverlapLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	// space for 1 card
	// the stack grows beyond via overlap
	return fyne.NewSize(Ctx.cardSizeWidth, Ctx.cardSizeHeight)
}
