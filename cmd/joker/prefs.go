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

const anteKey = "ante"
const betKey = "bet"
const walletKey = "wallet"
const pointsDealerKey = "pointsDealer"
const pointsPlayerKey = "pointsPlayer"
const street1Key = "street1"
const street2Key = "street2"
const street3Key = "street3"
const street4Key = "street4"

// global bound variables
var ante = binding.NewInt()
var bet = binding.NewInt()
var handState = binding.NewString()
var wallet = binding.NewInt()
var dealerPoints = binding.NewInt()
var playerPoints = binding.NewInt()
var street1bet = binding.NewBool()
var street2bet = binding.NewBool()
var street3bet = binding.NewBool()
var street4bet = binding.NewBool()

const defaultWallet = 100
const defaultBet = 2

var walletAdditions = []string{"10", "20", "50"}
var betAmounts = []string{"1", "2", "3", "4", "5"}

// initialize the global (and storage) variables.
func initBoundVariables() {
	ante.AddListener(binding.NewDataListener(func() {
		v, _ := ante.Get()
		Ctx.app.Preferences().SetInt(anteKey, v)
	}))
	bet.AddListener(binding.NewDataListener(func() {
		v, _ := bet.Get()
		Ctx.app.Preferences().SetInt(betKey, v)
	}))
	_ = ante.Set(Ctx.app.Preferences().IntWithFallback(anteKey, defaultBet))
	_ = bet.Set(Ctx.app.Preferences().IntWithFallback(betKey, defaultBet))

	wallet.AddListener(binding.NewDataListener(func() {
		v, _ := wallet.Get()
		Ctx.app.Preferences().SetInt(walletKey, v)
	}))
	_ = wallet.Set(Ctx.app.Preferences().IntWithFallback(walletKey, defaultWallet))

	dealerPoints.AddListener(binding.NewDataListener(func() {
		v, _ := dealerPoints.Get()
		Ctx.app.Preferences().SetInt(pointsDealerKey, v)
	}))
	playerPoints.AddListener(binding.NewDataListener(func() {
		v, _ := playerPoints.Get()
		Ctx.app.Preferences().SetInt(pointsPlayerKey, v)
	}))
	_ = dealerPoints.Set(Ctx.app.Preferences().IntWithFallback(pointsDealerKey, 0))
	_ = playerPoints.Set(Ctx.app.Preferences().IntWithFallback(pointsPlayerKey, 0))

	street1bet.AddListener(binding.NewDataListener(func() {
		v, _ := street1bet.Get()
		Ctx.app.Preferences().SetBool(street1Key, v)
	}))
	street2bet.AddListener(binding.NewDataListener(func() {
		v, _ := street2bet.Get()
		Ctx.app.Preferences().SetBool(street2Key, v)
	}))
	street3bet.AddListener(binding.NewDataListener(func() {
		v, _ := street3bet.Get()
		Ctx.app.Preferences().SetBool(street3Key, v)
	}))
	street4bet.AddListener(binding.NewDataListener(func() {
		v, _ := street4bet.Get()
		Ctx.app.Preferences().SetBool(street4Key, v)
	}))
	_ = street1bet.Set(Ctx.app.Preferences().BoolWithFallback(street1Key, false))
	_ = street2bet.Set(Ctx.app.Preferences().BoolWithFallback(street2Key, true))
	_ = street3bet.Set(Ctx.app.Preferences().BoolWithFallback(street3Key, false))
	_ = street4bet.Set(Ctx.app.Preferences().BoolWithFallback(street4Key, true))
	_ = handState.Set("Press \"DEAL\" to begin")
	_ = potData.Set("0")
}

func stringGet(bi binding.Int) string {
	v, _ := bi.Get()
	return strconv.Itoa(v)
}

func intGet(bi binding.Int) int {
	v, _ := bi.Get()
	return v
}

func boolGet(bi binding.Bool) bool {
	v, _ := bi.Get()
	return v
}

func addToWallet(v int) int {
	walletAmount, _ := wallet.Get()
	n := walletAmount + v
	_ = wallet.Set(n)
	return n
}

func addPoints(bi binding.Int, v int) int {
	amount, _ := bi.Get()
	n := amount + v
	_ = bi.Set(n)
	return n
}

// window size and theme modifications

func setTheme() {
	p := Ctx.app.Preferences()
	Ctx.width = float32(p.FloatWithFallback("sizeWidth", 250))
	Ctx.height = float32(p.FloatWithFallback("sizeHeight", 500))
	Ctx.offset = p.IntWithFallback("offset", 60)
	p.SetFloat("sizeWidth", float64(Ctx.width))
	p.SetFloat("sizeHeight", float64(Ctx.height))
	p.SetInt("offset", Ctx.offset)
	Ctx.theme = &jokerTheme{}

}

var _ fyne.Theme = (*jokerTheme)(nil)

type jokerTheme struct {
}

func (j jokerTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j jokerTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j jokerTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j jokerTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 1 // outer padding
	case theme.SizeNameInnerPadding:
		return 2 // spacing between label and widget
	}
	return theme.DefaultTheme().Size(name)
}
