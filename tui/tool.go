package tui

import tea "github.com/charmbracelet/bubbletea"

type ToolModel struct {
}

func NewToolModel() *ToolModel {
	return &ToolModel{}
}

func (p *ToolModel) Init() tea.Cmd {
	return nil
}

func (p *ToolModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (p *ToolModel) View() string {
	return ""
}
