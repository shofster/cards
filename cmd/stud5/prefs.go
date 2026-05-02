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

const walletKey = "wallet"
const anteKey = "ante"
const betKey = "bet"
const qualifyKey = "qualify"
const street1Key = "street1"
const street2Key = "street2"
const street3Key = "street3"

// global bound variables
var wallet = binding.NewInt()
var ante = binding.NewInt()
var bet = binding.NewInt()
var addWalletAddition = binding.NewInt()
var qualify = binding.NewBool() // true -> dealer has Ace,King high to qualify
var handState = binding.NewString()
var street1bet = binding.NewBool()
var street2bet = binding.NewBool()
var street3bet = binding.NewBool()

const defaultWallet = 100
const defaultBet = 2

var walletAdditions = []string{"10", "20", "50"}
var betAmounts = []string{"1", "2", "3", "4", "5"}

// initialize the global (and storage) variables.
func initBoundVariables() {
	wallet.AddListener(binding.NewDataListener(func() {
		v, _ := wallet.Get()
		Ctx.app.Preferences().SetInt(walletKey, v)
	}))
	ante.AddListener(binding.NewDataListener(func() {
		v, _ := ante.Get()
		Ctx.app.Preferences().SetInt(anteKey, v)
	}))
	bet.AddListener(binding.NewDataListener(func() {
		v, _ := bet.Get()
		Ctx.app.Preferences().SetInt(betKey, v)
	}))
	qualify.AddListener(binding.NewDataListener(func() {
		v, _ := qualify.Get()
		Ctx.app.Preferences().SetBool(qualifyKey, v)
	}))
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
	addWalletAddition.AddListener(binding.NewDataListener(func() {
		wal, _ := wallet.Get()
		add, _ := addWalletAddition.Get()
		_ = wallet.Set(wal + add)
	}))
	_ = wallet.Set(Ctx.app.Preferences().IntWithFallback(walletKey, defaultWallet))
	_ = ante.Set(Ctx.app.Preferences().IntWithFallback(anteKey, defaultBet))
	_ = bet.Set(Ctx.app.Preferences().IntWithFallback(betKey, defaultBet))
	_ = handState.Set("Press \"DEAL\" to begin")
	_ = qualify.Set(Ctx.app.Preferences().BoolWithFallback(qualifyKey, false))
	_ = street1bet.Set(Ctx.app.Preferences().BoolWithFallback(street1Key, true))
	_ = street2bet.Set(Ctx.app.Preferences().BoolWithFallback(street2Key, false))
	_ = street3bet.Set(Ctx.app.Preferences().BoolWithFallback(street3Key, true))
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

// window size and theme modifications

func setTheme() {
	p := Ctx.app.Preferences()
	Ctx.width = float32(p.FloatWithFallback("sizeWidth", 250))
	Ctx.height = float32(p.FloatWithFallback("sizeHeight", 400))
	Ctx.offset = p.IntWithFallback("offset", 45)
	p.SetFloat("sizeWidth", float64(Ctx.width))
	p.SetFloat("sizeHeight", float64(Ctx.height))
	p.SetInt("offset", Ctx.offset)
	Ctx.theme = &draw5Theme{}

}

var _ fyne.Theme = (*draw5Theme)(nil)

type draw5Theme struct {
}

func (j draw5Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j draw5Theme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j draw5Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j draw5Theme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
