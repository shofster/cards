package main

import "log"

/*

  File:    logger.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: simple logger. turn off or on
*/

func Logit(cb func() string) {
	if Ctx.devMode {
		log.Printf("%s", cb())
	}
}
