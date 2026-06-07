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

  File:    cribbage.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
 Description: Main program for Cribbage game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/cribbage/cribbage.exe ./cmd/cribbage


in games/cmd/cribbage:
fyne bundle -o splash.go --package main splash.png


in games/cmd/cribbage:
fyne package --release --icon=splash.png --id=com.scsi.cribbage

IntelliJ ---------
Run Name: go build cmd/cribbage
Run On: Local Machine
Run Kind : Directory
Directory: D:\charon\z\devl\projects\go\games\cmd\cribbage
Working: D:\charon\z\devl\projects\go\games\cmd\cribbage
Module: games

*/

func main() {
	cribbage := app.NewWithID("com.scsi.cribbage")
	Ctx.app = cribbage
	cribbage.Settings().SetTheme(cribbageTheme{})

	storage := cribbage.Storage().RootURI().Path()
	Ctx.mainWindow = cribbage.NewWindow("Cribbage (by Bob)")

	setTheme()

	Ctx.mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "cribbage") + ".log"
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
	Ctx.mainWindow.Resize(fyne.NewSize(900, 500))

	go func() { // show photo for a bit
		initGlobalVariables()
		content := initUI()
		time.Sleep(time.Millisecond * 3000)
		fyne.Do(func() {
			splash.RemoveAll()
			splash.Add(content)
			//mainWindow.Resize(fyne.NewSize(800, 500))
			splash.Refresh()
		})
	}()

	Ctx.mainWindow.ShowAndRun()
	cribbage.Quit()
}
