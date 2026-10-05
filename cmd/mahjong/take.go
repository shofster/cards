package main

import (
	"cards/common"

	"fyne.io/fyne/v2"
)

/*

  File:    take.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

const TAKE_ID = -1

// p.discardTileContainer.Objects - [1] is spcaer, [0] is the current discard
func (p *player) resetTakes(dim bool, code string) {
	p.history = nil
	for i := range maxHistory {
		ci := p.historyContainer.Objects[i+1].(*common.CardImage)
		ci.Replace(getTileResource(code))
	}
	p.takes = nil
	p.takeContainer.RemoveAll()
	fyne.Do(func() {
		p.historyContainer.Refresh()
		p.takeContainer.Refresh()
	})
}

// only PLAYER can click on aI's discard
func enableTake(b bool, cb func(int)) {
	ti := players[AI].discardTileContainer.Objects[0].(*common.CardImage)
	switch b {
	case true: // show clickable
		if players[AI].discardTile != nil {
			ti.Enabled = true
			ti.SetHighlite(true)
			ti.OnTap = func(ci *common.CardImage) {
				cb(TAKE_ID) // magic number
			}
		}
	case false: // hide clickable
		ti.Enabled = false
		ti.SetHighlite(false)
		ti.OnTap = nil
	}
	fyne.Do(func() {
		players[AI].discardTileContainer.Refresh()
	})
}
