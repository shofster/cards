package assets

import "fyne.io/fyne/v2"

/*

  File:    BN_resourceMap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:

fyne bundle -o cards.go --package assets ByronKnollCards

*/

var ResourceByCode = map[string]fyne.Resource{
	"2♣":      resource2ofclubsPng,
	"3♣":      resource3ofclubsPng,
	"4♣":      resource4ofclubsPng,
	"5♣":      resource5ofclubsPng,
	"6♣":      resource6ofclubsPng,
	"7♣":      resource7ofclubsPng,
	"8♣":      resource8ofclubsPng,
	"9♣":      resource9ofclubsPng,
	"10♣":     resource10ofclubsPng,
	"J♣":      resourceJackofclubsPng,
	"Q♣":      resourceQueenofclubsPng,
	"K♣":      resourceKingofclubsPng,
	"A♣":      resourceAceofclubsPng,
	"2♦":      resource2ofdiamondsPng,
	"3♦":      resource3ofdiamondsPng,
	"4♦":      resource4ofdiamondsPng,
	"5♦":      resource5ofdiamondsPng,
	"6♦":      resource6ofdiamondsPng,
	"7♦":      resource7ofdiamondsPng,
	"8♦":      resource8ofdiamondsPng,
	"9♦":      resource9ofdiamondsPng,
	"10♦":     resource10ofdiamondsPng,
	"J♦":      resourceJackofdiamondsPng,
	"Q♦":      resourceQueenofdiamondsPng,
	"K♦":      resourceKingofdiamondsPng,
	"A♦":      resourceAceofdiamondsPng,
	"2♥":      resource2ofheartsPng,
	"3♥":      resource3ofheartsPng,
	"4♥":      resource4ofheartsPng,
	"5♥":      resource5ofheartsPng,
	"6♥":      resource6ofheartsPng,
	"7♥":      resource7ofheartsPng,
	"8♥":      resource8ofheartsPng,
	"9♥":      resource9ofheartsPng,
	"10♥":     resource10ofheartsPng,
	"J♥":      resourceJackofheartsPng,
	"Q♥":      resourceQueenofheartsPng,
	"K♥":      resourceKingofheartsPng,
	"A♥":      resourceAceofheartsPng,
	"2♠":      resource2ofspadesPng,
	"3♠":      resource3ofspadesPng,
	"4♠":      resource4ofspadesPng,
	"5♠":      resource5ofspadesPng,
	"6♠":      resource6ofspadesPng,
	"7♠":      resource7ofspadesPng,
	"8♠":      resource8ofspadesPng,
	"9♠":      resource9ofspadesPng,
	"10♠":     resource10ofspadesPng,
	"J♠":      resourceJackofspadesPng,
	"Q♠":      resourceQueenofspadesPng,
	"K♠":      resourceKingofspadesPng,
	"A♠":      resourceAceofspadesPng,
	"reverse": resourceBackredburstPng,
	"joker":   resourceJokerredPng,
}
