package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	Name   string
	Action func() (tea.Model, tea.Cmd)
}

type MenuModel struct {
	Items  []Item
	Cursor int
}

func NewMenuModel() *MenuModel {
	return &MenuModel{
		Items: []Item{
			{
				Name: "SSH",
				Action: func() (tea.Model, tea.Cmd) {
					return NewSSHModel(), nil
				},
			},
			{
				Name: "VPN",
				Action: func() (tea.Model, tea.Cmd) {
					return NewVPNModel(), nil
				},
			},
			{
				Name: "Reverse Proxy",
				Action: func() (tea.Model, tea.Cmd) {
					return NewReverseModel(), nil
				},
			},
			{
				Name: "Tool",
				Action: func() (tea.Model, tea.Cmd) {
					return NewToolModel(), nil
				},
			},
		},
	}
}

func (p *MenuModel) Init() tea.Cmd {
	return nil
}

func (p *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return p, tea.Quit
		case "b": // pop
			return p, Pop()
		case "up":
			if p.Cursor > 0 {
				p.Cursor--
			}
		case "down":
			if p.Cursor < len(p.Items)-1 {
				p.Cursor++
			}
		case "enter": // push
			item := p.Items[p.Cursor]
			if item.Action != nil {
				model, innerCmd := item.Action()
				return p, PushWithCmd(model, innerCmd)
			}
		}
	}
	return p, nil
}

func (p *MenuModel) View() string {
	var b strings.Builder
	b.WriteString("\n\t\t[ Netools ]\n\n")

	for i, item := range p.Items {
		cursor := "  "
		if p.Cursor == i {
			cursor = "> "
		}
		s := fmt.Sprintf("%s%d. %s\n", cursor, i+1, item.Name)
		b.WriteString(s)
	}

	b.WriteString("\n[↑/↓ to change • Enter to select • q to quit]")

	return b.String()
}
