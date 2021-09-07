package oura

import (
	"github.com/gdamore/tcell"
)

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeHelpTextKeyboardControl(widget.ShowHelp)

	widget.SetKeyboardKey(tcell.KeyRight, widget.nextPage, "Next page")
	widget.SetKeyboardKey(tcell.KeyLeft, widget.prevPage, "Prev page")
}

func (widget *Widget) nextPage() {
	widget.idx++
	if widget.idx == len(widget.pageTypes) {
		widget.idx = 0
	}
	widget.Refresh()
}

func (widget *Widget) prevPage() {
	widget.idx--
	if widget.idx < 0 {
		widget.idx = len(widget.pageTypes) - 1
	}
	widget.Refresh()
}
