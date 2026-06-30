package main

import (
	"cards/common"
	"fmt"
)

/*

  File:    cellitem.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: struct with an ID that represents a Card's state and position.
*/

type Location int

const (
	Empty Location = iota
	Free
	Foundation
	Cards
)

// a combined int for use as a global ID
type CellItemId int

type CellItem struct {
	id        CellItemId // logical position
	card      common.Card
	scoreRank int // value 1 dest 10 - for scoring moves
	empty     bool
}

func NewCellItemId(loc Location, col, ix int) CellItemId {
	return CellItemId(int(loc<<16) | col<<8 | ix)
}

func (id CellItemId) location() Location {
	return Location(id >> 16)
}
func (id CellItemId) column() int {
	return int((id >> 8) & 0xff)
}
func (id CellItemId) index() int {
	return int(id & 0xff)
}
func decodeCellItemId(id CellItemId) (Location, int, int) {
	return id.location(), id.column(), id.index()
}

// pretty debug string
func (id CellItemId) String() string {
	switch id.location() {
	case Empty:
		return fmt.Sprintf("Empty:")
	case Free:
		return fmt.Sprintf("Free:%d", id.column())
	case Foundation:
		return fmt.Sprintf("Foundation:%d", id.column())
	case Cards:
		return fmt.Sprintf("Card:%d:%d", id.column(), id.index())
	}
	return ("CellItem:?")
}

func NewCellItem(location Location, col, ix int, card common.Card) *CellItem {
	ci := &CellItem{card: card}
	ci.id = NewCellItemId(location, col, ix)
	switch card.Rank {
	case common.Ace:
		ci.scoreRank = 1
	case common.Jack, common.Queen, common.King:
		ci.scoreRank = 10
	default:
		ci.scoreRank = int(card.Rank)
	}
	return ci
}

func (id CellItemId) IsEmpty() bool {
	return id.location() == Empty
}
func (id CellItemId) IsCardsCell() bool {
	return id.location() == Cards
}
func (id CellItemId) IsFoundationCell() bool {
	return id.location() == Foundation
}
func (id CellItemId) IsFreeCell() bool {
	return id.location() == Free
}

const EmptyCardIcon = "empty"

func emptyCard() common.Card {
	return common.Card{
		Rank:     common.Rank(-1),
		Suit:     common.NoSuit,
		Icon:     EmptyCardIcon,
		Resource: Ctx.cardMap["reverse"].Resource,
	}
}

func aceCard(col int) common.Card {
	card := common.Card{}
	switch col {
	case 0:
		card.Icon = "A♣"
	case 1:
		card.Icon = "A♦"
	case 2:
		card.Icon = "A♥"
	case 3:
		card.Icon = "A♠"
	}
	card.Resource = Ctx.cardMap[card.Icon].Resource
	return card
}
