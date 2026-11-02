package main

import (
	"cards/assets/checkers"
	"cards/common"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

var player *PlayerPanel
var bot *PlayerPanel
var checkerBoard *fyne.Container
var autoPlay bool
var playerTurn *widget.Label

func initUI() *fyne.Container {
	playerTurn = widget.NewLabel("PLAY")
	// necessary game variables
	initGame()
	return container.NewStack(initGameUI())
}

func initGameUI() *fyne.Container {
	playerControls := container.NewHBox()
	playerControls.Add(playerTurn)
	player = NewPlayerPanel(Ctx.playerName, Red, playerControls)
	botControls := container.NewHBox()
	botControls.Add(widget.NewLabel("Auto"))
	auto := common.NewToggle(func(b bool) {
		autoPlay = b
	})
	botControls.Add(auto)
	bot = NewPlayerPanel(Ctx.botName, Black, botControls)
	checkerBoard = newCheckerBoard()
	return container.NewBorder(nil, nil, player.content, bot.content, checkerBoard)
}

type PlayerPanel struct {
	suit      Suit
	nCaptured int
	captured  []*canvas.Image
	content   *fyne.Container
}

func NewPlayerPanel(title string, suit Suit, controls *fyne.Container) *PlayerPanel {
	pp := &PlayerPanel{
		suit:     suit,
		captured: make([]*canvas.Image, 12),
	}
	// up to 12 captured pieces
	grid := container.NewGridWithColumns(4)
	for i := range 12 {
		res := BlackCrownIcon
		if suit == Black {
			res = RedCrownIcon
		}
		img := canvas.NewImageFromResource(checkers.CheckerResource[res])
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(24, 24))
		pp.captured[i] = img
		grid.Add(img)
	}
	// left or right of the Border
	if len(title) > 8 {
		title = title[:8]
	}
	pp.content = container.NewVBox(
		container.NewCenter(widget.NewLabel(title)),
		container.NewCenter(widget.NewLabel(SuitCode(suit))),
		layout.NewSpacer(),
		controls,
		layout.NewSpacer(),
		container.NewCenter(widget.NewLabel("Captured:")),
		container.NewCenter(grid))

	return pp
}

func resetPlayerPanels() {
	player.clear()
	bot.clear()
}
func (pp *PlayerPanel) clear() {
	for i := range 12 {
		pp.setCaptured(i, EmptyIcon)
	}
	pp.nCaptured = 0
}
func (pp *PlayerPanel) addCapture(icon string) {
	pp.setCaptured(pp.nCaptured, icon)
	pp.nCaptured++
}
func (pp *PlayerPanel) setCaptured(slot int, icon string) {
	//fmt.Printf("setCaptured:suit %s, slot %d, icon <%s>",
	//	SuitCode(pp.suit), slot, icon)
	//defer fmt.Println()
	if r, ok := checkers.CheckerResource[icon]; ok {
		pp.captured[slot].Resource = r
		pp.captured[slot].Refresh()
		pp.content.Refresh()
		//fmt.Print(" - OK")
	}
}

func updateSquareImage(id int) {
	ci := CardImages.Objects[id].(*common.CardImage)
	if Board[id] == nil {
		ci.ReplaceImage(checkers.CheckerResource["empty"])
		fadeImage(ci, true)
	} else {
		ci.ReplaceImage(checkers.CheckerResource[Board[id].Icon])
		fadeImage(ci, false)
	}
	ci.Refresh()
}

func fadeImage(ci *common.CardImage, fade bool) {
	t := time.Millisecond * time.Duration(Ctx.fadeTime)
	if fade {
		common.DeflateImage(ci.Img, t)
	} else {
		common.InflateImage(ci.Img, t)
	}
	ci.Refresh()
}

func updatePlayerIndicator() {
	switch activePlayerSuit {
	case Red:
		playerTurn.Importance = widget.HighImportance
	case Black:
		playerTurn.Importance = widget.LowImportance
	}
	playerTurn.Refresh()
}

func borderImage(id int, show bool, color color.Color) {
	ci := CardImages.Objects[id].(*common.CardImage)
	ci.ShowBorder(show, color)
	ci.Refresh()
}

func dimSquare(id int, dim bool) {
	ci := CardImages.Objects[id].(*common.CardImage)
	ci.SetDimmed(dim)
	ci.Refresh()
}
