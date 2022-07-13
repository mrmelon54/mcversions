package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func initialModel() tea.Model {
	return WelcomeInitialModel()
}

func Launch() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}
