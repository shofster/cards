package main

import (
	"cards/assets/riichi/meld"
	"cards/common"
	"fmt"

	"fyne.io/fyne/v2"
)

/*

  File:    player.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

type player struct {
	id int // PLAYER = 0, AI = 1
	//	icon                 string
	tiles                []*Tile
	tileContainer        *fyne.Container // the current player Tiles
	history              []Tile
	historyContainer     *fyne.Container // the last discards dimmed
	discardTile          *Tile           // the actual discard Time
	discardTileContainer *fyne.Container // the actual discard image and a blank spacer
	takes                []Tile
	takeContainer        *fyne.Container // possible Melds
	hasPair              bool
	melds                []*Meld
	meldContainer        *fyne.Container // the player selected Melds
	seat                 string          // East (North) or South
	winner               bool
}

// dummy tiles to replace Tile images
func NewEmptyTile(id int) *Tile {
	return &Tile{
		Code:  "empty",
		Suit:  "",
		Rank:  0,
		Red:   false,
		ID:    id,
		Dummy: true,
	}
}
func NewBlankTile(id int) *Tile {
	return &Tile{
		Code:  "blank",
		Suit:  "",
		Rank:  0,
		Red:   false,
		ID:    id,
		Dummy: true,
	}
}
func NewReverseTile(id int) *Tile {
	return &Tile{
		Code:  "reverse",
		Suit:  "",
		Rank:  0,
		Red:   false,
		ID:    id,
		Dummy: true,
	}
}

// get an image resource from meld or configured tileMap
func getTileResource(code string) fyne.Resource {
	// check meld and special resources
	res, ok := meld.MeldResource[code]
	if ok {
		return res
	}
	res, ok = tileMap[code]
	if ok {
		return res
	}
	panic(fmt.Sprintf("getTileResource: code %s", code))
}

// initialize a hand and link images to the tiles in the hand
func (p *player) deal() {
	for i := 0; i < len(p.tiles); i++ {
		p.tiles[i].ID = i
	}
	p.melds = nil
	p.discardTile = nil
	p.meldContainer.RemoveAll()
	p.tileContainer.RemoveAll()
	p.hasPair = false
	p.winner = false
	p.resetTakes(true, "empty")
	p.disableDiscard()

	// images for the hand tiles
	p.tileContainer.RemoveAll()
	for i := range len(p.tiles) {
		id := p.tiles[i].ID
		p.tiles[i].Dummy = false
		res := getTileResource(p.tiles[i].Code)
		ci := common.NewCardImage(res, id)
		ci.ID = id
		p.tileContainer.Objects = append(p.tileContainer.Objects, ci)
	}
	p.tileContainer.Refresh()
}

func enablePlayerDiscardSelects(b bool, cb func(int)) {
	p := players[PLAYER]
	defer fyne.Do(func() {
		p.tileContainer.Refresh()
	})
	switch b {
	case true:
		for _, tile := range players[PLAYER].tiles {
			if !tile.Dummy { // ignore blank/empty tiles
				ci := p.tileContainer.Objects[tile.ID].(*common.CardImage)
				ci.Enabled = true
				ci.SetHighlite(true)
				ci.OnTap = func(ci *common.CardImage) {
					cb(ci.ID)
				}
			}
		}
	case false:
		for id := range p.tileContainer.Objects {
			ci := p.tileContainer.Objects[id].(*common.CardImage)
			ci.Enabled = false
			ci.SetHighlite(false)
			ci.OnTap = nil
		}
	}
}

/*   -----------------------------------------  */

func (p *player) addMeldPoints(m *Meld) {
	switch p.id {
	case 0:
		n := common.IntGet(playerScoreData)
		n += m.Score()
		common.IntSet(playerScoreData, n)
	case 1:
		n := common.IntGet(aiScoreData)
		n += m.Score()
		common.IntSet(aiScoreData, n)
	}
}

func (p *player) disableDiscard() {
	p.discardTile = nil
	ci := p.discardTileContainer.Objects[0].(*common.CardImage)
	ci.Replace(getTileResource("empty"))
	ci.SetDim(true)
	fyne.Do(func() {
		p.discardTileContainer.Refresh()
	})
}

func (p *player) takeDiscard() *Tile {
	defer Logit(func() string { return "\n" })
	Logit(func() string { return fmt.Sprintf("taking Discard (%d)", p.id) })
	if p.discardTile != nil && !p.discardTile.Dummy {
		Logit(func() string { return fmt.Sprintf(", took %s", p.discardTile.Code) })
		take := *p.discardTile
		take.ID = -1 // flag as the discard
		return &take
	}
	Logit(func() string { return fmt.Sprintf(", NO take") })
	return nil
}

// discard Tile[id] and set Tile and its image to empty
func (p *player) discard(tile *Tile) {
	var ops = make([]func(), 0)
	defer fyne.Do(func() {
		ops = append(ops, func() { p.discardTileContainer.Refresh() })
		ops = append(ops, func() { p.historyContainer.Refresh() })
		doWait(&ops, false)
	})
	//discard := p.tiles[id] // the Tile to trash
	p.discardTile = tile
	Logit(func() string {
		return fmt.Sprintf(
			"discard (%d) tile %s from slot %d \n", p.id, tile.Code, tile.ID)
	})
	dc := p.discardTileContainer.Objects[0].(*common.CardImage)
	ops = append(ops, func() { dc.SetDim(false) })
	ops = append(ops, func() { dc.Replace(getTileResource(tile.Code)) })
	// append to history , roll image
	// the current discard tile is p.discards[len(p.discards)-1]
	p.history = append(p.history, *tile)
	n := len(p.history)
	// ([0] and [6] are spacers)
	if n <= maxHistory { // does not include spacer
		ci := p.historyContainer.Objects[n].(*common.CardImage)
		ci.ID = tile.ID
		ops = append(ops, func() { ci.Replace(getTileResource(tile.Code)) })
	} else { // shift the images left 1 spot
		Logit(func() string {
			return fmt.Sprintf("discard (%d) n is %d", p.id, n)
		})
		for i := 1; i < maxHistory; i++ {
			dst := p.historyContainer.Objects[i].(*common.CardImage)
			src := p.historyContainer.Objects[i+1].(*common.CardImage)
			dst.ID = src.ID
			dst.Img.Resource = src.Img.Resource
			ops = append(ops, func() { dst.Img.Refresh() })
		}
		// fill the last slot
		Logit(func() string { return fmt.Sprintf(" fill %d\n", maxHistory) })
		ci := p.historyContainer.Objects[maxHistory].(*common.CardImage)
		ci.ID = tile.ID
		ops = append(ops, func() { ci.Replace(getTileResource(tile.Code)) })
	}
	// repalce the hand's Tile and image
	p.replaceHandTile(NewBlankTile(tile.ID))
}

// insert the tile into the discarded slot in the hand
func (p *player) insert(tile *Tile) {
	slot := p.discardTile.ID
	tile.ID = slot
	p.replaceHandTile(tile)
}

// replace the hand Tile and its image with the new Tile
// reset any attributes
func (p *player) replaceHandTile(tile *Tile) {
	if tile.Code != "blank" {
		Logit(func() string {
			return fmt.Sprintf("replaceHandTile (%d) slot %d tile: %s with %s\n",
				p.id, tile.ID, p.tiles[tile.ID].Code,
				tile.Code)
		})
	}
	defer fyne.Do(func() {
		p.tileContainer.Refresh()
	})
	p.tiles[tile.ID] = tile
	ci := p.tileContainer.Objects[tile.ID].(*common.CardImage)
	ci.Replace(getTileResource(tile.Code))
	ci.SetHighlite(false)
	ci.Enabled = false
}

// add a Meld for this hand
// replace the hand tile that aremembers of the meld
func (p *player) newMeld(meld *Meld) {
	defer fyne.Do(func() {
		p.meldContainer.Refresh()
	})

	// add meld to player's list
	p.melds = append(p.melds, meld)
	// add each tile image
	for _, t := range meld.tiles {
		ci := common.NewCardImage(getTileResource(t.Code), 0)
		p.meldContainer.Objects = append(p.meldContainer.Objects, ci)
	}
	// spacer
	blank := getTileResource("blank")
	sep := common.NewCardImage(blank, 0)
	p.meldContainer.Objects = append(p.meldContainer.Objects, sep)

	// blank out the tiles used in this meld
	for _, tile := range meld.tiles {
		if tile.ID != -1 {
			p.replaceHandTile(NewBlankTile(tile.ID))
		}
	}
}

// add a Meld image to the player's images
func (p *player) addTakeMeld(meld Meld) {
	defer fyne.Do(func() {
		p.takeContainer.Refresh()
	})
	ci := common.NewCardImage(getTileResource(meld.kind.String()), 0)
	p.takeContainer.Objects = append(p.takeContainer.Objects, ci)
	if len(p.takeContainer.Objects) > 6 {
		p.takeContainer.Objects = p.takeContainer.Objects[:6]
	}
}

// add a non-selectable Meld image to the hand
// call it a "Ron" type Meld - just for show.
func (p *player) addWinnerMeld() {
	defer fyne.Do(func() {
		p.takeContainer.Refresh()
	})
	ci := common.NewCardImage(getTileResource("ron"), 0)
	p.takeContainer.Objects = append(p.takeContainer.Objects, ci)
	if len(p.takeContainer.Objects) > 6 {
		p.takeContainer.Objects = p.takeContainer.Objects[:6]
	}
}

// determine if a Meld using hand and opponent's discard
func (p *player) detectMeld(discard *Tile) *Meld {
	m := DetectMeld(p.hasPair, p.tiles, discard)
	if m != nil {
		if m.kind == MeldPair {
			p.hasPair = true
		}
		SortTiles(m.tiles)
	}
	return m
}
