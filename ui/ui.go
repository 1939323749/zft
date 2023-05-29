package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rotisserie/eris"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"zft/utils"
)

var docStyle = lipgloss.NewStyle().Margin(0, 0)
var Errors error

type model struct {
	list    list.Model
	dir     string
	baseDir string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			if len(m.list.FilterInput.Value()) != 0 {
				m.list.FilterInput.SetValue(m.list.FilterInput.Value()[:len(m.list.FilterInput.Value())-1])
			}
			if m.dir == m.baseDir {
				return m, nil
			}
			m.dir = filepath.Dir(m.dir)
			items, err := utils.GetFiles(m.dir)
			if err != nil {
				Errors = eris.New(err.Error())
				return m, tea.Quit
			}
			m.list.SetItems(items)
			m.list.Title = filepath.Base(m.dir)
			return m, nil
		case " ":
			//Errors = eris.New(m.dir + "/" + m.list.SelectedItem().FilterValue())
			//return m, tea.Quit
			if m.list.FilterState().String() == "filtering" {
				m.list.FilterInput.SetValue(m.list.FilterInput.Value() + " ")
				m.list.FilterInput.CursorEnd()
				return m, nil
			}
			//?Cannot be used without adding m.list.FilteringEnabled()
			if runtime.GOOS == "darwin" && (m.list.FilterState().String() == "filter applied" || m.list.FilteringEnabled()) {
				cmd := exec.Command("qlmanage", "-p", m.dir+"/"+m.list.SelectedItem().FilterValue())
				err := cmd.Run()
				if err != nil {
					Errors = eris.New(err.Error())
					return m, tea.Quit
				}
			}
			return m, nil
		case "enter":
			if m.list.FilterState().String() == "filter applied" {
				m.list.ResetFilter()
			}
			if m.list.SelectedItem() == nil {
				Errors = eris.New("Selected item is nil.")
			}
			info, err := os.Stat(filepath.Join(m.dir, m.list.SelectedItem().FilterValue()))

			if err != nil {
				Errors = eris.New(err.Error())
				return m, tea.Quit
			}

			if info.IsDir() {
				newitems, err := utils.GetFiles(m.dir + "/" + m.list.SelectedItem().FilterValue())
				if err != nil {
					Errors = eris.New(err.Error())
					return m, tea.Quit
				}
				m.list.Title = filepath.Base(m.list.SelectedItem().FilterValue())
				m.dir = filepath.Join(m.dir, m.list.SelectedItem().FilterValue())
				m.list.SetItems(newitems)
				return m, nil
			} else {
				relpath, err := filepath.Rel(m.baseDir, m.dir)
				if err != nil {
					Errors = eris.New(err.Error())
					return m, tea.Quit
				}

				if relpath != "." {
					err = utils.UploadFile(relpath + "/" + m.list.SelectedItem().FilterValue())
					if err != nil {
						Errors = eris.New(err.Error())
					}
				} else {
					err = utils.UploadFile(m.list.SelectedItem().FilterValue())
					if err != nil {
						Errors = eris.New(err.Error())
					}
				}
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Run() {

	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		Errors = eris.New(fmt.Errorf("error getting current directory :%w", err).Error())
		os.Exit(1)
	}

	items, _ := utils.GetFiles(dir)

	// Create the list
	m := model{
		list:    list.New(items, list.NewDefaultDelegate(), 0, 0),
		dir:     dir,
		baseDir: dir,
	}
	m.list.Title = filepath.Base(dir)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
