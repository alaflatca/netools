package tui

import (
	"fmt"
	tui "sshtn/tui/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

type actionFunc func() (tea.Model, tea.Cmd)

type modelItem struct {
	title  string
	action actionFunc
}

type model struct {
	cursor   int
	items    []modelItem
	selected string
}

type changeModelMsg tea.Model

func NewMainMenu() tea.Model {
	return initialModel()
}

func initialModel() model {
	return model{
		items: []modelItem{
			{
				title: "SSH",
				action: func() (tea.Model, tea.Cmd) {
					return tui.NewSSHMenu(), nil
				},
			},
			{
				title: "VPN",
				action: func() (tea.Model, tea.Cmd) {
					return nil, nil
				},
			},
			{
				title: "Reverse Proxy",
				action: func() (tea.Model, tea.Cmd) {
					return nil, nil
				},
			},
			{
				title: "Tool",
				action: func() (tea.Model, tea.Cmd) {
					return nil, nil
				},
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				return m, func() tea.Msg {
					return changeModelMsg(tui.NewSSHMenu())
				}
			case 1:
				return m, tea.Quit
			case 2:
				return m, tea.Quit
			case 3:
				return m, tea.Quit
			}
		}
	case changeModelMsg:
		return msg, nil
	}

	return m, nil
}

func (m model) View() string {
	s := "\n\t\t[Network Tools]\n"
	s += "Select a tool:\n"

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, item.title)

	}
	s += "\n↑↓ to navigate • Enter to select • q to quit\n"

	if m.selected != "" {
		s += fmt.Sprintf("\nYou selected: %s\n", m.selected)
	}

	return s
}
