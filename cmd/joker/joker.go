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

  File:    joker.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/
/*
  Description: Main program for joker wild game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/joker/joker.exe ./cmd/joker


in games/cmd/joker:
fyne bundle -o splash.go --package main splash.png


in games/cmd/joker:
fyne package --release --icon=splash.png --id=com.scsi.joker

IntelliJ ---------
Run Name: go build cmd/joker
Run On: Local Machine
Run Kind : Directory
Directory: D:\charon\z\devl\projects\go\games\cmd\joker
Working: D:\charon\z\devl\projects\go\games\cmd\joker
Module: game
*/

func main() {
	joker := app.NewWithID("com.scsi.joker")
	Ctx.app = joker
	joker.Settings().SetTheme(jokerTheme{})

	storage := joker.Storage().RootURI().Path()
	mainWindow := joker.NewWindow("JOKER's Wild (by Bob)")

	setTheme()

	mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "joker") + ".log"
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
	mainWindow.SetContent(splash)
	mainWindow.Resize(fyne.NewSize(300, 400))

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

	mainWindow.ShowAndRun()
	joker.Quit()
}
