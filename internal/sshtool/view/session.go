package sshtool

import "github.com/gdamore/tcell/v2"

func (v *View) handleSessionKey(ev *tcell.EventKey) *tcell.EventKey {
	if v.sessionStdin == nil {
		return ev
	}

	switch ev.Key() {
	case tcell.KeyRune:
		_, _ = v.sessionStdin.Write([]byte(string(ev.Rune())))
		return nil
	case tcell.KeyEnter:
		_, _ = v.sessionStdin.Write([]byte("\n"))
		return nil
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		_, _ = v.sessionStdin.Write([]byte{0x08})
		return nil
	case tcell.KeyCtrlC:
		if v.sessionCancel != nil {
			v.sessionCancel()
		}
		return nil
	}

	return ev
}

func (v *View) AppendSessionLine(line string) {
	if v.sessionTerm == nil {
		return
	}
	v.sessionTerm.PrintLine(line)
}

func (v *View) ClearSession() {
	if v.sessionTerm == nil {
		return
	}
	v.sessionTerm.ClearScreen()
}
