package main

import (
	"cards/common"
	"fmt"

	"fyne.io/fyne/v2"
)

/*

  File:    stock.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: manage the clickability for PLAYER
*/

const DRAW_ID = -2

func enableStock(b bool, cb func(int)) {
	defer fyne.Do(func() {
		stockSlot.Refresh()
	})
	si := stockSlot.Objects[0].(*common.CardImage)
	switch b {
	case true: // show clickable
		si.Enabled = true
		si.SetHighlite(true)
		si.OnTap = func(ci *common.CardImage) {
			cb(DRAW_ID)
		}
	case false: // hide clickable
		si.Enabled = false
		si.SetHighlite(false)
		si.OnTap = nil
	}
}

func hiliteStock(b bool) {
	defer fyne.Do(func() {
		stockSlot.Refresh()
	})
	si := stockSlot.Objects[0].(*common.CardImage)
	si.Enabled = false
	si.SetHighlite(b)
	si.OnTap = nil
}

func stockImage(code string) {
	defer fyne.Do(func() {
		stockSlot.Refresh()
	})
	si := stockSlot.Objects[0].(*common.CardImage)
	si.Replace(getTileResource(code))
}

func stockDraw() *Tile {
	if len(stock) == 0 {
		stockImage("empty")
		return nil
	}
	defer fyne.Do(func() {
		stockSlot.Refresh()
	})
	si := stockSlot.Objects[0].(*common.CardImage)
	// current tile
	t := stock[0]
	stock = stock[1:] // set next
	si.Replace(getTileResource(t.Code))
	Logit(func() string { return fmt.Sprintf("stockDraw: card is %s\n", t.Code) })
	return t
}
