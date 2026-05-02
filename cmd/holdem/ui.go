package main

import (
	"fmt"
	"image/color"
	"strconv"

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

var anteSelect *widget.Select
var betSelect *widget.Select
var streetText *widget.Label
var dealerRow *fyne.Container
var flopRow *fyne.Container
var playerRow *fyne.Container

func initUI() *fyne.Container {

	// buttons for game play
	betButton = widget.NewButton(fmt.Sprintf("  Bet: %d  ", intGet(bet)), nil)
	//betButton.Disable()
	nextButton = widget.NewButton("    Next    ", nil)
	//nextButton.Disable()
	dealButton = widget.NewButton("    Deal    ", nil)
	streetText = widget.NewLabelWithData(streetData)
	streetText.Alignment = fyne.TextAlignCenter
	streetText.Alignment = fyne.TextAlignCenter

	initGame()

	settingsTab := container.NewTabItem("Settings", initSettingsUI())
	settingsTab.Icon = theme.SettingsIcon()
	title := " Play   ♣" + " ♦" + " ♥" + " ♠"
	gameTab := container.NewTabItem(title, makeGreenPanel(initGameUI()))
	gameTab.Icon = theme.MediaPlayIcon()
	tabs := container.NewAppTabs(settingsTab, gameTab)

	return container.NewStack(tabs)
}

func makeGreenPanel(content *fyne.Container) *fyne.Container {
	// Soft felt green — color.NRGBA{0, 160, 0, 255}
	// Bright emerald — color.NRGBA{20, 180, 20, 255}
	// Muted mint — color.NRGBA{60, 200, 120, 255}
	// Casino - color.NRGBA{0, 128, 0, 255}
	bg := canvas.NewRectangle(color.NRGBA{20, 180, 20, 50})
	return container.NewStack(bg, content)
}

// build the widgets for the actual game.
// initGame must previously be called.
func initGameUI() *fyne.Container {

	walletData := binding.NewString()
	_ = walletData.Set(stringGet(wallet))
	wallet.AddListener(binding.NewDataListener(func() {
		v, _ := wallet.Get()
		_ = walletData.Set(strconv.Itoa(v))
	}))
	walletLabel := widget.NewLabelWithData(walletData)
	walletLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	po := container.NewHBox(widget.NewLabel("Game: "),
		widget.NewLabelWithData(handState))
	top := container.NewHBox(layout.NewSpacer(), po, layout.NewSpacer(),
		widget.NewLabel("    Bet:"), widget.NewLabelWithData(potData),
		widget.NewLabel("    Wallet:"), walletLabel)

	//initial card images slots
	initCardFaces()

	content := container.NewGridWithRows(3, dealerRow, flopRow, playerRow)

	controlLabel := widget.NewLabel("Game Control")
	controlLabel.Alignment = fyne.TextAlignCenter
	controlLabel.TextStyle = fyne.TextStyle{Bold: true}
	bottom := container.NewHBox(widget.NewLabel("Street #"), streetText,
		layout.NewSpacer(), dealButton, layout.NewSpacer(), nextButton, betButton)

	board := container.NewBorder(top, bottom, nil, nil, content)
	return board
}

// build the auto refresh settings configuration.
func initSettingsUI() *fyne.Container {

	var bs2i = func(bs binding.String) int {
		s, _ := bs.Get()
		v, _ := strconv.Atoi(s)
		return v
	}

	anteData := binding.NewString()
	_ = anteData.Set(stringGet(ante))
	anteData.AddListener(binding.NewDataListener(func() {
		_ = ante.Set(bs2i(anteData))
		game.anteAmount = intGet(ante)
	}))
	anteSelect = widget.NewSelectWithData(betAmounts, anteData)

	betData := binding.NewString()
	_ = betData.Set(stringGet(bet))
	betData.AddListener(binding.NewDataListener(func() {
		_ = bet.Set(bs2i(betData))
		game.betAmount = intGet(bet)
	}))
	betSelect = widget.NewSelectWithData(betAmounts, betData)

	walletData := binding.NewString()
	_ = walletData.Set(stringGet(wallet))
	walletLabel := widget.NewLabelWithData(walletData)
	walletLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	addSelect := widget.NewSelect(walletAdditions, func(s string) {
		add, _ := strconv.Atoi(s)
		v := addToWallet(add)
		_ = walletData.Set(strconv.Itoa(v))
	})

	var action = func(s binding.Bool) string {
		b, _ := s.Get()
		if b {
			return "Optional"
		}
		return "Disallowed"
	}

	qual := "No"
	v, _ := qualify.Get()
	if v {
		qual = "Yes"
	}

	settings := widget.NewForm(
		widget.NewFormItem("Ante Amount / Hand", anteSelect),
		widget.NewFormItem("Bet Amount / Street", betSelect),
		widget.NewFormItem("", widget.NewLabel("Rules")),
		widget.NewFormItem("Dealer Must Qualify?", widget.NewLabel(qual)),
		widget.NewFormItem("Street 1 Bet", widget.NewLabel(action(street1bet))),
		widget.NewFormItem("Street 2 Bet", widget.NewLabel(action(street2bet))),
		widget.NewFormItem("Street 3 Bet", widget.NewLabel(action(street3bet))),
		widget.NewFormItem("", widget.NewLabel("Wallet")),
		widget.NewFormItem("Current Amount", walletLabel),
		widget.NewFormItem("Add Amount to Wallet", addSelect))

	return container.NewStack(settings)
}

func allowUpdateBets(b bool) {
	switch b {
	case true:
		anteSelect.Enable()
		betSelect.Enable()
	default:
		anteSelect.Disable()
		betSelect.Disable()
	}
}

// set image to card face or the reverse (if specified)
var d1 *canvas.Image
var d2 *canvas.Image
var p1 *canvas.Image
var p2 *canvas.Image
var f1 *canvas.Image
var f2 *canvas.Image
var f3 *canvas.Image
var j1 *canvas.Image // turn
var j2 *canvas.Image // river

func initCardFaces() {
	d1 = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	d2 = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	p1 = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	p2 = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	d1.FillMode = canvas.ImageFillContain
	d2.FillMode = canvas.ImageFillContain
	p1.FillMode = canvas.ImageFillContain
	p2.FillMode = canvas.ImageFillContain

	f1 = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	f2 = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	f3 = canvas.NewImageFromResource(cardMap["reverse"].Resource)
	f1.FillMode = canvas.ImageFillContain
	f2.FillMode = canvas.ImageFillContain
	f3.FillMode = canvas.ImageFillContain

	j1 = canvas.NewImageFromResource(cardMap["joker"].Resource)
	j2 = canvas.NewImageFromResource(cardMap["joker"].Resource)
	j1.FillMode = canvas.ImageFillContain
	j2.FillMode = canvas.ImageFillContain

	dealerRow = container.NewGridWithColumns(6, layout.NewSpacer(), layout.NewSpacer(), d1, d2,
		layout.NewSpacer(), layout.NewSpacer())
	flopRow = container.NewGridWithColumns(6, f1, f2, f3, layout.NewSpacer(), j1, j2)
	playerRow = container.NewGridWithColumns(6, layout.NewSpacer(), layout.NewSpacer(), p1, p2,
		layout.NewSpacer(), layout.NewSpacer())
}

func resetCardFaces() {
	d1.Resource = cardMap["reverse"].Resource
	d2.Resource = cardMap["reverse"].Resource
	p1.Resource = cardMap["reverse"].Resource
	p2.Resource = cardMap["reverse"].Resource
	f1.Resource = cardMap["reverse"].Resource
	f2.Resource = cardMap["reverse"].Resource
	f3.Resource = cardMap["reverse"].Resource
	j1.Resource = cardMap["reverse"].Resource
	j2.Resource = cardMap["reverse"].Resource
	dealerRow.Refresh()
	flopRow.Refresh()
	playerRow.Refresh()
}

func setCardFace(idx int, obverse *canvas.Image) {
	switch idx {
	case 0:

		d1.Resource = obverse.Resource
	case 1:
		d2.Resource = obverse.Resource

	case 2:
		p1.Resource = obverse.Resource
	case 3:
		p2.Resource = obverse.Resource

	case 4:
		f1.Resource = obverse.Resource
	case 5:
		f2.Resource = obverse.Resource
	case 6:
		f3.Resource = obverse.Resource

	case 7:
		j1.Resource = obverse.Resource
	case 8:
		j2.Resource = obverse.Resource
	}
	dealerRow.Refresh()
	flopRow.Refresh()
	playerRow.Refresh()
}
