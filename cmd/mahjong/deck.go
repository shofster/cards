package main

import (
	"cards/assets/riichi/meld"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

/*

  File:    deck.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: functions to support a "deck" of tiles.
*/

const (
	Manzu      = "Man"
	Pinzu      = "Pin"
	Souzu      = "Sou"
	DragonSuit = "Dragon"
	WindSuit   = "Wind"
)

type Tile struct {
	Code  string
	Suit  string // "Man", "Pin", "Sou"
	Rank  int    // 1–9
	Red   bool   // true = red 5, else false
	ID    int    // index in the hand
	Dummy bool
}

type Tiles []*Tile

func NewTile(code string) *Tile {
	suit, rank, red := parseTile(code)
	t := &Tile{
		Code: code,
		Suit: suit,
		Rank: rank,
		Red:  red,
	}
	return t
}

func BuildDeck(tileMap map[string]fyne.Resource) Tiles {
	var deck Tiles
	for code := range tileMap {
		t := NewTile(code)
		// playable tiles only
		if t.Suit != "?" && t.Rank != 0 {
			deck = append(deck, t)
		}
	}
	return deck
}

func parseTile(code string) (suit string, rank int, red bool) {
	// red suffix - 1 per suit
	red = strings.HasSuffix(code, "r")
	// base code / icon
	icon := code
	if red {
		icon = code[:len(code)-1]
	}
	r := []rune(icon)
	if len(r) == 0 {
		return "?", 0, red
	}
	u := int(r[0])

	switch { // rnak 1 to 9
	// Manzu 🀇–🀏
	case u >= 0x1F007 && u <= 0x1F00F:
		suit = Manzu
		rank = u - 0x1F006
	// Souzu 🀐–🀘
	case u >= 0x1F010 && u <= 0x1F018:
		suit = Souzu
		rank = u - 0x1F00F
	// Pinzu 🀙–🀡
	case u >= 0x1F019 && u <= 0x1F021:
		suit = Pinzu
		rank = u - 0x1F018
	// Winds 🀀–🀃  (4 - E,S,W,N)
	case u >= 0x1F000 && u <= 0x1F003:
		suit = WindSuit
		rank = u - 0x1EFF
	// Dragons 🀄–🀆  (3 - white, green, red)
	case u >= 0x1F004 && u <= 0x1F006:
		suit = DragonSuit
		rank = u - 0x1F003

		/*
			// Flowers 🀢–🀥
			case u >= 0x1F022 && u <= 0x1F025:
				suit = "Flower"
				rank = u - 0x1F021
			// Seasons 🀦–🀩
			case u >= 0x1F026 && u <= 0x1F029:
				suit = "Season"
				rank = u - 0x1F025
		*/

	default: // unused as playable tile
		suit = "?"
		rank = 0
	}

	return suit, rank, red
}

func SortTiles(t Tiles) {
	sort.Slice(t, func(i, j int) bool {
		if t[i].Suit != t[j].Suit {
			return t[i].Suit < t[j].Suit
		}
		if t[i].Rank != t[j].Rank {
			return t[i].Rank < t[j].Rank
		}
		// red after non-red
		if t[i].Red != t[j].Red {
			return !t[i].Red && t[j].Red // false < true
		}
		return false
	})
}

func ShuffleTiles(t Tiles) {
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	r.Shuffle(len(t), func(i, j int) {
		t[i], t[j] = t[j], t[i]
	})
}

// Deal nTiles to nHands, return remaining deck.
func DealTiles(deck Tiles, nHands, nCards int) ([]Tiles, Tiles) {
	hands := make([]Tiles, nHands)
	for i := range hands {
		hands[i] = make([]*Tile, 0)
	}
	idx := 0
	id := 0
	for range nCards {
		for h := range nHands {
			tile := deck[idx]
			tile.ID = id
			hands[h] = append(hands[h], tile)
			idx++ // index in the hand
			if idx >= len(deck) {
				return hands, nil
			}
		}
		id++
	}
	return hands, deck[idx:]
}

// dummy tiles to replace Tile images
func NewEmptyTile(id int) *Tile {
	return &Tile{
		Code:  "empty",
		Suit:  "",
		Rank:  0,
		Red:   false,
		ID:    id,
		Dummy: true,
	}
}
func NewBlankTile(id int) *Tile {
	return &Tile{
		Code:  "blank",
		Suit:  "",
		Rank:  0,
		Red:   false,
		ID:    id,
		Dummy: true,
	}
}

// get an image resource from meld or configured tileMap
func getTileResource(code string) fyne.Resource {
	// check meld and special resources
	res, ok := meld.MeldResource[code]
	if ok {
		return res
	}
	res, ok = tileMap[code]
	if ok {
		return res
	}
	panic(fmt.Sprintf("getTileResource: code %s", code))
}
