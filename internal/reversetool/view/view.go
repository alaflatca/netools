package reversetool

import "github.com/rivo/tview"

func NewView() (string, tview.Primitive) {
	view := tview.NewList().
		AddItem("Register", "", 1, nil).
		AddItem("Start", "", 2, nil).
		AddItem("Log", "", 3, nil)
	view.SetTitle("reverse")

	return "reverse", view
}
