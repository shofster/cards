package main

import (
	"cards/assets/riichi/black"
	"cards/assets/riichi/meld"
	"cards/assets/riichi/regular"
	"cards/common"
	"fmt"
	"image/color"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    ui.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
+------------------------------------------------------+
|                                                      |
|                       Grandma (TOP)                  |
|                                                      |
|   Grandma's Hand (13 tiles, upright)                 |
|                                                      |
|   Grandma's Melds (open sets, grouped)               |
|      [🀇🀇🀇]   [🀙🀚🀛]   [🀐🀐🀐🀐]                        |
|                                                      |
|   Grandma's Discards (face-up row)    (BOTTOM)       |
|                                                      |
+------------------------------------------------------+
|                                                      |
|   AI's Discards (face-up row)       (TOP)            |
|                                                      |
|   AI's Melds (open sets, grouped)                    |
|      [🀝🀝🀝]   [🀗🀘🀙]                                  |
|                                                      |
|   AI's Hand (13 tiles, hidden or backs)              |
|                                                      |
|                         AI        (BOTTOM)           |
|                                                      |
+------------------------------------------------------+
|                                                      |
|                 [ Face‑Down Stock ]                  |
|                     (shared pile)                    |
|                                                      |
|------------------------------------------------------|
*/

var actionButton *widget.Button

var windSlot *fyne.Container

var playerScoreData = binding.NewInt()
var aiScoreData = binding.NewInt()
var turnState = binding.NewString()

func initUI() *fyne.Container {
	// necessary game variables
	initGame()
	return container.NewStack(makeGreenPanel(initGameUI()))
}

func makeGreenPanel(content *fyne.Container) *fyne.Container {
	// Soft felt green — color.NRGBA{0, 160, 0, 255}
	// Bright emerald — color.NRGBA{20, 180, 20, 255}
	// Muted mint — color.NRGBA{60, 200, 120, 255}
	// Casino - color.NRGBA{0, 128, 0, 255}
	bg := canvas.NewRectangle(color.NRGBA{20, 180, 20, 50})
	return container.NewStack(bg, content)
}

func initGameUI() *fyne.Container {
	var space = Ctx.width + 6

	playerTitle := widget.NewLabel("🀀" + " - " + PLAYER)
	playerTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	playerHandLayout := &OverlapLayout{Offset: space, Vertical: false}
	players[PLAYER].tileContainer = container.New(playerHandLayout, buildHandImages(false)...)
	players[PLAYER].meldContainer = container.NewHBox(buildMeldExamples()...)
	playerLayout := container.NewGridWithRows(2,
		players[PLAYER].tileContainer,
		container.NewHBox(layout.NewSpacer(), players[PLAYER].meldContainer),
	)

	playerHistoryLayout := &OverlapLayout{Offset: space, Vertical: false}
	players[PLAYER].historyContainer = container.New(playerHistoryLayout, buildHistoryImages(false)...)

	playerDiscardLayout := &OverlapLayout{Offset: space, Vertical: false}
	players[PLAYER].discardTileContainer = container.New(playerDiscardLayout, buildDiscardImages()...)

	players[PLAYER].takeContainer = container.NewHBox(fakeDiscardMeldImages()...)

	playerPanel := container.NewBorder(
		container.NewCenter(playerTitle),
		container.NewHBox(
			players[PLAYER].historyContainer,
			layout.NewSpacer(),
			players[PLAYER].takeContainer,
			layout.NewSpacer(),
			players[PLAYER].discardTileContainer,
		),
		nil,
		nil,
		playerLayout)

	ppl := widget.NewLabel("")
	playerScoreData.AddListener(binding.NewDataListener(func() {
		ppl.SetText(fmt.Sprintf("%-6d", common.IntGet(playerScoreData)))
	}))
	ps := container.NewHBox(widget.NewLabel(NaiNai), ppl)

	//  AI - same rows, but vertically mirrored

	aiTitle := widget.NewLabel("🀁" + " - " + AI)
	aiTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	aiHandLayout := &OverlapLayout{Offset: space, Vertical: false}
	players[AI].tileContainer = container.New(aiHandLayout, buildHandImages(false)...)
	players[AI].meldContainer = container.NewHBox(buildMeldExamples()...)
	aiLayout := container.NewGridWithRows(2,
		container.NewHBox(layout.NewSpacer(), players[AI].meldContainer),
		players[AI].tileContainer,
	)

	aiHistoryLayout := &OverlapLayout{Offset: space, Vertical: false}
	players[AI].historyContainer = container.New(aiHistoryLayout, buildHistoryImages(true)...)

	aiDiscardLayout := &OverlapLayout{Offset: space, Vertical: false}
	players[AI].discardTileContainer = container.New(aiDiscardLayout, buildDiscardImages()...)

	players[AI].takeContainer = container.NewHBox(fakeDiscardMeldImages()...)

	aiPanel := container.NewBorder(
		container.NewHBox(
			players[AI].historyContainer,
			layout.NewSpacer(),
			players[AI].takeContainer,
			layout.NewSpacer(),
			players[AI].discardTileContainer,
		),
		container.NewCenter(aiTitle),
		nil,
		nil,
		aiLayout)

	// dynamic bound variables
	apl := widget.NewLabel("")
	aiScoreData.AddListener(binding.NewDataListener(func() {
		apl.SetText(fmt.Sprintf("%-6d", common.IntGet(aiScoreData)))
	}))
	pl := container.NewHBox(widget.NewLabel("  "+XiaoLong), apl)

	tl := container.NewHBox(widget.NewLabel("  Turn is: "),
		widget.NewLabelWithData(turnState))

	common.IntSet(playerScoreData, 0)
	common.IntSet(aiScoreData, 0)

	actionButton = widget.NewButtonWithIcon("DEAL ", theme.MediaPlayIcon(),
		func() {
			actionButton.SetText("")
			actionButton.Disable()
			start()
		})

	stockImage := common.NewCardImage(tileMap["back"], 0)
	stockSlot = container.NewStack(stockImage)
	wind := canvas.NewImageFromResource(tileMap["🀀"])
	wind.FillMode = canvas.ImageFillContain
	wind.SetMinSize(fyne.NewSize(Ctx.deckWidth, Ctx.deckHeight))
	windSlot = container.NewStack(wind)
	spacer := common.NewCardImage(meld.MeldResource["blank"], 0)

	control := container.NewHBox(container.NewHBox(ps, pl),
		layout.NewSpacer(), container.NewHBox(windSlot, tl, spacer,
			actionButton),
		layout.NewSpacer(), widget.NewLabel("    Stock: "), stockSlot, layout.NewSpacer())
	scroll := container.NewScroll(
		container.NewGridWithRows(2, playerPanel, aiPanel))
	scroll.Direction = fyne.ScrollHorizontalOnly

	board := container.NewBorder(nil, nil, nil, nil, scroll)

	var showScores = func() string {
		return fmt.Sprintf("%s - %d\n%s - %d",
			PLAYER, common.IntGet(playerScoreData),
			AI, common.IntGet(aiScoreData))
	}

	// PLAYER / AI turn toggle
	turnState.AddListener(binding.NewDataListener(func() {

		Logit(func() string {
			return fmt.Sprintf("\n  ----- turnState - running %s\n", common.StringGet(turnState))
		})
		switch common.StringGet(turnState) {
		case PLAYER:
			go flowGrandma()
		case AI:
			go flowDragon()
		case "won":
			fyne.Do(func() {
				actionButton.SetText("GAME OVER")
				actionButton.Disable()
			})
			dialog.ShowInformation("WINNER",
				fmt.Sprintf("%s\n%s\n", winner, showScores()), Ctx.mainWindow)
		case "end":
			fyne.Do(func() {
				actionButton.SetText("GAME OVER")
				actionButton.Disable()
			})
			dialog.ShowInformation("STOCK EXHAUSTED", showScores(), Ctx.mainWindow)
		}

	}))

	r := canvas.NewRectangle(theme.Color(theme.ColorNameForeground))
	r.SetMinSize(fyne.NewSize(0, 4))
	line := container.NewStack(r)
	cb := container.NewVBox(line, control, line)
	return container.NewBorder(cb,
		nil, nil, nil, board)
}

func buildHandImages(ai bool) []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 13)
	var tileMap map[string]fyne.Resource
	if ai {
		tileMap = black.BlackResource
	} else {
		tileMap = regular.RegularResource
	}
	keys := make([]string, 0, len(tileMap))
	for k := range tileMap {
		suit, _, _ := parseTile(k)
		if suit != "?" { // ignore not dealt tiles
			keys = append(keys, k)
		}
	}
	for i := range 13 {
		ci := common.NewCardImage(tileMap[keys[rand.Intn(len(keys))]], i)
		objs[i] = ci
	}
	return objs
}

func buildHistoryImages(ai bool) []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0)
	// images are at [1] to allow the empty slot
	placeholder := common.NewCardImage(getTileResource("blank"), 0)
	objs = append(objs, placeholder)
	for i := range maxHistory {
		res := regular.RegularResource["🀣"]
		if ai {
			res = black.BlackResource["🀆"]
		}
		ci := common.NewCardImage(res, i)
		// disable and dim
		ci.SetDimmed(true)
		objs = append(objs, ci)
	}
	return objs
}

func buildDiscardImages() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0)
	discard := meld.MeldResource["empty"]
	empty := common.NewCardImage(discard, emptyID)
	blank := meld.MeldResource["blank"]
	objs = append(objs, empty)
	objs = append(objs, common.NewCardImage(blank, blankID))
	return objs
}

func buildMeldExamples() []fyne.CanvasObject {
	return buildMeldImages(fakeMelds())
}

func fakeDiscardMeldImages() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0)
	objs = append(objs, common.NewCardImage(meld.MeldResource["pair"], 0))
	objs = append(objs, common.NewCardImage(meld.MeldResource["kan"], 1))
	objs = append(objs, common.NewCardImage(meld.MeldResource["chi"], 2))
	objs = append(objs, common.NewCardImage(meld.MeldResource["ron"], 3))
	return objs
}

func fakeMelds() []*Meld {
	var pair = Tiles{NewTile("🀇"), NewTile("🀇")}
	var pon = Tiles{NewTile("🀙"), NewTile("🀙"), NewTile("🀙")}
	chi := Tiles{NewTile("🀄"), NewTile("🀅"), NewTile("🀆")}
	var kan = Tiles{NewTile("🀖"), NewTile("🀖"), NewTile("🀖"), NewTile("🀖")}
	chi2 := Tiles{NewTile("🀄"), NewTile("🀅"), NewTile("🀆")}
	return []*Meld{
		NewMeld(MeldPair, true, pair),
		NewMeld(MeldPon, true, pon),
		NewMeld(MeldKan, true, kan),
		NewMeld(MeldChi, true, chi),
		NewMeld(MeldChi, true, chi2),
	}
}

func buildMeldImages(melds []*Meld) []fyne.CanvasObject {

	common.DefaultCardImageWidth = Ctx.meldWidth
	common.DefaultCardImageHeight = Ctx.meldHeight

	objs := make([]fyne.CanvasObject, 0)
	n := 0

	for _, m := range melds {
		// the meld icon separator
		var icon fyne.Resource
		switch m.kind {
		case MeldPon:
			icon = meld.MeldResource["pon"]
		case MeldKan:
			icon = meld.MeldResource["kan"]
		case MeldChi:
			icon = meld.MeldResource["chi"]
		}
		objs = append(objs, common.NewCardImage(icon, n))
		n++
		tiles := m.tiles
		switch m.kind {
		case MeldPon:
			// 3 tiles
			for i := range 3 {
				res := tileMap[tiles[i].Code]
				objs = append(objs, common.NewCardImage(res, n))
				n++
			}
		case MeldKan:
			// 4 tiles
			for i := range 4 {
				res := tileMap[tiles[i].Code]
				objs = append(objs, common.NewCardImage(res, n))
				n++
			}
		case MeldChi:
			// 3 sequential tiles
			for i := range 3 {
				res := tileMap[tiles[i].Code]
				objs = append(objs, common.NewCardImage(res, n))
				n++
			}
		}
	}

	for _, m := range melds {
		// the meld icon separator
		var icon fyne.Resource
		switch m.kind {
		case MeldPair:
			icon = meld.MeldResource["pair"]
		default:
			continue
		}
		b := meld.MeldResource["blank"]
		objs = append(objs, common.NewCardImage(b, n))
		n++
		objs = append(objs, common.NewCardImage(icon, n))
		n++
		tiles := m.tiles
		switch m.kind {
		case MeldPair:
			// 2 tiles
			for i := range 2 {
				res := tileMap[tiles[i].Code]
				objs = append(objs, common.NewCardImage(res, n))
				n++
			}
		}
	}

	return objs
}

// ------------------------------------------------------------------
//
// Layout for a horizontal or vertical "hands" - Tile images.
//
// ------------------------------------------------------------------

type OverlapLayout struct {
	Offset   float32
	Vertical bool // true = top down, false = left right
}

func (l *OverlapLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n == 0 {
		return
	}

	// East/West players (vertical hands)
	maxVerticalCardHeight := Ctx.height
	maxVerticalCardWidth := Ctx.width

	// North/South players (horizontal hands)
	maxHorizontalCardHeight := Ctx.height
	maxHorizontalCardWidth := Ctx.width

	var cardWidth, cardHeight float32

	if l.Vertical {
		// Vertical -- Height is the limiting dimension
		usableHeight := size.Height - float32(n-1)*l.Offset
		if usableHeight < 40 {
			usableHeight = 40
		}

		cardHeight = usableHeight
		cardWidth = cardHeight * (5.0 / 7.0)

		// Clamp to vertical maximums
		if cardHeight > maxVerticalCardHeight {
			cardHeight = maxVerticalCardHeight
			cardWidth = cardHeight * (5.0 / 7.0)
		}
		if cardWidth > maxVerticalCardWidth {
			cardWidth = maxVerticalCardWidth
			cardHeight = cardWidth * (7.0 / 5.0)
		}

	} else {
		// horizontal -- Width is the limiting dimension
		usableWidth := size.Width - float32(n-1)*l.Offset
		if usableWidth < 40 {
			usableWidth = 40
		}

		cardWidth = usableWidth
		cardHeight = cardWidth * (7.0 / 5.0)

		// Clamp to horizontal maximums
		if cardWidth > maxHorizontalCardWidth {
			cardWidth = maxHorizontalCardWidth
			cardHeight = cardWidth * (7.0 / 5.0)
		}
		if cardHeight > maxHorizontalCardHeight {
			cardHeight = maxHorizontalCardHeight
			cardWidth = cardHeight * (5.0 / 7.0)
		}
	}

	// anchor positions
	var x, y float32

	// Centering --
	if l.Vertical {
		// Center horizontally
		x = (size.Width - cardWidth) / 2
		y = 0
	} else {
		// Center vertically
		x = 0
		y = (size.Height - cardHeight) / 2
	}

	// Card size and positions.
	for i, c := range objects {
		c.Resize(fyne.NewSize(cardWidth, cardHeight))

		if l.Vertical {
			c.Move(fyne.NewPos(x, float32(i)*l.Offset))
		} else {
			c.Move(fyne.NewPos(float32(i)*l.Offset, y))
		}
	}
}

func (l *OverlapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	n := len(objects)
	if n == 0 {
		return fyne.NewSize(0, 0)
	}

	// East/West (vertical hands)
	maxVerticalCardHeight := Ctx.height
	maxVerticalCardWidth := Ctx.width

	// North/South (horizontal hands)
	maxHorizontalCardHeight := Ctx.height
	maxHorizontalCardWidth := Ctx.width

	if l.Vertical {
		// East/West
		h := maxVerticalCardHeight + float32(n-1)*l.Offset
		return fyne.NewSize(maxVerticalCardWidth, h)
	}

	// North/South
	w := maxHorizontalCardWidth + float32(n-1)*l.Offset
	return fyne.NewSize(w, maxHorizontalCardHeight)
}
