package common

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    tracker.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var _ desktop.Mouseable = (*Tracker)(nil)
var _ desktop.Hoverable = (*Tracker)(nil)

type Tracker struct {
	widget.BaseWidget
	LastPos fyne.Position
	Event   func(*desktop.MouseEvent)
	Click   func(bool, *desktop.MouseEvent) // down, event
}

func NewTracker() *Tracker {
	t := &Tracker{}
	t.ExtendBaseWidget(t)
	return t
}

/* ---- click events ---- */

func (t *Tracker) MouseDown(ev *desktop.MouseEvent) {
	//fmt.Printf("MouseDown: Button %d, %d, X%0f, Y%0f\n", ev.Button, ev.Modifier,
	//	ev.AbsolutePosition.X, ev.AbsolutePosition.Y)
	t.sendClick(true, ev)
}

func (t *Tracker) MouseUp(ev *desktop.MouseEvent) {
	//fmt.Printf("MouseUp: Button %d, %d, X%0f, Y%0f\n", ev.Button, ev.Modifier,
	//	ev.AbsolutePosition.X, ev.AbsolutePosition.Y)
	t.sendClick(false, ev)
}

func (t *Tracker) MouseMoved(ev *desktop.MouseEvent) {
	t.sendEvent(ev)
}

func (t *Tracker) MouseIn(ev *desktop.MouseEvent) {
	t.sendEvent(ev)
}

func (t *Tracker) MouseOut() {
	t.sendEvent(nil)
}

func (t *Tracker) sendEvent(ev *desktop.MouseEvent) {
	if ev != nil {
		t.LastPos = ev.Position
	}
	if t.Event != nil {
		t.Event(ev)
	}
}

func (t *Tracker) sendClick(down bool, ev *desktop.MouseEvent) {
	if ev != nil {
		t.LastPos = ev.Position
	}
	if t.Click != nil {
		t.Click(down, ev)
	}
}

func (t *Tracker) CreateRenderer() fyne.WidgetRenderer {
	r := canvas.NewRectangle(color.Transparent)
	return &trackerRenderer{
		rect: r,
		t:    t,
	}
}

type trackerRenderer struct {
	rect *canvas.Rectangle
	t    *Tracker
}

func (r *trackerRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *trackerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *trackerRenderer) Refresh() {
	canvas.Refresh(r.rect)
}

func (r *trackerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect}
}

func (r *trackerRenderer) Destroy() {}
