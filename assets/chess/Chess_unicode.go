package chess

/*

  File:    Chess_unicode.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: UNICODE's for dark / light pieces.
*/

var ChessUnicode = map[string]string{
	"gK":    "♚",
	"gQ":    "♛",
	"gR":    "♜",
	"gB":    "♝",
	"gHl":   "♞",
	"gHr":   "♞", // :-(( NO right facing unicode
	"gP":    "♟",
	"sK":    "♔",
	"sQ":    "♕",
	"sR":    "♖",
	"sB":    "♗",
	"sHl":   "♘",
	"sHr":   "♘", // ditto
	"sP":    "♙",
	"empty": " ",
}
