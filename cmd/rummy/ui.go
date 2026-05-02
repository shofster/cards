package main

import (
	"fmt"
	"games/common"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
*/

// --- BOARD STRUCTURE ---------------------------------------------------------

var dealerHand *fyne.Container
var dealerMeld binding.String
var dealerPlayButton *widget.Button
var playerHand *fyne.Container
var playerMeld binding.String
var playerKnockButton *widget.Button

var stock *common.CardImage
var discard *common.CardImage
var playCenter *fyne.Container
var scoreCenter *fyne.Container

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

	// Dealer fan (10 cards, face-down)
	dealerLayout := &OverlapLayout{Offset: 40, Vertical: false}
	dealerHand = container.New(dealerLayout, buildHandImages()...)
	dealerHand.Refresh()
	// Dealer meld row (hidden until reveal)
	dealerMeld = binding.NewString()
	dealerLabel := widget.NewLabelWithData(dealerMeld)

	// Player fan (10 cards, face-up)
	playerLayout := &OverlapLayout{Offset: 73, Vertical: false}
	playerHand = container.New(playerLayout, buildHandImages()...)
	playerHand.Refresh()
	// Player meld row (always visible)
	playerMeld = binding.NewString()
	playerLabel := widget.NewLabelWithData(playerMeld)
	_ = playerMeld.Set("Player Meld")

	// PLAY : Center row: stock + discard
	stock = common.NewCardImage(cardMap["reverse"].Resource, 0)
	stockCell := container.NewBorder(nil, widget.NewLabel("STOCK"), nil,
		nil, stock)
	discard = common.NewCardImage(cardMap["joker"].Resource, 1)
	discardCell := container.NewBorder(nil, widget.NewLabel("DISCARD"), nil,
		nil, discard)
	playCenter = container.NewCenter(container.NewHBox(stockCell, discardCell))

	// RESULTS : Center row: stats from hand
	resultGrid := container.NewGridWithRows(6)
	resultGrid.Add(container.NewCenter(widget.NewLabel("-- RESULTS --")))
	winner := widget.NewLabel("")
	winnerData.AddListener(binding.NewDataListener(func() {
		r, _ := winnerData.Get()
		winner.SetText(fmt.Sprintf("Winner is %s", r))
	}))
	resultGrid.Add(winner)
	outcome := widget.NewLabel("")
	outcomeData.AddListener(binding.NewDataListener(func() {
		r, _ := outcomeData.Get()
		outcome.SetText(fmt.Sprintf("Outcome: %s", r))
	}))
	resultGrid.Add(outcome)
	playerPoints := widget.NewLabel("")
	playerData.AddListener(binding.NewDataListener(func() {
		r, _ := playerData.Get()
		playerPoints.SetText(fmt.Sprintf("Player Points: %d", r))
	}))
	resultGrid.Add(playerPoints)
	dealerPoints := widget.NewLabel("")
	dealerData.AddListener(binding.NewDataListener(func() {
		r, _ := dealerData.Get()
		dealerPoints.SetText(fmt.Sprintf("Dealer Points: %d", r))
	}))
	resultGrid.Add(dealerPoints)
	diff := widget.NewLabel("")
	handResult.AddListener(binding.NewDataListener(func() {
		r, _ := handResult.Get()
		diff.SetText(fmt.Sprintf("Result Points: %d", r))
	}))
	resultGrid.Add(diff)

	scoreCenter = container.NewCenter(resultGrid)

	showScoreCenter()
	centerPanel := container.NewStack(playCenter, scoreCenter)

	deal := widget.NewButton("DEAL", func() {
		incHandsPlayed()
		_ = dealerMeld.Set("Dealer Meld")
		dealerHand.Refresh()
		showPlayCenter()
		deal()
	})
	playerKnockButton = widget.NewButton("KNOCK", func() {
		knock()
	})
	playerKnockButton.Disable()

	resultData := binding.NewString()
	result := widget.NewLabelWithData(resultData)
	handResult.AddListener(binding.NewDataListener(func() {
		r, _ := handResult.Get()
		_ = resultData.Set(fmt.Sprintf("Result: %d", r))
	}))

	playedData := binding.NewString()
	played := widget.NewLabelWithData(playedData)
	handsPlayed.AddListener(binding.NewDataListener(func() {
		r, _ := handsPlayed.Get()
		_ = playedData.Set(fmt.Sprintf("Played: %s", r))
	}))

	scoreData := binding.NewString()
	score := widget.NewLabelWithData(scoreData)
	handsScore.AddListener(binding.NewDataListener(func() {
		r, _ := handsScore.Get()
		_ = scoreData.Set(fmt.Sprintf("Total: %s", r))
	}))
	scoreCard := container.NewCenter(container.NewGridWithColumns(6, widget.NewLabel("    "),
		result, layout.NewSpacer(), played, layout.NewSpacer(), score))

	dealerPlayButton = widget.NewButton(" Play ", func() {
		dealerPlay()
	})
	dealerPlayButton.Disable()
	dealerMeldRow := container.NewHBox(dealerPlayButton, layout.NewSpacer(),
		dealerLabel, layout.NewSpacer())

	// Full board layout
	root := container.NewVBox(
		newHorizontal(false),
		dealerMeldRow,
		container.NewCenter(dealerHand),
		newHorizontal(true),
		centerPanel,
		newHorizontal(true),
		container.NewCenter(playerHand),
		container.NewHBox(playerLabel,
			layout.NewSpacer(), playerKnockButton),
		newHorizontal(false),
		container.NewHBox(deal, scoreCard, layout.NewSpacer()),
	)

	return root
}

func buildHandImages() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 10)
	for i := 0; i < 10; i++ {
		ci := common.NewCardImage(cardMap["joker"].Resource, i)
		objs[i] = ci
	}
	return objs
}

// "pretty" lines
func newHorizontal(gray bool) fyne.CanvasObject {
	if gray {
		r := canvas.NewRectangle(theme.Color(theme.ColorNameBackground))
		r.SetMinSize(fyne.NewSize(0, 8))
		return r
	} else {
		r := canvas.NewRectangle(theme.Color(theme.ColorNameForeground))
		r.SetMinSize(fyne.NewSize(0, 10))
		return r
	}

}

func resetHands() {
	for i := 0; i < 10; i++ {
		dealerHand.Objects[i] = common.NewCardImage(cardMap["reverse"].Resource, i)
		playerHand.Objects[i] = common.NewCardImage(cardMap["reverse"].Resource, i)
	}
}

func resetPlayerHand() {
	for i := 0; i < 10; i++ {
		ci := playerHand.Objects[i].(*common.CardImage)
		ci.SetDimmed(false)
		ci.SetHighlite(false)
		ci.Enabled = false
		ci.Refresh()
	}
}

func setPlayerCard(id int, card common.Card) {
	ci := playerHand.Objects[id].(*common.CardImage)
	ci.Img.Resource = card.Resource
	ci.Refresh()
}

func setDealerCard(id int, card common.Card) {
	ci := dealerHand.Objects[id].(*common.CardImage)
	ci.Img.Resource = card.Resource
	ci.Refresh()
}

func showPlayCenter() {
	playCenter.Show()
	scoreCenter.Hide()
}

func showScoreCenter() {
	playCenter.Hide()
	scoreCenter.Show()
}

/*

BOARD LAYOUT

DEALER ROWs
[ Dealer Melds Row ]
[ Dealer CARDS - Hidden until Gin / Knock ]

CENTER ROW
[ Stock ]   [ Discard ]

PLAYER ROWs
[ PlayerCards ]
[ Dealer Melds Row ]

CONTROL & STATUS ROW

*/

// ------------------------------------------------------------------
//
// Layout for a horizontal or vertical "hands" - Card images.
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
