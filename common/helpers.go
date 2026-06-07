package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    helpers.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: various functions.
*/

// wrapper to make a BOLD Label
func BoldLabel(l string) *widget.Label {
	lab := widget.NewLabel(l)
	lab.TextStyle = fyne.TextStyle{Bold: true}
	return lab
}

//
// wrappers for binding variables -- ignoring error
//

func StringGet(st binding.String) string {
	v, _ := st.Get()
	return v
}

func StringSet(st binding.String, v string) string {
	_ = st.Set(v)
	return v
}

func IntGet(bi binding.Int) int {
	v, _ := bi.Get()
	return v
}

func IntSet(st binding.Int, v int) int {
	_ = st.Set(v)
	return v
}

func BoolGet(bi binding.Bool) bool {
	v, _ := bi.Get()
	return v
}

func BoolSet(st binding.Bool, v bool) bool {
	_ = st.Set(v)
	return v
}
func Ternary[T any](cond bool, yes, no T) T {
	if cond {
		return yes
	}
	return no
}

// IfElse. Not a generic.
// closure potentially can cause side effects within its scope.
func IfElse(cond bool, yes, no func()) {
	if cond {
		yes()
	} else {
		no()
	}
}
