package tui

import (
	"context"
	reversetool "netools/internal/reversetool/view"
	sshtool "netools/internal/sshtool/view"
	vpntool "netools/internal/vpntool/view"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Start(ctx context.Context) error {
	app := tview.NewApplication()
	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.SetStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDefault))
		return false
	})

	ssh := sshtool.NewView(app)
	vpnName, vpnView := vpntool.NewView()
	reverseName, reverseView := reversetool.NewView()

	side := tview.NewList().
		AddItem(ssh.Name, "", '1', nil).
		AddItem(vpnName, "", '2', nil).
		AddItem(reverseName, "", '3', nil)
	side.SetBorder(true)
	side.SetTitle(" netools ")

	pages := tview.NewPages()
	pages.AddPage(ssh.Name, ssh.FocusDefault(), true, true)

	focusTarget := map[string]tview.Primitive{
		ssh.Name:    ssh.FocusDefault(),
		vpnName:     vpnView,
		reverseName: reverseView,
	}

	side.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		pages.SetTitle(" " + mainText + " ")
		pages.SwitchToPage(mainText)
	})

	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyLeft, tcell.KeyExit:
			app.SetFocus(side)
			return nil
		case tcell.KeyRight, tcell.KeyTab, tcell.KeyEnter:
			name, _ := pages.GetFrontPage()
			if page, ok := focusTarget[name]; ok {
				app.SetFocus(page)
			}

		}
		return ev
	})

	layout := tview.NewFlex().
		AddItem(side, 24, 0, true).
		AddItem(pages, 0, 1, false)

	if err := app.SetRoot(layout, true).Run(); err != nil {
		return err
	}
	return nil
}
