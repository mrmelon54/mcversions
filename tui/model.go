package tui

import (
	"code.mrmelon54.xyz/sean/go-mcversions"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadingState int

const (
	loadingStateFail loadingState = iota
	loadingStateDone
	loadingStateWait
)

type loadingStateLocalMsg loadingState
type loadingStateRemoteMsg loadingState

var (
	crossIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Render("✘")
	tickIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render("✔")
)

type model struct {
	bg          *modelBg
	spinner     spinner.Model
	err         error
	localState  loadingState
	remoteState loadingState
	loadLocal   chan struct{}
	loadRemote  chan struct{}
}

type modelBg struct {
	mcv *mcversions.MCVersions
}

func WelcomeInitialModel() tea.Model {
	var m model
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return m
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.bg == nil {
		m.bg = new(modelBg)
		m.bg.mcv, m.err = mcversions.NewMCVersions()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() (s string) {
	s += "\n"
	if m.err != nil {
		s += fmt.Sprintf("Oh no:\n%s\n\n", m.err.Error())
		return
	}
	s += outputLoadingState(m.localState, outputLoadingOptions{
		"Local cache invalid or outdated",
		"Loaded local cache",
		"Finding local cache",
		crossIcon,
		tickIcon,
		m.spinner.View(),
	})
	if m.localState == loadingStateWait {
		s += "\n\n"
		return
	}
	s += "\n"
	s += outputLoadingState(m.localState, outputLoadingOptions{
		"Failed to load launcher meta",
		"Loaded launcher meta",
		"Finding launcher meta",
		crossIcon,
		tickIcon,
		m.spinner.View(),
	})
	if m.remoteState == loadingStateWait {
		s += "\n\n"
		return
	}
	s += "\n"
	s += "  Choose action:"
	s += "\n\n"
	return
}

type outputLoadingOptions struct {
	failMsg  string
	doneMsg  string
	loadMsg  string
	failIcon string
	doneIcon string
	loadIcon string
}

func outputLoadingState(l loadingState, options outputLoadingOptions) string {
	switch l {
	case loadingStateFail:
		return fmt.Sprintf("  %s %s", options.failIcon, options.failMsg)
	case loadingStateDone:
		return fmt.Sprintf("  %s %s", options.doneIcon, options.doneMsg)
	case loadingStateWait:
		return fmt.Sprintf("  %s %s", options.loadIcon, options.loadMsg)
	}
	return ""
}

func fetchLocalCache(m model) tea.Cmd {
	return func() tea.Msg {
		return loadingStateLocalMsg(loadingStateWait)
	}
}

func fetchRemoteData(m model) tea.Cmd {
	return func() tea.Msg {
		return loadingStateRemoteMsg(loadingStateWait)
	}
}
