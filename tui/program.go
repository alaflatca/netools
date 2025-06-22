package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type popModelMsg struct{}
type pushModelMsg struct{ model tea.Model }

type ProgramModel struct {
	stack []tea.Model
}

func NewProgramModel() *ProgramModel {
	return &ProgramModel{
		stack: []tea.Model{NewMenuModel()},
	}
}

func (p *ProgramModel) Init() tea.Cmd {
	return nil
}

func (p *ProgramModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Fprintf(os.Stdout, "[DEBUG] Update called with msg: %T\n", msg)
	if len(p.stack) == 0 {
		return p, tea.Quit
	}

	switch m := msg.(type) {
	case popModelMsg:
		if len(p.stack) > 1 {
			p.stack = p.stack[:len(p.stack)-1]
			return p, nil
		}
	case pushModelMsg:
		p.stack = append(p.stack, m.model)
		return p, m.model.Init()
	}

	updated, cmd := p.stack[len(p.stack)-1].Update(msg)
	p.stack[len(p.stack)-1] = updated
	return p, cmd
}

func (p *ProgramModel) View() string {
	if len(p.stack) == 0 || p.stack[len(p.stack)-1] == nil {
		return "⚠️ No screen to display"
	}
	return p.stack[len(p.stack)-1].View()
}

func pop() tea.Cmd {
	return func() tea.Msg {
		return popModelMsg{}
	}
}

func push(m tea.Model) tea.Cmd {
	return func() tea.Msg {
		return pushModelMsg{model: m}
	}
}

func PushWithCmd(m tea.Model, cmd tea.Cmd) tea.Cmd {
	return tea.Batch(
		push(m),
		cmd,
	)
}
