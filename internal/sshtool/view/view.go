package sshtool

import "github.com/rivo/tview"

func NewView() (string, tview.Primitive) {
	view := tview.NewList().
		AddItem("add config", "server config register", '+', nil).
		AddItem("eleven", "dk server test", 1, nil)

	// view.SetTitle("ssh")
	// view := tview.NewList().
	// 	AddItem("Session", "remote session", 1, nil).
	// 	AddItem("Tunnel", "port forwading", 2, nil)
	// view.SetTitle("ssh")

	return "ssh", view
}
