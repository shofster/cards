package main

import (
	"cards/assets/chess"
	"cards/common"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
var chessBoard *fyne.Container
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
	player = NewPlayerPanel(Ctx.playerName, White, playerControls)
	botControls := container.NewHBox()
	botControls.Add(widget.NewLabel("Auto"))
	auto := common.NewToggle(func(b bool) {
		autoPlay = b
	})
	botControls.Add(auto)
	bot = NewPlayerPanel(Ctx.botName, Black, botControls)
	chessBoard = NewChessBoard()
	return container.NewBorder(nil, nil, player.content, bot.content, chessBoard)
}

type PlayerPanel struct {
	suit      Suit
	nCaptured int
	captured  []binding.String
	content   *fyne.Container
}

func NewPlayerPanel(title string, suit Suit, controls *fyne.Container) *PlayerPanel {
	pp := &PlayerPanel{
		suit:     suit,
		captured: make([]binding.String, 16),
	}
	// up to 16 captured pieces
	grid := container.NewGridWithColumns(4)
	for i := range 16 {
		pp.captured[i] = binding.NewString()
		g := widget.NewLabelWithData(pp.captured[i])
		grid.Add(g)
	}
	for i := range 8 {
		if suit == Black {
			common.StringSet(pp.captured[i], chess.ChessUnicode[whiteHonorsRow[i]])
		} else {
			common.StringSet(pp.captured[i], chess.ChessUnicode[blackHonorsRow[i]])
		}
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
	for i := range 16 {
		common.StringSet(pp.captured[i], " ")
	}
	pp.nCaptured = 0
}
func (pp *PlayerPanel) addCapture(piece string) {
	common.StringSet(pp.captured[pp.nCaptured], piece)
	pp.nCaptured++
}

// swap the Board images for source and destination
func swapBoardPieceImages(sourceId, destinationId int) {
	ci := CardImages.Objects[sourceId].(*common.CardImage)
	res := ci.Img.Resource
	fadeImage(ci, true)
	ci = CardImages.Objects[destinationId].(*common.CardImage)
	ci.ReplaceImage(res)
	ci.Refresh()
	fadeImage(ci, false)
}

func updatePlayerIndicator() {
	switch activePlayerSuit {
	case White:
		playerTurn.Importance = widget.HighImportance
	case Black:
		playerTurn.Importance = widget.LowImportance
	}
	playerTurn.Refresh()
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

func kingWarn() {
	// remove previous warning
	if blackKingWarn >= 0 {
		borderImage(blackKingWarn, false, common.BorderRed)
	}
	if whiteKingWarn >= 0 {
		borderImage(whiteKingWarn, false, common.BorderRed)
	}
	// show new one
	blackKingWarn = KingThreatened(Black)
	if blackKingWarn >= 0 {
		borderImage(blackKingWarn, true, common.BorderRed)
	}
	whiteKingWarn = KingThreatened(White)
	if whiteKingWarn >= 0 {
		borderImage(whiteKingWarn, true, common.BorderRed)
	}
}
