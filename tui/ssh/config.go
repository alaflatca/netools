package tui

import (
	"fmt"
	"sshtn/common"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type sshConfig struct {
	focusIndex int
	inputs     []textinput.Model
	errors     []string
}

func NewSSHConfig() *sshConfig {
	m := &sshConfig{
		inputs: make([]textinput.Model, 3),
		errors: make([]string, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.EchoMode = textinput.EchoNormal

		switch i {
		case 0:
			t.Prompt = "Username: "
			t.CharLimit = 30
			t.Focus()
		case 1:
			t.Prompt = "Keypath: "
			t.CharLimit = 30
		case 2:
			t.Prompt = "\n[ Save ]"
			t.CharLimit = 0
		}

		m.inputs[i] = t
	}

	return m
}

func (m *sshConfig) Init() tea.Cmd {
	return textinput.Blink
}

func (m *sshConfig) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				if m.validate() {
					username := m.inputs[0].Value()
					keypath := m.inputs[1].Value()
					fmt.Printf("Saving: %s / %s\n", username, keypath)
					return m, tea.Quit
				}
				return m, nil
			}

			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *sshConfig) validate() bool {
	valid := true
	username := strings.TrimSpace(m.inputs[0].Value())
	keypath := strings.TrimSpace(m.inputs[1].Value())

	m.errors[0] = ""
	m.errors[1] = ""

	if username == "" {
		m.errors[0] = "Username cannot be empty"
		valid = false
	}

	if keypath == "" {
		m.errors[1] = "Keypath cannot be empty"
		valid = false
	}

	if !common.FileExist(keypath) {
		m.errors[1] = fmt.Sprintf("%q is not exist", keypath)
		valid = false
	}

	return valid
}

func (m *sshConfig) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *sshConfig) View() string {
	var b strings.Builder

	b.WriteString("\n\t\t[SSH Tools]\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
		if err := m.errors[i]; err != "" {
			b.WriteString("  ⚠ " + err + "\n")
		}
	}
	b.WriteString("\n↑/↓ or tab: move • enter: save • esc/q: quit\n")

	return b.String()
}
