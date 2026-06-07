package main

import (
	"cards/common"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
)

/*

  File:    score.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

func scorePoints(who, points int, reason string) {
	duration := time.Second * 1
	switch who {
	case PLAYER:
		common.IntSet(Ctx.playerScore, common.IntGet(Ctx.playerScore)+points)
		spiral.advanceA(points)
		common.Toast(Ctx.mainWindow.Canvas(), fyne.Position{
			X: 10,
			Y: 10,
		}, fmt.Sprintf("PLAYER: %d (%s)", points, reason), common.FailColor, duration)
	case AI:
		common.IntSet(Ctx.aiScore, common.IntGet(Ctx.aiScore)+points)
		spiral.advanceB(points)
		common.Toast(Ctx.mainWindow.Canvas(), fyne.Position{
			X: 10,
			Y: 50,
		}, fmt.Sprintf("AI: %d (%s)", points, reason), common.InfoColor, duration)
	}
}

/*
-------------------------
| SCORING PHASE         |
-------------------------
Score hands in strict order:
 1. NON-DEALER hand
 2. DEALER hand
 3. CRIB (belongs to dealer this round)

For each hand:
  - Count 15s
  - Count pairs
  - Count runs
  - Flush (4 or 5)
  - His Nobs (Jack matching starter suit)

Add points to appropriate score
*/

func scoreAction() {

	dealButton.Enable()
	discardButtonState(-1)
	peggingButtonState(-1)
	scoreButton.Disable()
	showAll()

	//runSpiral = true
	//demo()

	if dealer == AI {
		// PLAYER hand
		p := scoreHand(keeps[PLAYER], starter)
		scorePoints(PLAYER, p, "Hand")
		if common.IntGet(Ctx.playerScore) < 121 {
			// AI hand
			a := scoreHand(keeps[AI], starter)
			// Crib
			c := scoreHand(crib, starter)
			scorePoints(AI, a+c, "Crib")
		}
	} else {
		// AI hand
		a := scoreHand(keeps[AI], starter)
		scorePoints(AI, a, "Hand")
		if common.IntGet(Ctx.aiScore) < 121 {
			// PLAYER hand
			p := scoreHand(keeps[PLAYER], starter)
			// Crib
			c := scoreHand(crib, starter)
			scorePoints(PLAYER, p+c, "Crib")
		}

	}
	// --- game over? ---
	if common.IntGet(Ctx.playerScore) >= 121 ||
		common.IntGet(Ctx.aiScore) >= 121 {
		dealButton.Disable()
		return
	}

}

// makeWinSVG creates an SVG fyne.Resource saying:
//
//	<name>
//	WON
func makeWinSVG(name string) fyne.Resource {
	svg := fmt.Sprintf(`
<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200">
  <rect x="0" y="0" width="400" height="200" fill="#101020"/>
  <text x="50%%" y="45%%"
        fill="#FFD700"
        font-size="32"
        text-anchor="middle"
        dominant-baseline="middle"
        font-family="sans-serif">%s</text>
  <text x="50%%" y="70%%"
        fill="#FFFFFF"
        font-size="28"
        text-anchor="middle"
        dominant-baseline="middle"
        font-family="sans-serif">WON</text>
</svg>`, name)

	return fyne.NewStaticResource("winner.svg", []byte(svg))
}
