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

  File:    chess.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

   Description: Main program for the chess game.

in cards:
go build -v -trimpath -ldflags="-s -w -H=windowsgui" -o cmd/chess/mahjong.exe ./cmd/chess


in cards/cmd/chess:
fyne bundle -o splash.go --package main splash.png
fyne bundle -o squares.go --package main board.png


in cards/cmd/chess:
fyne package --release --icon=splash.png --id=com.scsi.chess

*/

func main() {
	chess := app.NewWithID("com.scsi.chess")
	Ctx.app = chess
	chess.Settings().SetTheme(chessTheme{})

	storage := chess.Storage().RootURI().Path()
	Ctx.mainWindow = chess.NewWindow("Chess (by Bob)")

	setTheme()
	Ctx.devMode = false // debug logging
	SquareSize = Ctx.squareSize

	Ctx.mainWindow.SetIcon(resourceSplashPng)

	file := filepath.Join(storage, "chess") + ".log"
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
	w := float32(Ctx.squareSize)*8 + 220
	h := float32(Ctx.squareSize)*8 + 1
	Ctx.mainWindow.Resize(fyne.NewSize(w, h))
	Ctx.mainWindow.CenterOnScreen()

	go func() { // show splash for a bit
		initBoundVariables()
		content := initUI() // takes a few seconds to gen the board
		time.Sleep(time.Millisecond * 2)
		fyne.Do(func() {
			splash.RemoveAll()
			splash.Add(content)
			splash.Refresh()
		})
	}()

	Ctx.mainWindow.ShowAndRun()
	chess.Quit()
}
