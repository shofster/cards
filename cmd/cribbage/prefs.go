package main

import (
	"cards/common"
	"image/color"
	"os/user"
	"strings"

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
  Description:
*/

const (
	PLAYER int = 0 // Player
	AI     int = 1 // AI
)

type Global struct {
	app                  fyne.App
	mainWindow           fyne.Window
	theme                fyne.Theme
	cardMap              map[string]common.Card
	cardSizeWidth        float32
	cardSizeHeight       float32
	cardOffset           int
	pegSize              float32
	trackSize            fyne.Size
	dealerName           binding.String
	playerName, aiName   string
	handPoints           binding.Int
	playerScore, aiScore binding.Int
	debug                bool
}

var Ctx Global

const userNameKey = "player"
const aiNameKey = "ai"
const pegSizeKey = "pegSize"
const trackWidthKey = "trackWidth"
const trackHeightKey = "trackHeight"
const debugKey = "debug"

// initialize the global (and storage) variables.
func initGlobalVariables() {
	u, e := user.Current()
	if e != nil {
		u.Username = "Igor"
	}
	Ctx.playerName = Ctx.app.Preferences().StringWithFallback(userNameKey,
		shortUsername(u.Username))
	Ctx.app.Preferences().SetString(userNameKey, Ctx.playerName)
	Ctx.aiName = Ctx.app.Preferences().StringWithFallback(aiNameKey, "Suckling")
	Ctx.app.Preferences().SetString(aiNameKey, Ctx.aiName)
	Ctx.debug = Ctx.app.Preferences().BoolWithFallback(debugKey, false)

	Ctx.cardMap = common.CardMap() // icon vs Card

	s := Ctx.app.Preferences().IntWithFallback(pegSizeKey, 20)
	Ctx.app.Preferences().SetInt(pegSizeKey, s)
	Ctx.pegSize = float32(s)

	w := Ctx.app.Preferences().IntWithFallback(trackWidthKey, 500)
	h := Ctx.app.Preferences().IntWithFallback(trackHeightKey, 500)
	Ctx.app.Preferences().SetInt(trackWidthKey, w)
	Ctx.app.Preferences().SetInt(trackHeightKey, h)
	Ctx.trackSize = fyne.Size{Width: float32(w), Height: float32(h)}

	Ctx.dealerName = binding.NewString()
	common.StringSet(Ctx.dealerName, "Cut Deck")
	Ctx.handPoints = binding.NewInt()
	Ctx.playerScore = binding.NewInt()
	Ctx.aiScore = binding.NewInt()

}

func shortUsername(full string) string {
	i := strings.LastIndex(full, "\\")
	if i == -1 {
		return full
	}
	return full[i+1:]
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
	Ctx.theme = &cribbageTheme{}

}

var _ fyne.Theme = (*cribbageTheme)(nil)

type cribbageTheme struct {
}

func (j cribbageTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j cribbageTheme) Font(style fyne.TextStyle) fyne.Resource {
	// return theme.DefaultTheme().Font(fyne.TextStyle{Monospace: true})
	return theme.DefaultTheme().Font(style)
}

func (j cribbageTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j cribbageTheme) Size(name fyne.ThemeSizeName) float32 {
	base := theme.DefaultTheme().Size(name)
	switch name {
	case theme.SizeNameText:
		return base * 1.3 // or 1.3, 1.5 etc.
	}
	return base
}
