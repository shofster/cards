package main

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

/*

  File:    spiral.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var DefaultColorA = color.NRGBA{220, 50, 50, 30}
var DefaultColorB = color.NRGBA{50, 50, 220, 30}
var BaseColor = color.NRGBA{40, 128, 40, 30}

type SpiralTracks struct {
	c      *fyne.Container
	a, b   []*canvas.Circle
	posA   int
	posB   int
	colorA color.NRGBA
	colorB color.NRGBA
	Delay  int // animation delay in MS
}

func NewSpiralTracks(s float32, holes int, size fyne.Size) *SpiralTracks {
	if holes <= 0 {
		holes = 121
	}

	tracks := &SpiralTracks{
		colorA: DefaultColorA,
		colorB: DefaultColorB,
		Delay:  200,
	}

	// peg radius
	rp := s * 0.25

	// lane separation (radial)
	lane := float64(s) * 1.25 // 20 px for s=16 → visible

	// spiral radii
	rOuter := float64(size.Width) * 0.37
	rInner := rOuter * 0.33 // prevents inner spiral crunching

	turns := 2.5
	cx := float64(size.Width) / 2
	cy := float64(size.Height) / 2

	tracks.c = container.NewWithoutLayout()
	tracks.c.Resize(size)
	tracks.a = make([]*canvas.Circle, holes)
	tracks.b = make([]*canvas.Circle, holes)

	mk := func(col color.Color) *canvas.Circle {
		circ := canvas.NewCircle(col)
		circ.Resize(fyne.NewSize(2*rp, 2*rp))
		return circ
	}

	theta0 := math.Pi
	theta1 := math.Pi + 2*math.Pi*turns

	for i := 0; i < holes; i++ {
		t := float64(i) / float64(holes-1)

		// angle
		theta := theta0 + t*(theta1-theta0)

		// base Archimedean radius
		R := rInner + (rOuter-rInner)*t

		// radial offsets for A/B lanes
		rA := R + lane/2
		rB := R - lane/2

		// coordinates
		xA := cx + rA*math.Cos(theta)
		yA := cy + rA*math.Sin(theta)
		xB := cx + rB*math.Cos(theta)
		yB := cy + rB*math.Sin(theta)

		a := mk(BaseColor)
		b := mk(BaseColor)

		a.Move(fyne.NewPos(float32(xA)-rp, float32(yA)-rp))
		b.Move(fyne.NewPos(float32(xB)-rp, float32(yB)-rp))

		tracks.a[i] = a
		tracks.b[i] = b

		tracks.c.Add(a)
		tracks.c.Add(b)
	}

	return tracks
}

func (st *SpiralTracks) reset() {
	for i := range st.a {
		st.a[i].FillColor = BaseColor
	}
	st.posA = 0
	for i := range st.b {
		st.b[i].FillColor = BaseColor
	}
	st.posB = 0
	fyne.Do(func() { st.c.Refresh() })
}

func alphaForDist(dist int) uint8 {
	switch dist {
	case 0:
		return 255 // bright
	case 1:
		return 200 // glow
	default:
		return 80 // faint
	}
}

func (st *SpiralTracks) advanceA(to int) {
	to-- // make index
	to = min(max(to, 0), len(st.a)-1)
	st.posA = to
	for i := 0; i <= to; i++ {
		d := to - i
		a := alphaForDist(d)
		st.a[i].FillColor = color.NRGBA{
			R: st.colorA.R,
			G: st.colorA.G,
			B: st.colorA.B,
			A: a,
		}
	}
	fyne.Do(func() { st.c.Refresh() })
}

func (st *SpiralTracks) advanceB(to int) {
	to-- // make inde
	to = min(max(to, 0), len(st.b)-1)
	st.posB = to
	for i := 0; i <= to; i++ {
		d := to - i
		a := alphaForDist(d)
		st.b[i].FillColor = color.NRGBA{
			R: st.colorB.R,
			G: st.colorB.G,
			B: st.colorB.B,
			A: a,
		}
	}
	st.posB = to
	fyne.Do(func() { st.c.Refresh() })
}
