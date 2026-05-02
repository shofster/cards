package main

import (
	"image/color"
	"strconv"

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

const handsPlayedKey = "hands-played"
const handsScoreKey = "total-score"
const dealerKnockKey = "dealer-knock"

// global bound variables

var handsPlayed = binding.NewString()
var handsScore = binding.NewString()
var winnerData = binding.NewString()
var playerData = binding.NewInt()
var dealerData = binding.NewInt()
var handResult = binding.NewInt()
var outcomeData = binding.NewString()

var dealerKnock = 10
var devMode bool = false

// initialize the global (and storage) variables.
func initBoundVariables() {

	handsPlayed.AddListener(binding.NewDataListener(func() {
		v, _ := handsPlayed.Get()
		Ctx.app.Preferences().SetString(handsPlayedKey, v)
	}))
	handsScore.AddListener(binding.NewDataListener(func() {
		v, _ := handsScore.Get()
		Ctx.app.Preferences().SetString(handsScoreKey, v)
	}))
	handResult.AddListener(binding.NewDataListener(func() {
		r, _ := handResult.Get()
		s, _ := handsScore.Get()
		n, _ := strconv.Atoi(s)
		_ = handsScore.Set(strconv.Itoa(n + r))
	}))

	_ = handsPlayed.Set(Ctx.app.Preferences().StringWithFallback(handsPlayedKey, "0"))
	_ = handsScore.Set(Ctx.app.Preferences().StringWithFallback(handsScoreKey, "0"))
	dealerKnock = Ctx.app.Preferences().IntWithFallback(dealerKnockKey, 10)
	devMode = Ctx.app.Preferences().BoolWithFallback("dev-mode", false)
}

func incHandsPlayed() {
	s, _ := handsPlayed.Get()
	n, _ := strconv.Atoi(s)
	_ = handsPlayed.Set(strconv.Itoa(n + 1))
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
	Ctx.theme = &rummyTheme{}

}

var _ fyne.Theme = (*rummyTheme)(nil)

type rummyTheme struct {
}

func (j rummyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
	//	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j rummyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j rummyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j rummyTheme) Size(name fyne.ThemeSizeName) float32 {
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
