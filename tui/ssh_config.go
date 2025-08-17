package tui

import (
	"fmt"
	"log"
	"netools/internal/db"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepName step = iota
	stepIP
	stepPort
	stepKeyPath
	stepDesc
	stepDone
)

type SSHConfigModel struct {
	inputs []textinput.Model

	step    step
	name    string
	ip      string
	port    string
	keypath string
	desc    string

	err error
}

func NewSSHConfigModel() *SSHConfigModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter your name"
	nameInput.CharLimit = 40
	nameInput.Width = 40

	ipInput := textinput.New()
	ipInput.Placeholder = "Enter your ip"
	ipInput.CharLimit = 15
	ipInput.Width = 40

	portInput := textinput.New()
	portInput.Placeholder = "Enter your port"
	portInput.CharLimit = 5
	portInput.Width = 40

	pathInput := textinput.New()
	pathInput.Placeholder = "Enter your path"
	pathInput.CharLimit = 100
	pathInput.Width = 60

	descInput := textinput.New()
	descInput.Placeholder = "Enter your desc"
	descInput.CharLimit = 100
	descInput.Width = 60

	return &SSHConfigModel{
		inputs: []textinput.Model{nameInput, ipInput, portInput, pathInput, descInput},
		step:   stepName,
	}
}

func (m *SSHConfigModel) Init() tea.Cmd {
	m.inputs[0].Focus()
	return textinput.Blink
}

func (m *SSHConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			switch m.step {
			case stepName:
				m.name = m.inputs[stepName].Value()
				m.step = stepIP
				m.inputs[stepIP].Focus()
				m.inputs[stepName].Blur()
			case stepIP:
				m.ip = m.inputs[stepIP].Value()
				m.step = stepPort
				m.inputs[stepPort].Focus()
				m.inputs[stepIP].Blur()
			case stepPort:
				m.port = m.inputs[stepPort].Value()
				m.step = stepKeyPath
				m.inputs[stepKeyPath].Focus()
				m.inputs[stepPort].Blur()
			case stepKeyPath:
				m.keypath = m.inputs[stepKeyPath].Value()
				m.step = stepDesc
				m.inputs[stepDesc].Focus()
				m.inputs[stepKeyPath].Blur()
			case stepDesc:
				m.desc = m.inputs[stepDesc].Value()
				m.step = stepDone
				return m, tea.ClearScreen
			case stepDone:
				err := db.InsertSSHConfig(db.SSHConfig{
					Name:    m.name,
					IP:      m.ip,
					Port:    m.port,
					KeyPath: m.keypath,
					Desc:    m.desc,
				})
				if err != nil {
					log.Printf("[ssh] failed to config insert: %v", err)
					m.err = err
					return m, nil
				}
				return m, Pop()
			default:
				return m, Pop()
			}
		case tea.KeyCtrlC, tea.KeyCtrlQ, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	if m.step < stepDone {
		var cmd tea.Cmd
		m.inputs[m.step], cmd = m.inputs[m.step].Update(msg)
		return m, tea.Batch(cmd)
	}

	return m, nil
}

func (m *SSHConfigModel) View() string {
	switch m.step {
	case stepName:
		return fmt.Sprintf("Name:\n%s\n\n(press Enter to continue)", m.inputs[stepName].View())
	case stepIP:
		return fmt.Sprintf("IP:\n%s\n\n(press Enter to continue)", m.inputs[stepIP].View())
	case stepPort:
		return fmt.Sprintf("Port:\n%s\n\n(press Enter to continue)", m.inputs[stepPort].View())
	case stepKeyPath:
		return fmt.Sprintf("KeyPath:\n%s\n\n(press Enter to continue)", m.inputs[stepKeyPath].View())
	case stepDesc:
		return fmt.Sprintf("Desc:\n%s\n\n(press Enter to continue)", m.inputs[stepDesc].View())
	case stepDone:
		return fmt.Sprintf("Done!\n\nName:\t\t%s\nIP:\t\t%s\nPort:\t\t%s\nKeyPath:\t%s\nDesc:\t\t%s\n\npress any key to exit.",
			m.name, m.ip, m.port, m.keypath, m.desc)
	default:
		return "Unkown Step"
	}
}
