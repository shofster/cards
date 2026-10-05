package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

/*

  File:    mahjong.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
   Description: Main program for the mahjong game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/mahjong/mahjong.exe ./cmd/mahjong


in games/cmd/mahjong:
fyne bundle -o splash.go --package main splash.png


in games/cmd/mahjong:
fyne package --release --icon=splash.png --id=com.scsi.mahjong

IntelliJ ---------
Run Name: go build mahjong
Run On: Local Machine
Run Kind : Package
Package Path : cards/cmd/mahjong
Working: D:\charon\z\devl\projects\go\cards
Module: cards

*/

func main() {
	mahjong := app.NewWithID("com.scsi.mahjong")
	Ctx.app = mahjong
	mahjong.Settings().SetTheme(mahjongTheme{})

	storage := mahjong.Storage().RootURI().Path()
	Ctx.mainWindow = mahjong.NewWindow("Mahjong 4-1 (by Bob)")

	setTheme()
	Ctx.devMode = false

	Ctx.mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "mahjong") + ".log"
	var logger *os.File
	if logf, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666); err == nil {
		logger = logf
		defer func() {
			_ = logger.Close()
		}()
		log.SetOutput(logger)
	}

	photo := canvas.NewImageFromResource(resourceSplashPng)
	photo.FillMode = canvas.ImageFillContain
	splash := container.NewStack(photo)
	Ctx.mainWindow.SetContent(splash)
	Ctx.mainWindow.Resize(fyne.NewSize(1124, 700))
	Ctx.mainWindow.CenterOnScreen()

	go func() { // show photo for a bit
		initBoundVariables()
		content := initUI()
		time.Sleep(time.Millisecond * 3000)
		fyne.Do(func() {
			splash.RemoveAll()
			splash.Add(content)
			splash.Refresh()
		})
	}()

	Ctx.mainWindow.ShowAndRun()
	mahjong.Quit()
}
