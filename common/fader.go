package common

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

/*

  File:    fader.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: functions to dim or brighten an image.
*/

// brighten
func InflateImage(
	img *canvas.Image,
	duration time.Duration,
) {
	go func() {
		const fps = 60
		steps := int(duration.Seconds() * fps)
		delay := time.Second / fps
		for i := 0; i <= steps; i++ {
			t := float64(i) / float64(steps)
			fyne.Do(func() {
				img.Translucency = 1.0 - t
				img.Refresh()
			})
			time.Sleep(delay)
		}
	}()
}

// dim
func DeflateImage(
	img *canvas.Image,
	duration time.Duration,
) {
	go func() {
		const fps = 60
		steps := int(duration.Seconds() * fps)
		delay := time.Second / fps
		for i := 0; i <= steps; i++ {
			t := float64(i) / float64(steps)
			fyne.Do(func() {
				img.Translucency = t
				img.Refresh()
			})
			time.Sleep(delay)
		}
	}()
}
