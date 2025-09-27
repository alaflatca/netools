package vpntool

import "github.com/rivo/tview"

func NewView() (string, tview.Primitive) {
	view := tview.NewList().
		AddItem("Server", "vpn server start", 1, nil).
		AddItem("Client", "vpn client start", 2, nil)
	view.SetTitle("vpn")

	return "vpn", view
}
