package checkers

import "fyne.io/fyne/v2"

/*

  File:    Checkers_resourceMap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:

fyne bundle -o pieces.go --package checkers pieces


*/

const (
	BlackName = "BLACK"
	RedName   = "RED"
)

var CheckerResource = map[string]fyne.Resource{
	"black":      resourceBlackCircPng,
	"blackCrown": resourceBlackCircCrownPng,
	"red":        resourceRedHexPng,
	"redCrown":   resourceRedHexCrownPng,
	"empty":      resourceEmptyPng,
}
