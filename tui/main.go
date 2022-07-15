package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func Launch() {
	bg := NewBackground()
	m := InitialModel(bg)
	p := tea.NewProgram(m)
	bg.program = p
	if err := p.Start(); err != nil {
		fmt.Printf("Oh no the program crashed!\n\n")
		fmt.Println(err)
		os.Exit(1)
	}
}
