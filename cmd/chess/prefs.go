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
  Description: Settings and Globals
*/

type Global struct {
	app        fyne.App
	mainWindow fyne.Window
	theme      fyne.Theme
	playerName string
	botName    string
	squareSize int
	dragSize   int
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
	Ctx.dragSize = p.IntWithFallback("dragSize", 60)
	Ctx.fadeTime = p.IntWithFallback("fadeTime", 3000)
	p.SetInt("dragSize", Ctx.dragSize)
	p.SetInt("squareSize", Ctx.squareSize)
	p.SetInt("fadeTime", Ctx.fadeTime)
	Ctx.botName = p.StringWithFallback("botName", "Chesster")
	p.SetString("botName", Ctx.botName)
	u, err := user.Current()
	Ctx.playerName = "Player"
	if err == nil && u.Username != "" {
		Ctx.playerName = u.Username
	}
	Ctx.playerName = shortUsername(Ctx.playerName)
	Ctx.theme = &chessTheme{}
}

func shortUsername(full string) string {
	i := strings.LastIndex(full, "\\")
	if i == -1 {
		return full
	}
	return full[i+1:]
}

var _ fyne.Theme = (*chessTheme)(nil)

type chessTheme struct {
}

func (j chessTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
	//	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j chessTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j chessTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j chessTheme) Size(name fyne.ThemeSizeName) float32 {
	base := theme.DefaultTheme().Size(name)
	switch name {
	case theme.SizeNamePadding:
		return 1 // outer padding
	case theme.SizeNameInnerPadding:
		return 2 // spacing between label and widget
	case theme.SizeNameText:
		return base + 8
	}
	return base
}
