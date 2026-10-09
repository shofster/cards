package main

import "C"
import (
	"image/color"

	"fyne.io/fyne/v2"
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
	app                   fyne.App
	mainWindow            fyne.Window
	theme                 fyne.Theme
	width, height         float32
	deckWidth, deckHeight float32
	meldWidth, meldHeight float32
	offset                int
	tileTheme             string
	showAItiles           bool
	devMode               bool
	showHints             bool
}

var Ctx Global

// initialize the global (and storage) variables.
func initBoundVariables() {
}

const REGULAR = "regular"
const BLACK = "black"

// window sizes and theme modifications

func setTheme() {
	p := Ctx.app.Preferences()
	Ctx.width = float32(p.FloatWithFallback("sizeWidth", 80))
	Ctx.height = float32(p.FloatWithFallback("sizeHeight", 100))
	Ctx.deckWidth = float32(p.FloatWithFallback("deck_sizeWidth", 60))
	Ctx.deckHeight = float32(p.FloatWithFallback("deck_sizeHeight", 75))
	Ctx.meldWidth = float32(p.FloatWithFallback("meld_sizeWidth", 50))
	Ctx.meldHeight = float32(p.FloatWithFallback("meld_sizeHeight", 20))
	Ctx.offset = p.IntWithFallback("offset", 30)
	Ctx.tileTheme = p.StringWithFallback("tile_theme", REGULAR)
	Ctx.devMode = p.BoolWithFallback("dev-mode", true)
	Ctx.showAItiles = p.BoolWithFallback("show-ai", false)
	Ctx.showHints = p.BoolWithFallback("show-hints", true)
	p.SetFloat("sizeWidth", float64(Ctx.width))
	p.SetFloat("sizeHeight", float64(Ctx.height))
	p.SetInt("offset", Ctx.offset)
	p.SetString("tile_theme", Ctx.tileTheme)
	p.SetBool("show-ai", Ctx.showAItiles)
	p.SetBool("show-hints", Ctx.showHints)
	//Ctx.tileTheme = "black"

	Ctx.theme = &mahjongTheme{}
}

var _ fyne.Theme = (*mahjongTheme)(nil)

type mahjongTheme struct {
}

func (j mahjongTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
	//	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j mahjongTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j mahjongTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j mahjongTheme) Size(name fyne.ThemeSizeName) float32 {
	base := theme.DefaultTheme().Size(name)
	switch name {
	case theme.SizeNamePadding:
		return 1 // outer padding
	case theme.SizeNameInnerPadding:
		return 2 // spacing between label and widget
	case theme.SizeNameText:
		return base + 12
	}
	return base
}
