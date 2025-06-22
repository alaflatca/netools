package tui

import tea "github.com/charmbracelet/bubbletea"

type VPNModel struct {
}

func NewVPNModel() *VPNModel {
	return &VPNModel{}
}

func (p *VPNModel) Init() tea.Cmd {
	return nil
}

func (p *VPNModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (p *VPNModel) View() string {
	return ""
}
