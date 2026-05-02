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

  File:    jacks.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

  Description:  Main program for jacks or better game.

in games:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/jacks/jacks.exe ./cmd/jacks


in games/cmd/jacks:
fyne bundle -o splash.go --package main splash.png


in games/cmd/jacks:
fyne package --release --icon=splash.png --id=com.scsi.jacks

IntelliJ ---------
Run Name: go build cmd/jacks
Run On: Local Machine
Run Kind : Directory
Directory: D:\charon\z\devl\projects\go\games\cmd\jacks
Working: D:\charon\z\devl\projects\go\games\cmd\jacks
Module: games
*/

func main() {
	jacks := app.NewWithID("com.scsi.jacks")
	Ctx.app = jacks
	jacks.Settings().SetTheme(jacksTheme{})

	storage := jacks.Storage().RootURI().Path()
	mainWindow := jacks.NewWindow("Jacks or Better (by Bob)")

	setTheme()

	mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "jacks") + ".log"
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
	mainWindow.Resize(fyne.NewSize(Ctx.width, Ctx.height))

	go func() { // show photo for a bit
		initBoundVariables()
		initGame()
		content := initUI()
		time.Sleep(time.Millisecond * 3000)
		fyne.Do(func() {
			splash.Objects[0] = content
			splash.Refresh()
		})
	}()

	mainWindow.ShowAndRun()
	jacks.Quit()
}
