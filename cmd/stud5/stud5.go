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

  File:    stud5.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Main program for 5 card stud game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/stud5/stud5.exe ./cmd/stud5


in games/cmd/stud5:
fyne bundle -o splash.go --package main splash.png


in games/cmd/stud5:
fyne package --release --icon=splash.png --id=com.scsi.stud5

IntelliJ ---------
Run Name: go build cmd/stud5
Run On: Local Machine
Run Kind : Directory
Directory: D:\charon\z\devl\projects\go\games\cmd\stud5
Working: D:\charon\z\devl\projects\go\games\cmd\stud5
Module: game
*/

func main() {
	stud5 := app.NewWithID("com.scsi.stud5")
	Ctx.app = stud5
	stud5.Settings().SetTheme(draw5Theme{})

	storage := stud5.Storage().RootURI().Path()
	mainWindow := stud5.NewWindow("5 Card Stud (by Bob)")

	setTheme()

	mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "stud5") + ".log"
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
	mainWindow.Resize(fyne.NewSize(500, 350))

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
	stud5.Quit()
}
