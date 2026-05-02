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
}

var Ctx Global

const walletKey = "wallet"
const handsPlayedKey = "handsPlayed"
const handsWonKey = "handsWon"
const payoutKey = "payout"
const betKey = "bet"

// global bound variables
var wallet = binding.NewString()
var handsPlayed = binding.NewString()
var handsWon = binding.NewString()
var payout = binding.NewString()
var bet = binding.NewString()
var addWalletAdditions = binding.NewString()
var handResult = binding.NewString()

const defaultWallet = "1000"
const defaultPayout = "8/5"
const defaultBet = "1"

var walletAdditions = []string{"10", "50", "100"}
var betAmounts = []string{"1", "2", "3", "4", "5"}

// initialize the global (and storage) variables.
func initBoundVariables() {
	wallet.AddListener(binding.NewDataListener(func() {
		v, _ := wallet.Get()
		Ctx.app.Preferences().SetString(walletKey, v)
	}))
	handsPlayed.AddListener(binding.NewDataListener(func() {
		v, _ := handsPlayed.Get()
		Ctx.app.Preferences().SetString(handsPlayedKey, v)
	}))
	handsWon.AddListener(binding.NewDataListener(func() {
		v, _ := handsWon.Get()
		Ctx.app.Preferences().SetString(handsWonKey, v)
	}))
	payout.AddListener(binding.NewDataListener(func() {
		v, _ := payout.Get()
		Ctx.app.Preferences().SetString(payoutKey, v)
	}))
	bet.AddListener(binding.NewDataListener(func() {
		v, _ := bet.Get()
		Ctx.app.Preferences().SetString(betKey, v)
	}))
	addWalletAdditions.AddListener(binding.NewDataListener(func() {
		wal, _ := wallet.Get()
		add, _ := addWalletAdditions.Get()
		amt, _ := strconv.Atoi(wal)
		inc, _ := strconv.Atoi(add)
		s := strconv.Itoa(amt + inc)
		_ = wallet.Set(s)
	}))
	_ = wallet.Set(Ctx.app.Preferences().StringWithFallback(walletKey, defaultWallet))
	_ = handsPlayed.Set(Ctx.app.Preferences().StringWithFallback(handsPlayedKey, "0"))
	_ = handsWon.Set(Ctx.app.Preferences().StringWithFallback(handsWonKey, "0"))
	_ = payout.Set(Ctx.app.Preferences().StringWithFallback(payoutKey, defaultPayout))
	_ = bet.Set(Ctx.app.Preferences().StringWithFallback(betKey, defaultBet))
	_ = handResult.Set("Press \"DEAL\" to begin")
}

func setWonHand(won bool) {
	np, _ := handsPlayed.Get()
	nPlayed, _ := strconv.Atoi(np)
	_ = handsPlayed.Set(strconv.Itoa(nPlayed + 1))
	if won {
		nw, _ := handsWon.Get()
		nWon, _ := strconv.Atoi(nw)
		_ = handsWon.Set(strconv.Itoa(nWon + 1))
	}
}

func addBetAmountToWallet(wagerAmount int) {
	if wagerAmount < 0 {
		setWonHand(false)
	}
	wal, _ := wallet.Get()
	walletAmount, _ := strconv.Atoi(wal)
	_ = wallet.Set(strconv.Itoa(walletAmount + wagerAmount))
	if wagerAmount > 0 {
		setWonHand(true)
	}
}

// window size and theme modifications

func setTheme() {
	p := Ctx.app.Preferences()
	Ctx.width = float32(p.FloatWithFallback("sizeWidth", 900))
	Ctx.height = float32(p.FloatWithFallback("sizeHeight", 400))
	p.SetFloat("sizeWidth", float64(Ctx.width))
	p.SetFloat("sizeHeight", float64(Ctx.height))
	Ctx.theme = &jacksTheme{}

}

var _ fyne.Theme = (*jacksTheme)(nil)

type jacksTheme struct {
}

func (j jacksTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (j jacksTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (j jacksTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (j jacksTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
