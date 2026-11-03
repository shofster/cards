package main

import (
	"image/color"
	"os/user"
	"strings"

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
	app        fyne.App
	mainWindow fyne.Window
	theme      fyne.Theme
	playerName string
	botName    string
	squareSize int
	cursorSize int
	fadeTime   int
	devMode    bool
}

var Ctx Global

// initialize the global (and storage) variables.
func initBoundVariables() {
}

// window sizes and theme modifications

func setTheme() {
	p := Ctx.app.Preferences()
	Ctx.squareSize = p.IntWithFallback("squareSize", 100)
	Ctx.devMode = p.BoolWithFallback("dev-mode", true)
	Ctx.cursorSize = p.IntWithFallback("cursorSize", 60)
	Ctx.fadeTime = p.IntWithFallback("fadeTime", 3000)
	p.SetInt("cursorSize", Ctx.cursorSize)
	p.SetInt("squareSize", Ctx.squareSize)
	p.SetInt("fadeTime", Ctx.fadeTime)
	Ctx.botName = p.StringWithFallback("botName", "Checkers")
	p.SetString("botName", Ctx.botName)
	u, err := user.Current()
	Ctx.playerName = "Player"
	if err == nil && u.Username != "" {
		Ctx.playerName = u.Username
	}
	Ctx.playerName = shortUsername(Ctx.playerName)
	Ctx.theme = &checkersTheme{}
}

func shortUsername(full string) string {
	i := strings.LastIndex(full, "\\")
	if i == -1 {
		return full
	}
	return full[i+1:]
}

var _ fyne.Theme = (*checkersTheme)(nil)

type checkersTheme struct {
}

func (j checkersTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
	//	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j checkersTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j checkersTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j checkersTheme) Size(name fyne.ThemeSizeName) float32 {
	base := theme.DefaultTheme().Size(name)
	switch name {
	case theme.SizeNamePadding:
		return 1 // outer padding
	case theme.SizeNameInnerPadding:
		return 2 // spacing between label and widget
	case theme.SizeNameText:
		return base + 8
		//case theme.SizeNameHeadingText:
		// return float32(Ctx.cursorSize)
	}
	return base
}
