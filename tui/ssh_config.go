package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepName step = iota
	stepKeyPath
	stepDone
)

type SSHConfigModel struct {
	inputs   []textinput.Model
	focusIdx int
	step     step
	name     string
	keypath  string
}

func NewSSHConfigModel() *SSHConfigModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter your name"
	nameInput.CharLimit = 40
	nameInput.Width = 40

	pathInput := textinput.New()
	pathInput.Placeholder = "Enter your path"
	pathInput.CharLimit = 100
	pathInput.Width = 60

	return &SSHConfigModel{
		inputs:   []textinput.Model{nameInput, pathInput},
		focusIdx: 0,
		step:     stepName,
	}
}

func (m *SSHConfigModel) Init() tea.Cmd {
	m.inputs[0].Focus()
	// return textinput.Blink
	return nil
}

func (m *SSHConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			switch m.step {
			case stepName:
				m.name = m.inputs[0].Value()
				m.step = stepKeyPath
				m.focusIdx = 1
				m.inputs[1].Focus()
				m.inputs[0].Blur()
			case stepKeyPath:
				m.keypath = m.inputs[1].Value()
				m.step = stepDone
			case stepDone:
				return m, Pop()
			default:
				return m, Pop()
			}
		case tea.KeyCtrlC, tea.KeyCtrlQ, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *SSHConfigModel) View() string {
	switch m.step {
	case stepName:
		return fmt.Sprintf("Name:\n%s\n\n(press Enter to continue)", m.inputs[0].View())
	case stepKeyPath:
		return fmt.Sprintf("KeyPath:\n%s\n\n(press Enter to continue)", m.inputs[1].View())
	case stepDone:
		return fmt.Sprintf("Done!\n\nName: %s\nKeyPath: %s\n\npress any key to exit.", m.name, m.keypath)
	default:
		return "Unkown Step"
	}
}
