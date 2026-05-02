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

  File:    rummy.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Main program for the rummy game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/rummy/rummy.exe ./cmd/rummy


in games/cmd/rummy:
fyne bundle -o splash.go --package main splash.png


in games/cmd/singles:
fyne package --release --icon=splash.png --id=com.scsi.rummy

IntelliJ ---------
Run Name: go build cmd/rummy
Run On: Local Machine
Run Kind : Directory
Directory: D:\charon\z\devl\projects\go\games\cmd\rummy
Working: D:\charon\z\devl\projects\go\games\cmd\rummy
Module: game
*/

func main() {
	singles := app.NewWithID("com.scsi.rummy")
	Ctx.app = singles
	singles.Settings().SetTheme(rummyTheme{})

	storage := singles.Storage().RootURI().Path()
	mainWindow := singles.NewWindow("Rummy (by Bob)")

	setTheme()

	mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "rummy") + ".log"
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
	mainWindow.Resize(fyne.NewSize(700, 500))

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
	singles.Quit()
}
