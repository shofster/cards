package main

import (
	"cards/common"
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
  Description:
*/

type Global struct {
	app            fyne.App
	mainWindow     fyne.Window
	theme          fyne.Theme
	cardMap        map[string]common.Card
	cardSizeWidth  float32
	cardSizeHeight float32
	cardOffset     int
	debug          bool
}

var Ctx Global

const debugKey = "debug"

// initialize the global (and storage) variables.
func initGlobalVariables() {
	Ctx.debug = Ctx.app.Preferences().BoolWithFallback(debugKey, false)
	Ctx.app.Preferences().SetBool(debugKey, Ctx.debug)
	Ctx.cardMap = common.CardMap() // icon vs Card
}

// window size and theme modifications

func setTheme() {
	p := Ctx.app.Preferences()
	Ctx.cardSizeWidth = float32(p.FloatWithFallback("cardSizeWidth", 90))
	Ctx.cardSizeHeight = float32(p.FloatWithFallback("cardSizeHeight", 110))
	Ctx.cardOffset = p.IntWithFallback("cardOffset", 30)
	p.SetFloat("cardSizeWidth", float64(Ctx.cardSizeWidth))
	p.SetFloat("cardSizeHeight", float64(Ctx.cardSizeHeight))
	p.SetInt("cardOffset", Ctx.cardOffset)
	Ctx.theme = &freecellTheme{}

}

var _ fyne.Theme = (*freecellTheme)(nil)

type freecellTheme struct {
}

func (j freecellTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j freecellTheme) Font(style fyne.TextStyle) fyne.Resource {
	// return theme.DefaultTheme().Font(fyne.TextStyle{Monospace: true})
	return theme.DefaultTheme().Font(style)
}

func (j freecellTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j freecellTheme) Size(name fyne.ThemeSizeName) float32 {
	base := theme.DefaultTheme().Size(name)
	//switch name {
	//case theme.SizeNameText:
	//	return base * 1.3 // or 1.3, 1.5 etc.
	//}
	return base
}
