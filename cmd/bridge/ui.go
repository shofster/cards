package main

import (
	"cards/common"
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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

// Container for the 4 hands
var northStack *fyne.Container
var southStack *fyne.Container
var eastStack *fyne.Container
var westStack *fyne.Container

// Container for the 3 Cards of currrent trick
var southTrick *canvas.Image
var westTrick *canvas.Image
var northTrick *canvas.Image
var eastTrick *canvas.Image

// play control buttons
var dealButton *widget.Button
var playButton *widget.Button

// player help -- next action
var nextAction = binding.NewString()
var dealAction = "Press Deal to Begin"
var contractAction = "Select a Contract"
var playAction = "Press Play for South Lead"

// buttons to commit the selection ("played") of the "laid" Card
var southCommit *widget.Button
var westCommit *widget.Button
var northCommit *widget.Button
var eastCommit *widget.Button

var scoreCard *common.ScoreCard

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

// build the widgets for the actual game.
// initGame must previously be called.
func initGameUI() *fyne.Container {

	// the actual image containers for each Player's hand
	buildHandsImages()

	// Variables for a new hand to be dealt and played
	levelSelect := widget.NewSelectWithData(contractLevels, contractLevel)
	suitSelect := widget.NewSelectWithData(contractSuits, contractSuit)
	contractLevel.AddListener(binding.NewDataListener(setContract))
	contractSuit.AddListener(binding.NewDataListener(setContract))
	levelSelect.Disable()
	suitSelect.Disable()

	// clicking Deal enables contract selection
	dealButton = widget.NewButton("Deal", func() {
		resetPlayedCards()
		deal()
		levelSelect.SetSelected(common.StringSet(contractLevel, ""))
		suitSelect.SetSelected(common.StringSet(contractSuit, ""))
		levelSelect.Enable()
		suitSelect.Enable()
		playButton.Disable()
		common.StringSet(nextAction, contractAction)
	})
	// play is enabled when a contract is selected
	playButton = widget.NewButton("Play", nil)
	playButton.Disable()
	playButton.OnTapped = func() {
		playButton.Disable()
		levelSelect.Disable()
		suitSelect.Disable()
		play()
	}

	controls := container.NewCenter(container.NewHBox(widget.NewLabel("     "),
		dealButton, widget.NewLabel("     "),
		widget.NewLabel("Level:"), levelSelect, widget.NewLabel("     "),
		widget.NewLabel("Suit:"), suitSelect, widget.NewLabel("     "),
		playButton, widget.NewLabel("     ")))

	//  Commit buttons to press after trick card is played
	//playerCommit := widget.NewLabel("Player -> ")
	southCommit = widget.NewButton(playerToPosition(South), nil)
	southCommit.Disable()
	westCommit = widget.NewButton(playerToPosition(West), nil)
	westCommit.Disable()
	northCommit = widget.NewButton(playerToPosition(North), nil)
	northCommit.Disable()
	eastCommit = widget.NewButton(playerToPosition(East), nil)
	eastCommit.Disable()

	na := widget.NewLabelWithData(nextAction)
	na.TextStyle = fyne.TextStyle{Italic: true}
	next := container.NewCenter(na)
	controlBox := container.NewVBox(next, newHorizontal(),
		newHorizontal(), controls, newHorizontal())
	common.StringSet(nextAction, dealAction)

	// build the center area where tricks are displayed
	southTrick = canvas.NewImageFromResource(cardMap["joker"].Resource)
	southTrick.FillMode = canvas.ImageFillContain
	southTrickCell := container.NewBorder(nil, southCommit, nil,
		nil, southTrick)
	westTrick = canvas.NewImageFromResource(cardMap["joker"].Resource)
	westTrick.FillMode = canvas.ImageFillContain
	westTrickCell := container.NewBorder(nil, westCommit, nil,
		nil, westTrick)
	northTrick = canvas.NewImageFromResource(cardMap["joker"].Resource)
	northTrick.FillMode = canvas.ImageFillContain
	northTrickCell := container.NewBorder(nil, northCommit, nil,
		nil, northTrick)
	eastTrick = canvas.NewImageFromResource(cardMap["joker"].Resource)
	eastTrick.FillMode = canvas.ImageFillContain
	eastTrickCell := container.NewBorder(nil, eastCommit, nil,
		nil, eastTrick)
	tricks := container.NewGridWithColumns(4,
		southTrickCell, westTrickCell, northTrickCell,
		eastTrickCell)

	// current and previous game results
	//setTrickResult(-1, South)
	tricksBox := container.NewBorder(container.NewCenter(buildTrickResults()),
		nil, nil, nil, tricks)

	center := container.NewBorder(nil, controlBox, nil, nil, tricksBox)

	playArea := container.NewBorder(container.NewCenter(northStack), container.NewCenter(southStack),
		westStack, eastStack, center)

	scoreArea := buildScoreArea()

	view := container.NewBorder(nil, nil, nil, scoreArea, playArea)

	return container.NewBorder(newHorizontal(), newHorizontal(), newVertical(),
		newVertical(), view)
}

// "pretty" lines
func newHorizontal() fyne.CanvasObject {
	r := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	r.SetMinSize(fyne.NewSize(0, 5))
	return r
}
func newVertical() fyne.CanvasObject {
	r := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	r.SetMinSize(fyne.NewSize(5, 0))
	return r
}

func setContract() {
	l := common.StringGet(contractLevel)
	s := common.StringGet(contractSuit)
	if l != "" && s != "" {
		playButton.Enable()
		common.StringSet(nextAction, playAction)
		common.StringSet(contract, fmt.Sprintf("%s %s", l, s))
	}
}

// build all the hands images
func buildHandsImages() {
	northStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.offset),
		Vertical: false},
		buildHandImages()...)
	southStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.offset),
		Vertical: false},
		buildHandImages()...)
	eastStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.offset),
		Vertical: true},
		buildHandImages()...)
	westStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.offset),
		Vertical: true},
		buildHandImages()...)
}

func buildHandImages() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 13)
	for i := 0; i < 13; i++ {
		c := common.NewCardImage(cardMap["reverse"].Resource, i)
		objs[i] = c
		c.Refresh()
	}
	return objs
}

func buildTrickResults() *fyne.Container {
	trickStates = common.NewCardStrip(13)
	b := container.NewHBox(widget.NewLabel("Trick Winner:"))

	for s := 0; s < 13; s++ {
		b.Add(trickStates.Labels[s])
	}

	return b
}

func buildScoreArea() *fyne.Container {

	scoreBox := container.NewVBox()

	contractBid := widget.NewLabelWithData(contract)
	scoreBox.Add(container.NewCenter(
		container.NewHBox(widget.NewLabel("Contract:"), contractBid)))

	scoreBox.Add(container.NewCenter(
		widget.NewLabel("Cards Played:")))

	keys := []string{"#", "S", "W", "N", "E"}
	scoreCard = common.NewScoreCard(keys)
	scoreBox.Add(scoreCard.Grid)

	scoreBox.Add(newHorizontal())

	scoreBox.Add(widget.NewLabel("Deal Result:"))
	dealStatus := widget.NewLabel("")
	dealResult.AddListener(binding.NewDataListener(func() {
		dealStatus.SetText(common.StringGet(dealResult))
	}))
	scoreBox.Add(dealStatus)

	scoreBox.Add(newHorizontal())
	above := widget.NewLabel("")
	scoreBox.Add(container.NewHBox(widget.NewLabel("Above:"), above))
	aboveLine.AddListener(binding.NewDataListener(func() {
		above.SetText(strconv.Itoa((common.IntGet(aboveLine))))
	}))
	below := widget.NewLabel("")
	scoreBox.Add(container.NewHBox(widget.NewLabel("Below:"), below))
	belowLine.AddListener(binding.NewDataListener(func() {
		below.SetText(strconv.Itoa((common.IntGet(belowLine))))
	}))
	scoreBox.Add(newHorizontal())

	historyBid := widget.NewLabel("0")
	scoreBox.Add(container.NewHBox(widget.NewLabel("Contracts Bid:"), historyBid))
	contractsBid.AddListener(binding.NewDataListener(func() {
		historyBid.SetText(strconv.Itoa(common.IntGet(contractsBid)))
	}))

	historyWon := widget.NewLabel("0")
	scoreBox.Add(container.NewHBox(widget.NewLabel("Contracts Made:"), historyWon))
	contractsWon.AddListener(binding.NewDataListener(func() {
		historyWon.SetText(strconv.Itoa(common.IntGet(contractsWon)))
	}))

	// force labels to have previous saved values
	common.IntSet(contractsBid, common.IntGet(contractsBid))
	common.IntSet(contractsWon, common.IntGet(contractsWon))

	return container.NewCenter(scoreBox)
}

// set up the played cards widgets for a "new trick"
func resetPlayedCards() {
	for _, w := range []*canvas.Image{southTrick, westTrick, northTrick, eastTrick} {
		w.Resource = cardMap["reverse"].Resource
		fyne.Do(func() {
			w.Refresh()
		})
	}
}

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
