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
	ip      string
	port    string
	keyPath string
	desc    string
}

type SSHModel struct {
	items  []sshItem
	cursor int

	page      int
	limit     int
	totalPage int
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
		m.items = append(m.items, sshItem{name: config.Name, ip: config.IP, port: config.Port, keyPath: config.KeyPath, desc: config.Desc})
	}

	m.page = 0
	m.limit = 5
	if m.limit > len(m.items) {
		m.limit = len(m.items)
	}
	m.totalPage = len(configs) / 5
	if len(configs)%5 > 0 {
		m.totalPage += 1
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
			if m.cursor < m.limit-1 {
				m.cursor++
			}
		case "right":
			if m.page < m.totalPage-1 {
				m.page++
			}
		case "left":
			if m.page > 0 {
				m.page--
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

	page := m.page * m.limit
	limit := page + m.limit

	for i, item := range m.items[page:limit] {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		var s string
		if i == 0 {
			s = fmt.Sprintf("%s%s\n", cursor, item.name)
		} else {
			if len(item.desc) > 0 {
				s = fmt.Sprintf("%s%s[%s:%s],\t%s (%s)\n", cursor, item.name, item.ip, item.port, item.keyPath, item.desc)
			} else {
				s = fmt.Sprintf("%s%s[%s:%s],\t%s\n", cursor, item.name, item.ip, item.port, item.keyPath)
			}
		}
		b.WriteString(s)
	}

	pages := "\n  "
	for dot := range m.totalPage {
		if dot == m.page {
			pages += "•"
		} else {
			pages += "."
		}
	}
	b.WriteString(pages)

	return b.String()
}
