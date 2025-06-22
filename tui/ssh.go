package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type sshItem struct {
	name string
	path string
}

type SSHModel struct {
	items  []sshItem
	cursor int
}

func NewSSHModel() *SSHModel {
	return &SSHModel{
		items: []sshItem{
			{
				name: "+ add config",
				path: "",
			},
		},
	}
}

func (m *SSHModel) Init() tea.Cmd {
	return nil
}

func (m *SSHModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "b":
			return m, pop()
		}
	}
	return m, nil
}

func (m *SSHModel) View() string {
	var b strings.Builder
	b.WriteString("\t\t[ SSH ]\n")

	for i, item := range m.items {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		var s string
		if i == 0 {
			s = fmt.Sprintf("%s%s\n", cursor, item.name)
		} else {
			s = fmt.Sprintf("%s%s / path:%s\n", cursor, item.name, item.path)
		}
		b.WriteString(s)
	}

	return b.String()
}
