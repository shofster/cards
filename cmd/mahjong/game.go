package main

import (
	"cards/assets/riichi/black"
	"cards/assets/riichi/regular"
	"cards/common"
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
)

/*

  File:    game.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

const (
	Grandma  = "NaiNai"
	Dragon   = "XiaoLong"
	NaiNai   = "♡" //  paternal grandma
	XiaoLong = "🐉" // little dragon
)
const PLAYER = Grandma
const AI = Dragon

const (
	maxHistory = 5
	emptyID    = 997
	blankID    = 998
)

var (
	tileMap        map[string]fyne.Resource
	stock          Tiles
	players        map[string]*player
	winner         string
	stockSlot      *fyne.Container
	duration       time.Duration
	discardMessage = "Choose a Discard"
)

func initGame() {
	duration = time.Millisecond * 1000
	players = make(map[string]*player)
	players[PLAYER] = &player{
		id:      0, // for discard slot
		tiles:   make([]*Tile, 0),
		melds:   make(Melds, 0),
		history: make([]Tile, 0),
	}
	players[AI] = &player{
		id:      1,
		tiles:   make([]*Tile, 0),
		melds:   make(Melds, 0),
		history: make([]Tile, 0),
	}

	switch Ctx.tileTheme {
	case BLACK:
		tileMap = black.BlackResource
	case REGULAR:
		tileMap = regular.RegularResource
	default:
		panic("Invalid TileTheme")
	}

}

func start() {
	common.IntSet(playerScoreData, 0)
	common.IntSet(aiScoreData, 0)
	common.StringSet(turnState, PLAYER)
	si := stockSlot.Objects[0].(*common.CardImage)
	si.OnTap = func(ci *common.CardImage) {
		Logit(func() string { return fmt.Sprintf("stock clicked (%d)\n", ci.ID) })
	}
	deal()
}

func deal() {
	deck := BuildDeck(tileMap)
	deck2 := BuildDeck(tileMap)
	deck3 := BuildDeck(tileMap)
	deck4 := BuildDeck(tileMap)
	tileDeck := append(append(append(deck, deck2...), deck3...), deck4...)
	ShuffleTiles(tileDeck)

	var hands []Tiles // 13 tiles for each player
	hands, stock = DealTiles(tileDeck, 2, 13)
	enableStock(false, nil)

	players[PLAYER].tiles = append(Tiles(nil), hands[0]...)
	players[AI].tiles = append(Tiles(nil), hands[1]...)
	SortTiles(players[PLAYER].tiles)
	SortTiles(players[AI].tiles)

	players[PLAYER].deal()
	players[PLAYER].id = 0
	players[PLAYER].seat = "🀀" // East. use "🀃" for North
	players[AI].id = 1
	players[AI].seat = "🀁" // South.
	players[AI].deal()
	if !Ctx.showAItiles {
		for i := range 13 {
			ci := players[AI].tileContainer.Objects[i].(*common.CardImage)
			ci.Replace(tileMap["back"])
		}
	}
	fyne.Do(func() {
		players[PLAYER].tileContainer.Refresh()
		players[AI].tileContainer.Refresh()
	})
}

// win if 4 melds and 1 pair
func (p *player) checkWin() bool {
	Logit(func() string {
		return fmt.Sprintf("checkWin: (%d) pair %t count %d\n", p.id, p.hasPair, len(p.melds))
	})
	if !p.hasPair || len(p.melds) != 5 {
		return false
	}
	pairs := 0
	melds := 0
	for _, m := range p.melds {
		switch m.kind {
		case MeldPair:
			pairs++
		case MeldKan, MeldPon, MeldChi:
			melds++
		}
	}
	return pairs == 1 && melds == 4
}

func flowGrandma() bool {
	var meld *Meld
	var drawTile *Tile

	initFlow(fmt.Sprintf("%s %s", Dragon, XiaoLong))
	actionButton.OnTapped = func() {
		common.StringSet(turnState, AI)
	}

	var clicked func(id int)
	clicked = func(id int) {
		enableStock(false, nil)
		enableTake(false, nil)
		enablePlayerDiscardSelects(true, clicked)
		// Grandma clicked something
		switch id {
		case -999: // drawTile is already set
			enableStock(false, nil)
			hiliteStock(false) // show
			enableTake(false, nil)
			toast(discardMessage)
		case DRAW_ID: // selected the draw option
			Logit(func() string { return fmt.Sprintf("* clicked: DRAW_ID -  ") })
			toast(discardMessage)
		case TAKE_ID:
			drawTile = players[AI].discardTile
			toast(discardMessage)

		default: // id is the tile to replace
			Logit(func() string { return fmt.Sprintf("clicked (0) on slot %d \n", id) })
			enableStock(false, nil)
			enableTake(false, nil)
			drawTile.ID = id
			players[PLAYER].discard(players[PLAYER].tiles[id])
			players[PLAYER].insert(drawTile)
			enablePlayerDiscardSelects(false, nil)
			fyne.Do(func() {
				actionButton.Enable()
			})
		}
	}

	// PLAYER MUST use take if it makes a meld
	// otherwise choose AI's discard OR draw
	enableStock(false, nil)
	enableTake(false, nil)
	enablePlayerDiscardSelects(false, nil)

	// analyize the AI discard, if one
	if take := players[AI].takeDiscard(); take != nil {
		take.ID = -1
		// if meld: add to melds and remove & blank melded tiles
		if meld = players[PLAYER].detectMeld(take); meld != nil {
			players[PLAYER].addTakeMeld(*meld)
			players[PLAYER].newMeld(meld)
			players[PLAYER].addMeldPoints(meld)
			players[PLAYER].disableDiscard()
			players[AI].disableDiscard()
			stockImage("back") // NO tile taken
			if players[PLAYER].checkWin() {
				players[PLAYER].addWinnerMeld()
				winner = PLAYER
				common.IntSet(playerScoreData, common.IntGet(playerScoreData)+50)
				common.StringSet(turnState, "won")
				return false
			}
			fyne.Do(func() {
				actionButton.Enable()
			})
		} else {
			//  take exists, but NO meld.
			//  choose either draw or take before discard select.
			Logit(func() string { return fmt.Sprintf("* take, NO meld - ") })
			toast("Draw or Take?")
			if drawTile = stockDraw(); drawTile == nil {
				common.StringSet(turnState, "end")
				return false
			}
			// now have the discard and a drawn tile
			enableStock(true, clicked)
			enableTake(true, clicked)
		}
	} else { // normal draw, select a discard
		Logit(func() string { return fmt.Sprintf("* Normal DRAW - ") })
		if drawTile = stockDraw(); drawTile == nil {
			common.StringSet(turnState, "end")
			return false
		}
		//  select a discard to replace it
		clicked(-999) // magic number for this clicked
	}
	return true
}

func flowDragon() bool {

	initFlow(fmt.Sprintf("%s %s", Grandma, NaiNai))
	actionButton.OnTapped = func() {
		common.StringSet(turnState, PLAYER)
	}

	var selectDiscard = func() *Tile {
		var candidates []int
		p := players[AI]
		for _, t := range p.tiles {
			if !t.Dummy { // ignore blank/empty tiles
				candidates = append(candidates, t.ID)
			}
		}
		if len(candidates) == 0 {
			return nil
		}
		id := rand.Intn(len(candidates))
		Logit(func() string {
			return fmt.Sprintf("AI selectDiscard (%d) tile %s %d\n", p.id,
				p.tiles[candidates[id]].Code, p.tiles[candidates[id]].ID)
		})
		return p.tiles[id]

	}

	var dragonDrawDiscard = func() bool {
		if tile := stockDraw(); tile != nil {
			id := selectDiscard()
			players[AI].discard(id)
			players[AI].insert(tile)
			return true
		}
		return false
	}

	if take := players[PLAYER].takeDiscard(); take != nil {
		// if meld: add to melds and remove & blank melded tiles
		take.ID = -1 // mark the taken Tile
		if meld := players[AI].detectMeld(take); meld != nil {
			//toast(fmt.Sprintf("Meld: %s", meld.kind.String()))
			players[AI].addTakeMeld(*meld)
			players[AI].newMeld(meld)
			players[AI].addMeldPoints(meld)
			players[AI].disableDiscard()
			players[PLAYER].disableDiscard()
			stockImage("back") // NO tile taken
			if players[AI].checkWin() {
				players[AI].addWinnerMeld()
				winner = AI
				common.IntSet(aiScoreData, common.IntGet(aiScoreData)+50)
				common.StringSet(turnState, "won")
				return false
			}
		} else {
			if !dragonDrawDiscard() {
				common.StringSet(turnState, "end")
				return false
			}
		}
	} else {
		if !dragonDrawDiscard() {
			common.StringSet(turnState, "end")
			return false
		}
	}

	fyne.Do(func() {
		actionButton.Enable()
	})
	return true
}

func toast(s string) {
	toastWithDuration(s, duration)
}
func toastWithDuration(s string, d time.Duration) {
	common.Toast(Ctx.mainWindow.Canvas(), fyne.NewPos(5, 5), s, common.InfoColor, d)
}

func initFlow(who string) {
	var ops = make([]func(), 0)
	ops = append(ops, func() {
		actionButton.SetText(fmt.Sprintf("Click for: %s", who))
	})
	ops = append(ops, func() {
		actionButton.Disable()
	})
	doWait(&ops, true)

}

// ai's random choice
func (p *player) selectDiscard() *Tile {
	var candidates []int
	for _, t := range p.tiles {
		if !t.Dummy { // ignore blank/empty tiles
			candidates = append(candidates, t.ID)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	id := rand.Intn(len(candidates))
	Logit(func() string {
		return fmt.Sprintf("AI selectDiscard (%d) tile %s %d\n", p.id,
			p.tiles[candidates[id]].Code, p.tiles[candidates[id]].ID)
	})
	return p.tiles[id]
}
