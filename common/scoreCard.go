package common

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    scoreCard.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: ScoreCard — N x M reactive card / trick indicators -- Labels With Data.
*/

type ScoreCard struct {
	Keys   []string
	States map[string][]binding.String
	Empty  string
	Grid   *fyne.Container
}

func NewScoreCard(keys []string) *ScoreCard {
	sc := &ScoreCard{
		Keys:   keys,
		States: make(map[string][]binding.String),
		Empty:  " * ",
	}

	sc.Grid = container.NewGridWithColumns(len(keys))
	for k := 0; k < len(keys); k++ {
		sc.Grid.Add(widget.NewLabel(
			fmt.Sprintf(" %s ", sc.Keys[k])))
	}

	// allocate slices for each position
	for _, key := range sc.Keys {
		sc.States[key] = make([]binding.String, 13)
	}

	for i := 0; i < 13; i++ {
		// "#" trick numner
		tn := binding.NewString()
		StringSet(tn, fmt.Sprintf("%3d", i+1))
		sc.States[sc.Keys[0]][i] = tn
		sc.Grid.Add(widget.NewLabelWithData(tn))

		// suits : South (,West), North, and East
		for _, pos := range sc.Keys[1:] {
			bs := binding.NewString()
			bs.Set(sc.Empty)
			sc.States[pos][i] = bs
			sc.Grid.Add(widget.NewLabelWithData(bs))
		}
	}

	return sc
}

// Set a specific indicator symbol / string.
func (sc *ScoreCard) SetCard(trick int, pos string, v string) {
	_ = sc.States[pos][trick].Set(v)
}

// Reset all indicators to blank
func (sc *ScoreCard) Reset() {
	for _, pos := range sc.Keys[1:] {
		for i := 0; i < 13; i++ {
			sc.States[pos][i].Set(sc.Empty)
		}
	}
}
