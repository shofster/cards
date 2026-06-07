package main

import (
	"cards/common"
	"fmt"
	"image/color"
	"time"

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

var (
	spiral        *SpiralTracks
	runSpiral     bool
	aiStack       *fyne.Container
	aiIcon        *widget.Icon
	playerStack   *fyne.Container
	playerIcon    *widget.Icon
	starterStack  *fyne.Container
	cribStack     *fyne.Container
	dealButton    *widget.Button
	discardButton *widget.Button
	peggingButton *widget.Button
	scoreButton   *widget.Button
)

func initUI() *fyne.Container {
	return container.NewStack(makeGreenPanel(initGameUI()))
}

func makeGreenPanel(content *fyne.Container) *fyne.Container {
	// Soft felt green — color.NRGBA{0, 160, 0, 255}
	// Bright emerald — color.NRGBA{20, 180, 20, 255}
	// Muted mint — color.NRGBA{60, 200, 120, 255}
	// Casino - color.NRGBA{0, 128, 0, 255}
	bg := canvas.NewRectangle(color.NRGBA{20, 180, 20, 50})
	return container.NewStack(bg, content)
	//return content
}

// build the widgets for the actual game.
// initGame must previously be called.
func initGameUI() *fyne.Container {

	spiral = NewSpiralTracks(Ctx.pegSize, 121, Ctx.trackSize)
	playerScore := widget.NewLabel("")
	Ctx.playerScore.AddListener(binding.NewDataListener(func() {
		playerScore.SetText(fmt.Sprintf("%s : %4d",
			Ctx.playerName, common.IntGet(Ctx.playerScore)))
	}))
	aiScore := widget.NewLabel("")
	Ctx.aiScore.AddListener(binding.NewDataListener(func() {
		aiScore.SetText(fmt.Sprintf("%s : %4d",
			Ctx.aiName, common.IntGet(Ctx.aiScore)))
	}))

	dealButton = widget.NewButton("DEAL", func() {
		deal()
	})
	dealer := widget.NewLabel("")
	Ctx.dealerName.AddListener(binding.NewDataListener(func() {
		dealer.SetText(fmt.Sprintf("(%s)",
			common.StringGet(Ctx.dealerName)))
	}))

	scores := container.NewHBox()
	scores.Add(dealButton)
	scores.Add(dealer)
	scores.Add(widget.NewLabel("  Scores - "))
	scores.Add(playerScore)
	scores.Add(aiScore)
	scoreBoard := container.NewBorder(nil, scores, nil, nil, spiral.c)

	// PLAYER row - 6 cards, 2 hidden
	playerStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.cardOffset) * 1.8,
		Vertical: false},
		buildHandImages(6)...)
	playerRow := container.NewStack(container.NewCenter(playerStack))

	// AI row - 6 cards hidden
	aiStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.cardOffset) * 1.8,
		Vertical: false},
		buildHandImages(6)...)
	aiRow := container.NewStack(container.NewCenter(aiStack))

	// CRIB row - starter, 4 card crib (discards), control buttons
	starterStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.cardOffset),
		Vertical: false},
		buildHandImages(1)...)
	pileCards := container.NewStack(starterStack)
	cribStack = container.New(&OverlapLayout{
		Offset:   float32(Ctx.cardOffset),
		Vertical: false},
		buildHandImages(4)...)
	cribCards := container.NewStack(cribStack)
	pc := container.NewHBox(pileCards, cribCards)

	discardButton = widget.NewButtonWithIcon("Discard", theme.DeleteIcon(), func() {
		discardAction(PLAYER)
	})
	discardButton.Disable()
	peggingButton = widget.NewButtonWithIcon("Pegging", theme.DeleteIcon(), func() {
		peggingAction(PLAYER)
	})
	peggingButton.Disable()
	scoreButton = widget.NewButton("Score", func() {
		scoreAction()
	})
	control := container.NewVBox(
		discardButton,
		peggingButton,
		scoreButton,
	)
	scoreButton.Disable()

	// cards pegged in order (and counted) and hand count
	lanes = make([]*common.CardStrip, 2)
	lanes[AI] = common.NewCardStrip(8)
	lanes[PLAYER] = common.NewCardStrip(8)
	ar := container.NewHBox()
	pr := container.NewHBox()
	for i := 0; i < 8; i++ {
		ar.Add(lanes[AI].Labels[i])
		pr.Add(lanes[PLAYER].Labels[i])
	}

	cf := "Current Count: %4d"
	handCount := widget.NewLabel("")
	Ctx.handPoints.AddListener(binding.NewDataListener(func() {
		handCount.SetText(fmt.Sprintf(cf, common.IntGet(Ctx.handPoints)))
	}))
	played := container.NewVBox(ar, container.NewCenter(handCount), pr)

	cribRow := container.NewBorder(nil, played, nil, control, pc)

	playBoard := container.NewVBox()
	playBoard.Add(aiRow)
	playBoard.Add(container.NewCenter(cribRow))
	playBoard.Add(playerRow)
	cardPlayed(-1, ".x.")

	demo()
	return container.NewGridWithColumns(2, scoreBoard, playBoard)

}

func cardPlayed(card int, icon string) {
	if card == -1 {
		for i := range 8 {
			lanes[PLAYER].SetCard(i, fmt.Sprintf("%3s", icon))
			lanes[AI].SetCard(i, fmt.Sprintf("%3s", icon))
		}
		return
	}
}

func buildHandImages(n int) []fyne.CanvasObject {
	common.DefaultCardImageWidth = Ctx.cardSizeWidth
	common.DefaultCardImageHeight = Ctx.cardSizeHeight
	objs := make([]fyne.CanvasObject, 0)
	for i := 0; i < n; i++ {
		c := common.NewCardImage(Ctx.cardMap["reverse"].Resource, i)
		objs = append(objs, c)
		c.Refresh()
	}
	return objs
}

// manage player discard button
func discardButtonState(state int) {
	switch state {
	case 0:
		discardButton.SetText("Discards?")
		discardButton.Enable()
	case 1:
		discardButton.SetText("Discard")
		discardButton.Enable()
	case 2:
		discardButton.Disable()
		fyne.Do(func() {
			time.Sleep(time.Millisecond * 300)
			discardButton.Enable()
		})
	default:
		discardButton.SetText("Discard")
		discardButton.Disable()
	}
}

// manage player pegging button
func peggingButtonState(state int) {
	switch state {
	case 0:
		peggingButton.SetText("Play a Card")
		peggingButton.Enable()
	case 1:
		peggingButton.SetText("GO")
		peggingButton.Enable()
	case 2:
		peggingButton.Disable()
		fyne.Do(func() {
			time.Sleep(time.Millisecond * 300)
			peggingButton.Enable()
		})
	default:
		peggingButton.SetText("Pegging")
		peggingButton.Disable()
	}
}

/*
playAnimation([]fyne.Resource{
    resourceFrame0Png,
    resourceFrame1Png,
    resourceFrame2Png,
    resourceFrame3Png,
}, myImageWidget)

*/

func playAnimation(imgs []fyne.Resource, target *canvas.Image) {
	go func() {
		for {
			for _, r := range imgs {
				target.Resource = r
				target.Refresh()
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

func demo() {
	redPos := 0
	bluePos := 0
	runSpiral = true
	//spiral.Delay = 500
	go func() {
		for {
			if !runSpiral {
				spiral.reset()
				return
			}
			time.Sleep(time.Second)
			spiral.advanceA(redPos)
			spiral.advanceB(bluePos)
			redPos += 10
			bluePos += 8
			if redPos > 120 {
				spiral.reset()
				redPos = 0
				bluePos = 0
			}
		}
	}()
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
	maxVerticalCardHeight := Ctx.cardSizeHeight
	maxVerticalCardWidth := Ctx.cardSizeWidth

	// North/South players (horizontal hands)
	maxHorizontalCardHeight := Ctx.cardSizeHeight
	maxHorizontalCardWidth := Ctx.cardSizeWidth

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
	maxVerticalCardHeight := Ctx.cardSizeHeight
	maxVerticalCardWidth := Ctx.cardSizeWidth

	// North/South (horizontal hands)
	maxHorizontalCardHeight := Ctx.cardSizeHeight
	maxHorizontalCardWidth := Ctx.cardSizeWidth

	if l.Vertical {
		// East/West
		h := maxVerticalCardHeight + float32(n-1)*l.Offset
		return fyne.NewSize(maxVerticalCardWidth, h)
	}

	// North/South
	w := maxHorizontalCardWidth + float32(n-1)*l.Offset
	return fyne.NewSize(w, maxHorizontalCardHeight)
}
