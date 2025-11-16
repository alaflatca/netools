package tui

import (
	"context"
	"netools/internal/db"
	sshtool "netools/internal/sshtool/view"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Start(ctx context.Context, db *db.DB) error {
	const (
		pageMain    = "main"
		pageSSH     = "ssh"
		pageVPN     = "vpn"
		pageReverse = "reverse"
		pageTools   = "tools"
	)

	app := tview.NewApplication()
	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.SetStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDefault))
		return false
	})

	menu := tview.NewList().
		AddItem("ssh", "", '1', nil).
		AddItem("vpn", "", '2', nil).
		AddItem("reverse proxy", "", '3', nil).
		AddItem("tools", "", '4', nil)
	menu.SetBorder(true).SetTitle("  netools  ")

	// logging도 추가 필요
	sshPage := sshtool.NewView(app, db)

	pages := tview.NewPages()
	pages.AddPage(pageMain, menu, true, true)
	pages.AddPage(pageSSH, sshPage.Root, true, false)

	focusTarget := map[string]tview.Primitive{
		pageMain:    menu,
		pageSSH:     sshPage.Root,
		pageVPN:     nil,
		pageReverse: nil,
		pageTools:   nil,
	}

	switchTo := func(name string) {
		pages.SwitchToPage(name)
		if v, ok := focusTarget[name]; ok {
			app.SetFocus(v)
		}
	}

	menu.SetSelectedFunc(func(i int, main, secondary string, shortcut rune) {
		pages.SwitchToPage(secondary)
	})

	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyEscape:
			switchTo(pageMain)
			return nil
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyRune:
			switch ev.Rune() {
			case '1':
				switchTo(pageSSH)
				return nil
			case '2':
				switchTo(pageVPN)
				return nil
			case '3':
				switchTo(pageReverse)
				return nil
			case '4':
				switchTo(pageTools)
				return nil
			}
		}
		return ev
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
