package ui

import (
	"fmt"
	"github.com/1939323749/zft/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rotisserie/eris"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
			// 如果选择的项目为空，则退出
			selectedItem := m.list.SelectedItem()
			if selectedItem == nil {
				Errors = eris.New("Selected item is nil.")
				return m, tea.Quit
			}

			// 获取选择的项目的完整路径
			selectedPath := filepath.Join(m.dir, selectedItem.FilterValue())
			info, err := os.Stat(selectedPath)

			if err != nil {
				Errors = eris.New(err.Error())
				return m, tea.Quit
			}

			// 如果选择的项目是目录
			if info.IsDir() {
				// 如果应用了过滤器，则重置过滤器
				if m.list.FilterState().String() == "filter applied" {
					m.list.ResetFilter()
				}

				newitems, err := utils.GetFiles(selectedPath)
				if err != nil {
					Errors = eris.New(err.Error())
					return m, tea.Quit
				}
				m.list.Title = filepath.Base(selectedItem.FilterValue())
				m.dir = selectedPath
				m.list.SetItems(newitems)
				return m, nil
			} else {
				// 如果选择的项目是文件
				relpath, err := filepath.Rel(m.baseDir, m.dir)
				if err != nil {
					Errors = eris.New(err.Error())
					return m, tea.Quit
				}

				if relpath != "." {
					err = utils.UploadFile(relpath + "/" + selectedItem.FilterValue())
					if err != nil {
						Errors = eris.New(err.Error())
					}
				} else {
					err = utils.UploadFile(selectedItem.FilterValue())
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

	items, err := utils.GetFiles(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
