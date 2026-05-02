package main

import (
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

func initUI() *fyne.Container {

	settingsTab := container.NewTabItem("Settings", initSettingsUI())
	settingsTab.Icon = theme.SettingsIcon()
	paytableTab := container.NewTabItem("Payouts", initPaytableUI())
	paytableTab.Icon = theme.DocumentIcon()
	title := " Play   ♣" + " ♦" + " ♥" + " ♠"
	gameTab := container.NewTabItem(title, makeGreenPanel(initGameUI()))
	gameTab.Icon = theme.MediaPlayIcon()
	tabs := container.NewAppTabs(settingsTab, paytableTab, gameTab)

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

// build the auto refresh settings configuration.
func initSettingsUI() *fyne.Container {

	walletLabel := widget.NewLabelWithData(wallet)
	walletLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	game.betAmount = widget.NewSelectWithData(betAmounts, bet)
	game.payoutType = widget.NewSelectWithData(payTables, payout)

	settings := widget.NewForm(
		widget.NewFormItem("Payout Type", game.payoutType),
		widget.NewFormItem("        Bet Amount / Hand", game.betAmount),
		widget.NewFormItem("", widget.NewLabel("")),
		widget.NewFormItem("Current Wallet", walletLabel),
		widget.NewFormItem("Add Amount to Wallet",
			widget.NewSelectWithData(walletAdditions, addWalletAdditions)),
		widget.NewFormItem("", widget.NewLabel("")),
		widget.NewFormItem("Statistics:", widget.NewLabel("")),
		widget.NewFormItem("Hands Played", widget.NewLabelWithData(handsPlayed)),
		widget.NewFormItem("Winners", widget.NewLabelWithData(handsWon)),
	)

	return container.NewStack(settings)
}

// build the auto refresh table of payouts.
func initPaytableUI() *fyne.Container {
	grid := container.NewGridWithColumns(6)
	payout.AddListener(binding.NewDataListener(func() {
		v, _ := payout.Get()
		setPayTable(PayTables[v], grid)
	}))
	v, _ := payout.Get()
	setPayTable(PayTables[v], grid)
	po := widget.NewLabelWithData(payout)
	po.TextStyle = fyne.TextStyle{
		Bold:   true,
		Italic: true,
	}
	return container.NewBorder(po, nil, nil, nil, grid)
}

func setPayTable(table map[HandRank][]int, grid *fyne.Container) {
	grid.RemoveAll()
	grid.Add(widget.NewLabel("Rank"))
	for i := 1; i < 6; i++ {
		grid.Add(widget.NewLabel("  " + strconv.Itoa(i)))
	}
	for hr := RoyalFlush; hr >= JacksOrBetter; hr-- {
		t := widget.NewLabel(hr.String())
		t.TextStyle = fyne.TextStyle{Bold: true}
		grid.Add(t)
		for _, v := range table[hr] {
			grid.Add(widget.NewLabel(strconv.Itoa(v)))
		}
	}
}

// build the widgets for the actual game.
// initGame must previously be called.
func initGameUI() *fyne.Container {

	top := container.NewHBox(layout.NewSpacer(),
		widget.NewLabel("Payout:"), widget.NewLabelWithData(payout),
		widget.NewLabel("    Bet:"), widget.NewLabelWithData(bet),
		widget.NewLabel("    Wallet:"), widget.NewLabelWithData(wallet))

	dealButton = widget.NewButton("DEAL", func() {
		startGame()
		deal()
	})
	drawButton = widget.NewButton("DRAW", func() {
		draw()
	})
	drawButton.Disable()

	result := widget.NewLabelWithData(handResult)
	result.TextStyle = fyne.TextStyle{
		Bold:      true,
		Monospace: false,
	}
	bottom := container.NewBorder(nil, nil, dealButton, drawButton,
		container.NewCenter(result))

	// build the 5 card cells for a hand
	content := container.NewGridWithColumns(5)
	for h := 0; h < 5; h++ {
		check := widget.NewCheck("keep", nil)
		cont, image := initCell(check)
		check.OnChanged = func(v bool) {
			setCardFace(image, "joker", v)
		}
		content.Add(cont)
		cardCells[h].image = image
		cardCells[h].action = check
	}

	board := container.NewBorder(top, bottom, nil, nil,
		content)
	return board
}

func initCell(check *widget.Check) (*fyne.Container, *fyne.Container) {
	resource := canvas.NewImageFromResource(cardMap["reverse"].Resource)
	image := container.NewStack(resource)
	return container.NewBorder(nil, container.NewCenter(check), nil, nil, image), image
}

func setCardFace(image *fyne.Container, icon string, obverse bool) {
	image.RemoveAll()
	if obverse {
		image.Add(canvas.NewImageFromResource(cardMap[icon].Resource))
	} else {
		image.Add(canvas.NewImageFromResource(cardMap["reverse"].Resource))
	}
	image.Refresh()
}
