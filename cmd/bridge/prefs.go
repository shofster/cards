package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
)

/*

  File:    prefs.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Settings and Globals
*/

type Global struct {
	app           fyne.App
	theme         fyne.Theme
	width, height float32
	offset        int
}

var Ctx Global

const contractsBidKey = "bid"
const contractsWonKey = "won"
const showEastKey = "show-east"
const showWestKey = "show-west"

var contractLevels = []string{"1", "2", "3", "4", "5", "6", "7"}
var contractSuits = []string{"♣", "♦", "♥", "♠", "NoTrump"}

// global bound variables
var contractLevel = binding.NewString()
var contractSuit = binding.NewString()
var contract = binding.NewString()
var dealResult = binding.NewString()
var aboveLine = binding.NewInt()
var belowLine = binding.NewInt()

var contractsBid = binding.NewInt()
var contractsWon = binding.NewInt()
var showEast = binding.NewBool()
var showWest = binding.NewBool()

// initialize the global (and storage) variables.
func initBoundVariables() {

	contractsBid.AddListener(binding.NewDataListener(func() {
		v, _ := contractsBid.Get()
		Ctx.app.Preferences().SetInt(contractsBidKey, v)
	}))
	contractsWon.AddListener(binding.NewDataListener(func() {
		v, _ := contractsWon.Get()
		Ctx.app.Preferences().SetInt(contractsWonKey, v)
	}))

	_ = contractLevel.Set("")
	_ = contractSuit.Set("")
	_ = contract.Set("?")
	_ = dealResult.Set("  ----  RESULTS  ----")
	_ = contractsBid.Set(Ctx.app.Preferences().IntWithFallback(contractsBidKey, 0))
	_ = contractsWon.Set(Ctx.app.Preferences().IntWithFallback(contractsWonKey, 0))
	_ = showEast.Set(Ctx.app.Preferences().BoolWithFallback(showEastKey, false))
	_ = showEast.Set(false)
	_ = showWest.Set(Ctx.app.Preferences().BoolWithFallback(showWestKey, false))
	_ = showWest.Set(false)
	_ = aboveLine.Set(0)
	_ = belowLine.Set(0)

}

// window size and theme modifications

func setTheme() {
	p := Ctx.app.Preferences()
	Ctx.width = float32(p.FloatWithFallback("sizeWidth", 80))
	Ctx.height = float32(p.FloatWithFallback("sizeHeight", 100))
	Ctx.offset = p.IntWithFallback("offset", 30)
	p.SetFloat("sizeWidth", float64(Ctx.width))
	p.SetFloat("sizeHeight", float64(Ctx.height))
	p.SetInt("offset", Ctx.offset)
	Ctx.theme = &singlesTheme{}

}

var _ fyne.Theme = (*singlesTheme)(nil)

type singlesTheme struct {
}

func (j singlesTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j singlesTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j singlesTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j singlesTheme) Size(name fyne.ThemeSizeName) float32 {
	base := theme.DefaultTheme().Size(name)
	switch name {
	case theme.SizeNamePadding:
		return 1 // outer padding
	case theme.SizeNameInnerPadding:
		return 2 // spacing between label and widget
	case theme.SizeNameText:
		return base + 4
	}
	return base
}
