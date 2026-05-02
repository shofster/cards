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

  File:    bridge.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Main program for the bridge game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/bridge/bridge.exe ./cmd/bridge


in games/cmd/bridge:
fyne bundle -o splash.go --package main splash.png


in games/cmd/singles:
fyne package --release --icon=splash.png --id=com.scsi.bridge

IntelliJ ---------
Run Name: go build cmd/bridge
Run On: Local Machine
Run Kind : Directory
Directory: D:\charon\z\devl\projects\go\games\cmd\bridge
Working: D:\charon\z\devl\projects\go\games\cmd\bridge
Module: game
*/

func main() {
	singles := app.NewWithID("com.scsi.bridge")
	Ctx.app = singles
	singles.Settings().SetTheme(singlesTheme{})

	storage := singles.Storage().RootURI().Path()
	mainWindow := singles.NewWindow("Bridge (by Bob)")

	setTheme()

	mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "bridge") + ".log"
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
	mainWindow.Resize(fyne.NewSize(800, 550))

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
