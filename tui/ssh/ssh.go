package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type sshModel struct {
	cursor  int
	choices []string
}

func NewSSHMenu() tea.Model {
	return sshModel{
		choices: []string{
			"+ add config",
		},
	}
}

func (m sshModel) Init() tea.Cmd {
	return nil
}

func (m sshModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "K":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				return NewSSHConfig(), nil
			case 1:
				return m, tea.Quit
			case 2:
				return m, tea.Quit
			case 3:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m sshModel) View() string {
	s := "\n\t\t[SSH Tools]\n"
	s += "Select a config:\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)

	}
	s += "\n↑↓ to navigate • Enter to select • q to quit\n"

	// if m.selected != "" {
	// 	s += fmt.Sprintf("\nYou selected: %s\n", m.selected)
	// }

	return s
}
