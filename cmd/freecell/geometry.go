package main

import (
	"log"

	"fyne.io/fyne/v2"
)

/*

  File:    geometry.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: methods to compute and evaluate the mouse position
    against the image layout(s).
*/

type Geometry struct {
	FreeGrid       [4]CardImageArea
	FoundationGrid [4]CardImageArea
	CardGrids      [8][]CardImageArea
}

type CardImageArea struct {
	position fyne.Position
	size     fyne.Size
}

// build the current board Geometry (absolute X,Y) - changes during play !!
func buildGeometry() {

	for i, obj := range freeGrid.Objects {
		abs := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj)
		geometry.FreeGrid[i] = CardImageArea{
			position: abs,
			size:     obj.Size(),
		}
	}
	for i, obj := range foundationGrid.Objects {
		abs := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj)
		geometry.FoundationGrid[i] = CardImageArea{
			position: abs,
			size:     obj.Size(),
		}
	}

	for i := 0; i < 8; i++ {
		geometry.CardGrids[i] = nil
		obj := cardsGrid.Objects[i]
		column, ok := obj.(*fyne.Container)
		if !ok {
			log.Println("ERROR: cardsGrid.Objects[", i, "] is not a container:", obj)
			continue
		}
		if len(column.Objects) == 0 {
			continue
		}
		geometry.CardGrids[i] = make([]CardImageArea, 0, len(column.Objects))

		for _, cardObj := range column.Objects {
			abs := fyne.CurrentApp().Driver().AbsolutePositionForObject(cardObj)
			geometry.CardGrids[i] = append(geometry.CardGrids[i], CardImageArea{
				position: abs,
				size:     cardObj.Size(),
			})
		}
	}

}

func isInside(ia CardImageArea, mouse fyne.Position) bool {
	return mouse.X >= ia.position.X &&
		mouse.X <= ia.position.X+ia.size.Width &&
		mouse.Y >= ia.position.Y &&
		mouse.Y <= ia.position.Y+ia.size.Height
}

// find the card index in the card columns visible area
func cardAt(cards []CardImageArea, mouse fyne.Position) int {
	n := len(cards)
	if n == 0 {
		return -1
	}
	// layout parameters
	width := Ctx.cardSizeWidth
	height := Ctx.cardSizeHeight
	off := float32(Ctx.cardOffset) // vertical overlap amount

	// column's top-left corner
	px := cards[0].position.X
	py := cards[0].position.Y

	// start at the last card - visible
	y := py + float32(n-1)*off

	for i := n - 1; i >= 0; i-- {
		// top card - full height
		// others   - only the exposed part
		visibleH := height
		if i < n-1 {
			visibleH = off
		}
		area := CardImageArea{
			position: fyne.NewPos(px, y),
			size:     fyne.NewSize(width, visibleH),
		}
		if isInside(area, mouse) {
			return i
		}
		// look backwards
		y -= off
	}

	return -1
}
