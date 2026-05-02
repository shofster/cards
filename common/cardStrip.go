package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    cardStrip.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: CardStrip — N reactive card / trick indicators -- Labels With Data.
*/

type CardStrip struct {
	States binding.StringList
	Labels []*widget.Label
	Empty  string
	nCards int
}

func NewCardStrip(nCards int) *CardStrip {
	cs := &CardStrip{
		States: binding.NewStringList(),
		Labels: make([]*widget.Label, nCards),
		nCards: nCards,
		Empty:  " ",
	}

	// force each element a binding.String
	for i := 0; i < cs.nCards; i++ {
		cs.States.Append(cs.Empty)
	}

	for i := 0; i < cs.nCards; i++ {
		raw, _ := cs.States.GetItem(i)
		item := raw.(binding.String)

		// bind to the state for formatting the label's string
		formatted := binding.NewString()
		l := widget.NewLabelWithData(formatted)
		item.AddListener(binding.NewDataListener(func() {
			v, _ := item.Get()
			formatted.Set("[" + v + "]")
			//l.Refresh()
		}))

		l.Alignment = fyne.TextAlignCenter
		l.TextStyle = fyne.TextStyle{Bold: true}
		cs.Labels[i] = l
	}

	return cs
}

// Set a specific indicator symbol / string.
func (cs *CardStrip) SetCard(i int, v string) {
	cs.States.SetValue(i, v)
}

// Reset all indicators to blank
func (cs *CardStrip) Reset() {
	for i := 0; i < cs.nCards; i++ {
		cs.States.SetValue(i, cs.Empty)
	}
}
