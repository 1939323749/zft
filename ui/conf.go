package ui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"github.com/rotisserie/eris"
	"os"
	"strings"
	"zft/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Save ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Save"))
)

type settingModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func initialModel() settingModel {
	m := settingModel{
		inputs: make([]textinput.Model, 4),
	}
	conf, err := utils.GetConf()
	if err != nil {
		Errors = eris.New(err.Error())
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			if conf.Operator != "" {
				t.Placeholder = fmt.Sprintf("%-15s%s", "operator:", conf.Operator)
			} else {
				t.Placeholder = "operator"
			}

			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			if conf.Secret != "" {
				t.Placeholder = fmt.Sprintf("%-15s%s", "secret:", conf.Secret)
			} else {
				t.Placeholder = "secret"
			}

			t.CharLimit = 64
		case 2:
			if conf.Bucket != "" {
				t.Placeholder = fmt.Sprintf("%-15s%s", "bucket:", conf.Bucket)
			} else {
				t.Placeholder = "bucket"
			}
			t.CharLimit = 64
		case 3:
			if conf.Bucketurl != "" {
				t.Placeholder = fmt.Sprintf("%-15s%s", "bucket url:", conf.Bucketurl)
			} else {
				t.Placeholder = "bucket url"
			}
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}

func (m settingModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m settingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			conf, err := utils.GetConf()
			if err != nil {
				Errors = eris.New(err.Error())
			}
			// Did the user press enter while the Save button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				upyun := utils.Upyun{}

				if value := m.inputs[0].Value(); value != "" {
					upyun.Operator = value
				} else {
					upyun.Operator = conf.Operator
				}
				if value := m.inputs[1].Value(); value != "" {
					upyun.Secret = value
				} else {
					upyun.Secret = conf.Secret
				}
				if value := m.inputs[2].Value(); value != "" {
					upyun.Bucket = value
				} else {
					upyun.Bucket = conf.Bucket
				}
				if value := m.inputs[3].Value(); value != "" {
					Errors = eris.New(value)
					upyun.Bucketurl = value
				} else {
					upyun.Bucketurl = conf.Bucketurl
				}

				err := utils.SetConf(upyun)
				if err != nil {
					Errors = fmt.Errorf("failed to set configuration: %v", err)
				}
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
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
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *settingModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m settingModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func Conf() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
