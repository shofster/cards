package chess

import "fyne.io/fyne/v2"

/*

  File:    Chess_resourceMap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:

fyne bundle -o pieces.go --package chess pieces

<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100">
    <g>
        <rect x="0" y="0" width="100" height="100" fill="none" stroke="black" stroke-width="1"/>
    </g>
</svg>



*/

var ChessResource = map[string]fyne.Resource{
	"gK":    resourceGKingPng,
	"gQ":    resourceGQueenPng,
	"gB":    resourceGBishopPng,
	"gR":    resourceGRookPng,
	"gP":    resourceGPawnPng,
	"gHl":   resourceGKnightPng,
	"gHr":   resourceGKnightRPng,
	"sK":    resourceSKingPng,
	"sQ":    resourceSQueenPng,
	"sB":    resourceSBishopPng,
	"sR":    resourceSRookPng,
	"sP":    resourceSPawnPng,
	"sHl":   resourceSKnightPng,
	"sHr":   resourceSKnightRPng,
	"empty": resourceEmptyPng,
}
