package regular

import "fyne.io/fyne/v2"

/*

  File:    Regular_reourceMap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:

run in cards\assets\riichi\regular
fyne bundle -o regular.go --package regular .


*/

var RegularResource = map[string]fyne.Resource{
	// Back / Blank / Front
	"back":  resourceBackPng,
	"blank": resourceBlankPng,
	"front": resourceFrontPng,
	// Dragons
	"🀄": resourceHakuPng,  // White
	"🀅": resourceHatsuPng, // Green
	"🀆": resourceChunPng,  //Red
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
	"🀋r": resourceMan5DoraPng, // red five
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
	"🀝r": resourcePin5DoraPng, // red five
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
	"🀔r": resourceSou5DoraPng, // red five
	"🀕":  resourceSou6Png,
	"🀖":  resourceSou7Png,
	"🀗":  resourceSou8Png,
	"🀘":  resourceSou9Png,
	// Flowers
	"🀢": resourceFlower1Png, // Flowers 1-4
	"🀣": resourceFlower2Png,
	"🀤": resourceFlower3Png,
	"🀥": resourceFlower4Png,
	// Seasons
	"🀦": resourceSeason1Png, // Seasons 1-4
	"🀧": resourceSeason2Png,
	"🀨": resourceSeason3Png,
	"🀩": resourceSeason4Png,
}

var RegularName = map[string]string{
	// Dragons
	"🀄": "Dragon-White",
	"🀅": "Dragon-Green",
	"🀆": "Dragon-Red",
	// Winds
	"🀀": "Wind-East",
	"🀁": "Wind-South",
	"🀂": "Wind-West",
	"🀃": "Wind-North",
	// Manzu (Characters)
	"🀇":  "Manzu-1",
	"🀈":  "Manzu-2",
	"🀉":  "Manzu-3",
	"🀊":  "Manzu-4",
	"🀋":  "Manzu-5",
	"🀋r": "Manzu-5-Red",
	"🀌":  "Manzu-6",
	"🀍":  "Manzu-7",
	"🀎":  "Manzu-8",
	"🀏":  "Manzu-9",
	// Pinzu (Dots)
	"🀙":  "Pinzu-1",
	"🀚":  "Pinzu-2",
	"🀛":  "Pinzu-3",
	"🀜":  "Pinzu-4",
	"🀝":  "Pinzu-5",
	"🀝r": "Pinzu-5-Red",
	"🀞":  "Pinzu-6",
	"🀟":  "Pinzu-7",
	"🀠":  "Pinzu-8",
	"🀡":  "Pinzu-9",
	// Souzu (Bamboo)
	"🀐":  "Souzu-1",
	"🀑":  "Souzu-2",
	"🀒":  "Souzu-3",
	"🀓":  "Souzu-4",
	"🀔":  "Souzu-5",
	"🀔r": "Souzu-5-Red",
	"🀕":  "Souzu-6",
	"🀖":  "Souzu-7",
	"🀗":  "Souzu-8",
	"🀘":  "Souzu-9",
	// Flowers
	"🀢": "Flower-1",
	"🀣": "Flower-2",
	"🀤": "Flower-3",
	"🀥": "Flower-4",
	// Seasons
	"🀦": "Season-1",
	"🀧": "Season-2",
	"🀨": "Season-3",
	"🀩": "Season-4",
}

var ReverseRegularResource = map[fyne.Resource]string{
	// Back / Blank / Front
	resourceBackPng:  "Back",
	resourceBlankPng: "Blank",
	resourceFrontPng: "Front",
	// Dragons
	resourceHakuPng:  "Dragon-White",
	resourceHatsuPng: "Dragon-Green",
	resourceChunPng:  "Dragon-Red",
	// Winds
	resourceTonPng:  "Wind-East",
	resourceNanPng:  "Wind-South",
	resourceShaaPng: "Wind-West",
	resourcePeiPng:  "Wind-North",
	// Manzu (Characters)
	resourceMan1Png:     "Manzu-1",
	resourceMan2Png:     "Manzu-2",
	resourceMan3Png:     "Manzu-3",
	resourceMan4Png:     "Manzu-4",
	resourceMan5Png:     "Manzu-5",
	resourceMan5DoraPng: "Manzu-5-Red",
	resourceMan6Png:     "Manzu-6",
	resourceMan7Png:     "Manzu-7",
	resourceMan8Png:     "Manzu-8",
	resourceMan9Png:     "Manzu-9",
	// Pinzu (Dots)
	resourcePin1Png:     "Pinzu-1",
	resourcePin2Png:     "Pinzu-2",
	resourcePin3Png:     "Pinzu-3",
	resourcePin4Png:     "Pinzu-4",
	resourcePin5Png:     "Pinzu-5",
	resourcePin5DoraPng: "Pinzu-5-Red",
	resourcePin6Png:     "Pinzu-6",
	resourcePin7Png:     "Pinzu-7",
	resourcePin8Png:     "Pinzu-8",
	resourcePin9Png:     "Pinzu-9",

	// Souzu (Bamboo)
	resourceSou1Png:     "Souzu-1",
	resourceSou2Png:     "Souzu-2",
	resourceSou3Png:     "Souzu-3",
	resourceSou4Png:     "Souzu-4",
	resourceSou5Png:     "Souzu-5",
	resourceSou5DoraPng: "Souzu-5-Red",
	resourceSou6Png:     "Souzu-6",
	resourceSou7Png:     "Souzu-7",
	resourceSou8Png:     "Souzu-8",
	resourceSou9Png:     "Souzu-9",

	// Flowers
	resourceFlower1Png: "Flower-1",
	resourceFlower2Png: "Flower-2",
	resourceFlower3Png: "Flower-3",
	resourceFlower4Png: "Flower-4",

	// Seasons
	resourceSeason1Png: "Season-1",
	resourceSeason2Png: "Season-2",
	resourceSeason3Png: "Season-3",
	resourceSeason4Png: "Season-4",
}
