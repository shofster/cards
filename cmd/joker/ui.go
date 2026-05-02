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
var playerPointsText *widget.Label
var dealerPointsText *widget.Label
var dealerStack *fyne.Container
var playerStack *fyne.Container

func initUI() *fyne.Container {

	// buttons for game play
	betButton = widget.NewButton(
		fmt.Sprintf("Bet: %d", intGet(bet)), nil)
	betButton.Disable()
	nextButton = widget.NewButton("Next", nil)
	nextButton.Disable()
	dealButton = widget.NewButton("Deal", nil)
	streetText = widget.NewLabelWithData(streetData)
	streetText.Alignment = fyne.TextAlignCenter

	playerPointsText = widget.NewLabelWithData(playerPointsData)
	playerPointsText.Alignment = fyne.TextAlignCenter

	dealerPointsText = widget.NewLabelWithData(dealerPointsData)
	dealerPointsText.Alignment = fyne.TextAlignCenter

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
	walletLabel := widget.NewLabelWithData(walletData)
	walletLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	wallet.AddListener(binding.NewDataListener(func() {
		_ = walletData.Set(stringGet(wallet))
	}))

	hs := container.NewHBox(widget.NewLabel("Game: "),
		widget.NewLabelWithData(handState))

	top := container.NewHBox(layout.NewSpacer(), hs, layout.NewSpacer(),
		widget.NewLabel("  Bet: "), widget.NewLabelWithData(potData),
		widget.NewLabel("  Wallet: "), walletLabel, layout.NewSpacer())

	content := container.NewGridWithColumns(3)

	statusLabel := widget.NewLabel("Game Status")
	statusLabel.Alignment = fyne.TextAlignCenter
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}
	slt := container.NewVBox(layout.NewSpacer(), layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), widget.NewLabel("Street #"), streetText),
		container.NewHBox(widget.NewLabel("  Player Points  "), playerPointsText),
		container.NewHBox(widget.NewLabel("  Dealer Points  "), dealerPointsText))
	s := widget.NewLabel("Game Control")
	s.Alignment = fyne.TextAlignCenter
	s.TextStyle = fyne.TextStyle{Bold: true}
	flow := container.NewVBox(slt, s, dealButton, nextButton, betButton)
	control := container.NewBorder(statusLabel, nil, nil, nil, flow)
	content.Add(control)

	//cardMap := common.CardMap()

	d1 := canvas.NewImageFromResource(cardMap["5♠"].Resource)
	d2 := canvas.NewImageFromResource(cardMap["J♥"].Resource)
	d3 := canvas.NewImageFromResource(cardMap["8♦"].Resource)
	d4 := canvas.NewImageFromResource(cardMap["8♣"].Resource)
	d5 := canvas.NewImageFromResource(cardMap["K♣"].Resource)
	d1.FillMode = canvas.ImageFillContain
	d2.FillMode = canvas.ImageFillContain
	d3.FillMode = canvas.ImageFillContain
	d4.FillMode = canvas.ImageFillContain
	d5.FillMode = canvas.ImageFillContain

	dealerStack = container.New(&OverlapLayout{Offset: float32(Ctx.offset)},
		d1, d2, d3, d4, d5)
	dealerLabel := widget.NewLabel("Dealer")
	dealerLabel.Alignment = fyne.TextAlignCenter
	dealerLabel.TextStyle = fyne.TextStyle{Bold: true}
	dealer := container.NewVBox(dealerLabel,
		container.NewStack(dealerStack))
	content.Add(dealer)

	p1 := canvas.NewImageFromResource(cardMap["A♥"].Resource)
	p2 := canvas.NewImageFromResource(cardMap["10♣"].Resource)
	p3 := canvas.NewImageFromResource(cardMap["3♠"].Resource)
	p4 := canvas.NewImageFromResource(cardMap["7♦"].Resource)
	p5 := canvas.NewImageFromResource(cardMap["joker"].Resource)
	p1.FillMode = canvas.ImageFillContain
	p2.FillMode = canvas.ImageFillContain
	p3.FillMode = canvas.ImageFillContain
	p4.FillMode = canvas.ImageFillContain
	p5.FillMode = canvas.ImageFillContain
	playerStack = container.New(&OverlapLayout{Offset: float32(Ctx.offset)},
		p1, p2, p3, p4, p5)
	playerLabel := widget.NewLabel("Player")
	playerLabel.Alignment = fyne.TextAlignCenter
	playerLabel.TextStyle = fyne.TextStyle{Bold: true}
	player := container.NewVBox(playerLabel,
		container.NewStack(playerStack))
	content.Add(player)

	board := container.NewBorder(top, nil, nil, nil,
		content)
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
		betButton.Text = fmt.Sprintf("Bet: %d", game.betAmount)
	}))
	betSelect = widget.NewSelectWithData(betAmounts, betData)

	walletData := binding.NewString()
	_ = walletData.Set(stringGet(wallet))
	walletLabel := widget.NewLabelWithData(walletData)
	walletLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	wallet.AddListener(binding.NewDataListener(func() {
		_ = walletData.Set(stringGet(wallet))
	}))
	addSelect := widget.NewSelect(walletAdditions, func(s string) {
		add, _ := strconv.Atoi(s)
		v := addToWallet(add)
		_ = walletData.Set(strconv.Itoa(v))
	})
	playerPointsData := binding.NewString()
	_ = playerPointsData.Set(stringGet(playerPoints))
	playerPointsLabel := widget.NewLabelWithData(playerPointsData)
	playerPointsLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	playerPoints.AddListener(binding.NewDataListener(func() {
		_ = playerPointsData.Set(stringGet(playerPoints))
	}))

	dealerPointsData := binding.NewString()
	_ = dealerPointsData.Set(stringGet(dealerPoints))
	dealerPointsLabel := widget.NewLabelWithData(dealerPointsData)
	dealerPointsLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	dealerPoints.AddListener(binding.NewDataListener(func() {
		_ = dealerPointsData.Set(stringGet(dealerPoints))
	}))

	var action = func(s binding.Bool) string {
		b, _ := s.Get()
		if b {
			return "Optional"
		}
		return "Disallowed"
	}

	settings := widget.NewForm(
		widget.NewFormItem("Ante Amount / Deal", anteSelect),
		widget.NewFormItem("Bet Amount / Street", betSelect),
		widget.NewFormItem("", widget.NewLabel("Rules:")),
		widget.NewFormItem(
			"Card Points",
			widget.NewLabel(
				"Ace->Six, Eight->Ten = 1->6, 8->10 Points"+
					"\nSeven = 4 Points, Steal 3 Points"+
					"\nJack  = Steal 3 Points"+
					"\nQueen = 12 Points"+
					"\nKing  = Double Next Card Points"+
					"\nJoker = Current Hand Points")),
		widget.NewFormItem("Street 1 Bet", widget.NewLabel(action(street1bet))),
		widget.NewFormItem("Street 2 Bet", widget.NewLabel(action(street2bet))),
		widget.NewFormItem("Street 3 Bet", widget.NewLabel(action(street3bet))),
		widget.NewFormItem("Street 4 Bet", widget.NewLabel(action(street4bet))),
		widget.NewFormItem("", widget.NewLabel("Player:")),
		widget.NewFormItem("Player Wallet", walletLabel),
		widget.NewFormItem("Add Amount to Wallet", addSelect),
		widget.NewFormItem("Accumulated Points", playerPointsLabel),
		widget.NewFormItem("Dealer Points", dealerPointsLabel))

	return container.NewStack(container.NewHBox(settings, layout.NewSpacer()))
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
func setCardFace(dealer bool, street int, cell cardCell, obverse *canvas.Image) {
	stack := playerStack
	if dealer {
		stack = dealerStack
	}

	img := stack.Objects[street].(*canvas.Image)
	if obverse != nil {
		img.Resource = obverse.Resource

	} else {
		img.Resource = cell.image.Resource

	}

	img.Refresh()
	stack.Refresh()
}

type OverlapLayout struct {
	Offset float32
}

func (l *OverlapLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n == 0 {
		return
	}

	// Compute max card height that fits the stack
	usableHeight := size.Height - float32(n-1)*l.Offset
	if usableHeight < 80 { // prevent collapse
		usableHeight = 80
	}

	// Maintain 5:7 aspect ratio ( for card image, visually OK)
	cardHeight := usableHeight
	cardWidth := cardHeight * (5.0 / 7.0)

	// Clamp width if window is narrow
	if cardWidth > size.Width {
		cardWidth = size.Width
		cardHeight = cardWidth * (7.0 / 5.0)
	}

	// Center horizontally
	x := (size.Width - cardWidth) / 2

	for i, o := range objects {
		o.Resize(fyne.NewSize(cardWidth, cardHeight))
		o.Move(fyne.NewPos(x, float32(i)*l.Offset))
	}
}

// MinSize - objects not used -- fixed size. will always be 5 in this app
func (l *OverlapLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	// Must be non-zero or Border center collapses
	return fyne.NewSize(Ctx.width, Ctx.height)
}
