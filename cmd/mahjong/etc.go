package main

import (
	"log"

	"fyne.io/fyne/v2"
)

/*

  File:    etc.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: simple helper functions,
*/

// simple logger. turn off or on
func Logit(cb func() string) {
	if Ctx.devMode {
		log.Printf("%s", cb())
	}
}

/*
	ops must be wrapped with code similar to:
	var ops = make([]func(), 0)
	defer fyne.Do(func() {
		doWait(&ops, false)
	})
	....................
		ops = append(ops, func() { widget.DoSomething })
		ops = append(ops, func() { someContainer.Refresh() })
		ops = append(ops, func() { someOther.Refresh() })
*/
// wait on a series of FYNE.io non-UI thread updates
func doWait(ops *[]func(), wait bool) {
	if len(*ops) == 0 {
		return
	}
	if wait {
		for _, op := range *ops {
			fyne.DoAndWait(op)
		}
	} else {
		for _, op := range *ops {
			fyne.Do(op)
		}
	}
}
