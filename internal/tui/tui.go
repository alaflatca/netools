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
		// 이전 프레임 내용 싹 지우기
		s.Clear()

		s.SetStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDefault))
		return false
	})
	menu := tview.NewList().
		AddItem("ssh", pageSSH, '1', nil).
		AddItem("vpn", pageVPN, '2', nil).
		AddItem("reverse proxy", pageReverse, '3', nil).
		AddItem("tools", pageTools, '4', nil)
	menu.SetBorder(true).SetTitle("  netools  ")

	// logging도 추가 필요

	sshPage := sshtool.NewView(app, db)

	vpnPage := tview.NewTextView()
	vpnPage.SetText("vpn: not implemented yet")
	vpnPage.SetBorder(true)

	reversePage := tview.NewTextView()
	reversePage.SetText("reverse proxy: not implemented yet")
	reversePage.SetBorder(true)

	toolsPage := tview.NewTextView()
	toolsPage.SetText("tools: not implemented yet")
	toolsPage.SetBorder(true)

	pages := tview.NewPages()
	pages.AddPage(pageMain, menu, true, true)
	pages.AddPage(pageSSH, sshPage.Primitive(), true, false)
	pages.AddPage(pageVPN, vpnPage, true, false)
	pages.AddPage(pageReverse, reversePage, true, false)
	pages.AddPage(pageTools, toolsPage, true, false)

	focusTarget := map[string]tview.Primitive{
		pageMain:    menu,
		pageSSH:     sshPage.Primitive(),
		pageVPN:     vpnPage,
		pageReverse: reversePage,
		pageTools:   toolsPage,
	}

	switchTo := func(name string) {
		pages.SwitchToPage(name)
		if v, ok := focusTarget[name]; ok && v != nil {
			app.SetFocus(v)
		}
	}

	menu.SetSelectedFunc(func(i int, main, secondary string, shortcut rune) {
		if secondary != "" {
			switchTo(secondary)
		}
	})

	app.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyEscape:
			switchTo(pageMain)
			return nil
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		}
		return ev
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		return err
	}
	return nil
}
