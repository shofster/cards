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

  File:    freecell.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Main program for FreeCell game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/freecell/freecell.exe ./cmd/freecell


in games/cmd/freecell:
fyne bundle -o splash.go --package main splash.png


in games/cmd/freecell:
fyne package --release --icon=splash.png --id=com.scsi.freecell

IntelliJ ---------
Run Name: go build cmd/freecell
Run On: Local Machine
Run Kind : Directory
Directory: D:\charon\z\devl\projects\go\games\cmd\freecell
Working: D:\charon\z\devl\projects\go\games\cmd\freecell
Module: cards

*/

func main() {
	freecell := app.NewWithID("com.scsi.freecell")
	Ctx.app = freecell
	freecell.Settings().SetTheme(freecellTheme{})

	storage := freecell.Storage().RootURI().Path()
	Ctx.mainWindow = freecell.NewWindow("FreeCell (by Bob)")

	setTheme()

	Ctx.mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "freecell") + ".log"
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
	freecell.Quit()
}
