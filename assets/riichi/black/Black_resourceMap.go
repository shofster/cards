package black

import "fyne.io/fyne/v2"

/*

  File:    Black_resourceMap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:

run in cards\assets\riichi\black
fyne bundle -o black.go --package black .


*/

var BlackResource = map[string]fyne.Resource{
	// Back / Blank / Front
	"back":  resourceBackPng,
	"blank": resourceBlankPng,
	"front": resourceFrontPng,
	// Dragons
	"🀄": resourceHakuPng,  // White Dragon
	"🀅": resourceHatsuPng, // Green Dragon
	"🀆": resourceChunPng,  // Red Dragon
	// Winds
	"🀀": resourceTonPng,  // East
	"🀁": resourceNanPng,  // South
	"🀂": resourceShaaPng, // West
	"🀃": resourcePeiPng,  // North
	// Manzu (Characters / Cracks)
	"🀇":  resourceMan1Png, // Manzu 1-9
	"🀈":  resourceMan2Png,
	"🀉":  resourceMan3Png,
	"🀊":  resourceMan4Png,
	"🀋":  resourceMan5Png,
	"🀋r": resourceMan5DoraPng, // Red 5
	"🀌":  resourceMan6Png,
	"🀍":  resourceMan7Png,
	"🀎":  resourceMan8Png,
	"🀏":  resourceMan9Png,
	// Pinzu (Dots)
	"🀙":  resourcePin1Png, // Pinzu 1-9
	"🀚":  resourcePin2Png,
	"🀛":  resourcePin3Png,
	"🀜":  resourcePin4Png,
	"🀝":  resourcePin5Png,
	"🀝r": resourcePin5DoraPng, // Red 5
	"🀞":  resourcePin6Png,
	"🀟":  resourcePin7Png,
	"🀠":  resourcePin8Png,
	"🀡":  resourcePin9Png,
	// Souzu (Bamboo)
	"🀐":  resourceSou1Png, // Souzu 1-9
	"🀑":  resourceSou2Png,
	"🀒":  resourceSou3Png,
	"🀓":  resourceSou4Png,
	"🀔":  resourceSou5Png,
	"🀔r": resourceSou5DoraPng, // Red 5
	"🀕":  resourceSou6Png,
	"🀖":  resourceSou7Png,
	"🀗":  resourceSou8Png,
	"🀘":  resourceSou9Png,
}
