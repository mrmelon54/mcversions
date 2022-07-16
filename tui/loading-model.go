package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	crossIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Render("✘")
	tickIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render("✔")
)

type LoadingModel struct {
	spinner spinner.Model
	options outputLoadingOptions
	state   loadingState
	err     error
}

type outputLoadingOptions struct {
	failMsg string
	doneMsg string
	loadMsg string
}

// loadingState is an enum to store the state of the loading model
type loadingState int

const (
	loadingStateUnknown loadingState = iota
	loadingStateFail
	loadingStateDone
	loadingStateWait
)

type loadingStateMsg struct {
	state loadingState
	err   error
}

// =====

func NewLoadingModel(options outputLoadingOptions) LoadingModel {
	var l LoadingModel
	l.options = options

	l.spinner = spinner.New()
	l.spinner.Spinner = spinner.Dot
	l.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return l
}

func (m LoadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoadingModel) Update(msg tea.Msg) (LoadingModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case loadingStateMsg:
		m.state = msg.state
		m.err = msg.err
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m LoadingModel) View() (s string) {
	switch m.state {
	case loadingStateFail:
		s += fmt.Sprintf("  %s %s", crossIcon, m.options.failMsg)
	case loadingStateDone:
		s += fmt.Sprintf("  %s %s", tickIcon, m.options.doneMsg)
	case loadingStateWait:
		s += fmt.Sprintf("  %s %s", m.spinner.View(), m.options.loadMsg)
	}
	if m.err != nil {
		s += "\n    Error: " + m.err.Error()
	}
	return
}
