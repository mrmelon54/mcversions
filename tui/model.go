package tui

import (
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

var (
	crossIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Render("✘")
	tickIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render("✔")
)

type WelcomeModel struct {
	spinner     spinner.Model
	localState  loadingState
	remoteState loadingState
}

func WelcomeInitialModel() (model WelcomeModel) {
	model.spinner = spinner.New()
	model.spinner.Spinner = spinner.Dot
	model.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return
}

func (m WelcomeModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m WelcomeModel) View() (s string) {
	s += "\n\n"
	outputLoadingState(m.localState, outputLoadingOptions{
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
	outputLoadingState(m.localState, outputLoadingOptions{
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
	s += "  Enter version now!!"
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
		return fmt.Sprintf("  %s %s", options.failIcon, options.doneMsg)
	case loadingStateDone:
		return fmt.Sprintf("  %s %s", options.doneIcon, options.doneMsg)
	case loadingStateWait:
		return fmt.Sprintf("  %s %s", options.loadIcon, options.loadMsg)
	}
	return ""
}
