package meld

import "fyne.io/fyne/v2"

/*

  File:    Meld_reourceMap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:

run in cards\assets\riichi\meld
fyne bundle -o meld.go --package meld .


*/

var MeldResource = map[string]fyne.Resource{
	"pair":  resourcePairPng,
	"pon":   resourcePonPng,
	"kan":   resourceKanPng,
	"chi":   resourceChiPng,
	"blank": resourceBlankPng,
	"empty": resourceEmptyPng,
	"ron":   resourceRonPng,
	"draw":  resourceDrawPng,
}
