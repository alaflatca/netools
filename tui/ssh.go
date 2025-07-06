package tui

import (
	"fmt"
	"log"
	"netools/internal/db"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type sshItem struct {
	name    string
	keyPath string
}

type SSHModel struct {
	items  []sshItem
	cursor int
}

func NewSSHModel() *SSHModel {
	return &SSHModel{}
}

func (m *SSHModel) Init() tea.Cmd {
	configs, err := db.SelectSSHConfigs()
	if err != nil {
		log.Printf("[ssh] failed to select configs: %v", err)
	}

	m.items = m.items[:0]
	m.items = append(m.items, sshItem{name: "+ add config", keyPath: ""})
	for _, config := range configs {
		m.items = append(m.items, sshItem{name: config.Name, keyPath: config.KeyPath})
	}

	return nil
}

func (m *SSHModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "b":
			return m, Pop()
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				return m, Push(NewSSHConfigModel())
			}
		}
	}
	return m, nil
}

func (m *SSHModel) View() string {
	var b strings.Builder
	b.WriteString("\t\t[ SSH ]\n")

	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		var s string
		if i == 0 {
			s = fmt.Sprintf("%s%s\n", cursor, item.name)
		} else {
			s = fmt.Sprintf("%s%s, %s\n", cursor, item.name, item.keyPath)
		}
		b.WriteString(s)
	}

	return b.String()
}
