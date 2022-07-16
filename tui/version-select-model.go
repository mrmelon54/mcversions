package tui

import (
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var arrowIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#5555ff")).Bold(true).Render("▸")

type VersionSelectModel struct {
	bg            *Background
	input         textinput.Model
	version       *structure.PistonMetaVersionData
	action        int
	vSelChoice    int
	vChoice       int
	cacheRelease  *structure.PistonMetaVersionData
	cacheSnapshot *structure.PistonMetaVersionData
	acCache       []*structure.PistonMetaVersionData
	acValid       bool
	packageState  loadingState
	packageData   *structure.PistonMetaPackage
}

func NewVersionSelectModel(bg *Background) VersionSelectModel {
	in := textinput.New()
	in.Prompt = ""
	in.Width = 25
	return VersionSelectModel{
		bg:         bg,
		input:      in,
		vSelChoice: -1,
	}
}

func (m VersionSelectModel) Init() tea.Cmd {
	return nil
}

func (m VersionSelectModel) Update(msg tea.Msg) (VersionSelectModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.action {
	case 0: // select release, snapshot or search
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if m.vSelChoice < 0 {
				switch msg.Type {
				case tea.KeyUp:
					m.vChoice--
					if m.vChoice < 0 {
						m.vChoice = 0
					}
				case tea.KeyDown:
					m.vChoice++
					if m.vChoice > 2 {
						m.vChoice = 2
					}
				case tea.KeyEnter:
					m.vSelChoice = m.vChoice
					if m.vSelChoice == 2 {
						m.action = 1
						cmds = append(cmds, textinput.Blink, m.input.Focus())
					} else {
						m.action = 2
						switch m.vSelChoice {
						case 0:
							cmds = append(cmds, func() tea.Msg {
								return setVersionMsg{m.latestRelease()}
							})
						case 1:
							cmds = append(cmds, func() tea.Msg {
								return setVersionMsg{m.latestSnapshot()}
							})
						}
					}
				}
			}
		}
	case 1: // search mode
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter {
				m.action = 2
				m.input.Blur()
				version := m.genMatchingVersion(m.input.Value())
				cmds = append(cmds, func() tea.Msg {
					return setVersionMsg{version}
				})
			} else {
				v := m.input.Value()
				m.input, cmd = m.input.Update(msg)
				v2 := m.input.Value()
				if v != v2 {
					m.acCache, m.acValid = m.genAutocompleteVersions(m.input.Value())
				}
				cmds = append(cmds, cmd)
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m VersionSelectModel) View() (s string) {
	if m.vSelChoice < 0 {
		s += "  Version:\n"
		s += outputChoiceOptions(m.vChoice, []string{
			"Latest release" + m.formatVersionInList(m.latestRelease()),
			"Latest snapshot" + m.formatVersionInList(m.latestSnapshot()),
			"Search for version",
		})
	} else {
		s += "  Version: "
		switch m.vSelChoice {
		case 0:
			s += "Latest release" + m.formatVersionInList(m.latestRelease())
		case 1:
			s += "Latest snapshot" + m.formatVersionInList(m.latestSnapshot())
		case 2:
			if m.action == 1 {
				if m.acValid {
					m.input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#55ff55")).Bold(true)
					s += m.input.View()
				} else {
					m.input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.NoColor{}).Bold(false)
					s += m.input.View()
				}
				for _, i := range m.acCache {
					s += fmt.Sprintf("\n           %s", i.ID)
				}
			} else {
				s += m.input.Value()
			}
		}
	}
	return
}

func outputChoiceOptions(choice int, choices []string) (s string) {
	for i := range choices {
		if i > 0 {
			s += "\n"
		}
		a := " "
		if i == choice {
			a = arrowIcon
		}
		s += fmt.Sprintf("  %s %s", a, choices[i])
	}
	return
}

func (m VersionSelectModel) latestRelease() *structure.PistonMetaVersionData {
	if m.bg == nil {
		return nil
	}
	if m.cacheRelease == nil {
		m.cacheRelease = m.bg.getLatestRelease()
		if m.cacheRelease == nil {
			m.cacheRelease = nil
		}
	}
	return m.cacheRelease
}

func (m VersionSelectModel) latestSnapshot() *structure.PistonMetaVersionData {
	if m.bg == nil {
		return nil
	}
	if m.cacheSnapshot == nil {
		m.cacheSnapshot = m.bg.getLatestSnapshot()
		if m.cacheSnapshot == nil {
			m.cacheSnapshot = nil
		}
	}
	return m.cacheSnapshot
}

func (m VersionSelectModel) formatVersionInList(version *structure.PistonMetaVersionData) string {
	if version == nil {
		return " (unknown)"
	}
	return fmt.Sprintf(" (%s)", version.ID)
}

func (m VersionSelectModel) genAutocompleteVersions(id string) ([]*structure.PistonMetaVersionData, bool) {
	if m.bg == nil {
		return nil, false
	}
	v, err := structure.NewPistonMetaId(id)
	if err != nil {
		return []*structure.PistonMetaVersionData{}, false
	}

	ver := m.bg.getAllVersions()
	var out []*structure.PistonMetaVersionData
	var valid bool
	for _, i := range ver {
		if i.ID.Equal(v) {
			out = append([]*structure.PistonMetaVersionData{i}, out...)
			valid = true
		} else if strings.Contains(i.ID.String(), v.String()) {
			out = append(out, i)
			if len(out) > 5 {
				out = out[:5]
			}
		}
	}
	return out, valid
}

func (m VersionSelectModel) genMatchingVersion(id string) *structure.PistonMetaVersionData {
	if m.bg == nil {
		return nil
	}
	v, err := structure.NewPistonMetaId(id)
	if err != nil {
		return nil
	}

	ver := m.bg.getAllVersions()
	for _, i := range ver {
		if i.ID.Equal(v) {
			return i
		}
	}
	return nil
}
