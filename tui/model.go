package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	bg         *Background
	appState   int
	err        error
	loadLocal  LoadingModel
	loadRemote LoadingModel
	vSelect    VersionSelectModel
}

type errMsg error
type loadingStateLocalMsg loadingStateMsg
type loadingStateRemoteMsg loadingStateMsg

func InitialModel(bg *Background) Model {
	var m Model
	m.bg = bg
	m.loadLocal = NewLoadingModel(outputLoadingOptions{
		"Local cache invalid or outdated",
		"Loaded local cache",
		"Finding local cache",
	})
	m.loadRemote = NewLoadingModel(outputLoadingOptions{
		"Failed to load launcher meta",
		"Loaded launcher meta",
		"Finding launcher meta",
	})
	m.vSelect = NewVersionSelectModel(bg)
	return m
}

func (m Model) Init() tea.Cmd {
	go m.bg.fetchLocalData()
	return tea.Batch(m.loadLocal.Init(), m.loadRemote.Init(), m.vSelect.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}
	case loadingStateLocalMsg:
		m.loadLocal, cmd = m.loadLocal.Update(loadingStateMsg(msg))
		cmds = append(cmds, cmd)
		if m.appState == 0 {
			switch m.loadLocal.state {
			case loadingStateDone:
				m.appState = 2
			case loadingStateFail:
				m.loadRemote.state = loadingStateWait
				go m.bg.fetchRemoteData()
				m.appState = 1
			}
		}
	case loadingStateRemoteMsg:
		m.loadRemote, cmd = m.loadRemote.Update(loadingStateMsg(msg))
		cmds = append(cmds, cmd)
		if m.loadRemote.state == loadingStateDone && m.appState == 1 {
			m.appState = 2
		}
	case errMsg:
		m.err = msg
	case spinner.TickMsg:
		m.loadLocal, cmd = m.loadLocal.Update(msg)
		cmds = append(cmds, cmd)
		m.loadRemote, cmd = m.loadRemote.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.appState == 2 {
		m.vSelect, cmd = m.vSelect.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() (s string) {
	s += "\n"
	if m.err != nil {
		s += fmt.Sprintf("Oh no:\n%s\n\n", m.err.Error())
		return
	}
	s += m.loadLocal.View()
	if m.appState > 0 && m.loadRemote.state != loadingStateUnknown {
		s += "\n"
		s += m.loadRemote.View()
	}
	if m.appState > 1 {
		s += "\n\n"
		s += m.vSelect.View()
	}
	s += "\n\n"
	return
}
