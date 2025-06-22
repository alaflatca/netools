package tui

import tea "github.com/charmbracelet/bubbletea"

type ReverseModel struct {
}

func NewReverseModel() *ReverseModel {
	return &ReverseModel{}
}

func (p *ReverseModel) Init() tea.Cmd {
	return nil
}

func (p *ReverseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (p *ReverseModel) View() string {
	return ""
}
