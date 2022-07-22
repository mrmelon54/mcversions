package tui

import (
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

var (
	arrowDlIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaa00")).Bold(true).Render("▸")
	footerText  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	footerDot   = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(" • ")
)

type Model struct {
	bg         *Background
	appState   int
	err        error
	loadLocal  LoadingModel
	loadRemote LoadingModel
	vSelect    VersionSelectModel
	loadPack   LoadingModel
	version    *structure.PistonMetaVersionData
	pack       *structure.PistonMetaPackage
	dl         []*DownloadInfo
	dlChoice   int
	loadDl     LoadingModel
	dlNow      int
	dlTotal    int
}

type errMsg struct {
	err error
}
type loadingStateLocalMsg loadingStateMsg
type loadingStateRemoteMsg loadingStateMsg
type loadingDownloadMsg struct {
	state loadingStateMsg
	now   int
	total int
}
type setVersionMsg struct {
	version *structure.PistonMetaVersionData
}
type setPackageMsg struct {
	state loadingStateMsg
	pack  *structure.PistonMetaPackage
}

type DownloadInfo struct {
	id   string
	name string
	data *structure.PistonMetaPackageDownloadsData
}

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
	m.loadPack = NewLoadingModel(outputLoadingOptions{
		failMsg: "Failed to load package meta",
		doneMsg: "Loaded package meta",
		loadMsg: "Finding package meta",
	})
	m.loadDl = NewLoadingModel(outputLoadingOptions{
		failMsg: "Failed to download file",
		doneMsg: "Downloaded file successfully",
		loadMsg: "Downloading...",
	})
	return m
}

func (m Model) Init() tea.Cmd {
	go m.bg.fetchLocalData()
	return tea.Batch(m.loadLocal.Init(), m.loadRemote.Init(), m.vSelect.Init(), m.loadPack.Init(), m.loadDl.Init())
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
	case loadingDownloadMsg:
		if msg.total == 1 {
			m.loadDl.options = outputLoadingOptions{
				failMsg: "Failed to download file",
				doneMsg: "Downloaded file successfully",
				loadMsg: "Downloading...",
			}
		} else {
			m.loadDl.options = outputLoadingOptions{
				failMsg: fmt.Sprintf("Failed to download file: %d / %d", msg.now, msg.total),
				doneMsg: fmt.Sprintf("Downloaded files successfully: %d / %d", msg.now, msg.total),
				loadMsg: fmt.Sprintf("Downloading... %d / %d", msg.now, msg.total),
			}
		}
		m.loadDl, cmd = m.loadDl.Update(msg.state)
		cmds = append(cmds, cmd)
		if m.loadDl.state == loadingStateDone && m.appState == 5 {
			cmds = append(cmds, tea.Quit)
		}
	case spinner.TickMsg:
		m.loadLocal, cmd = m.loadLocal.Update(msg)
		cmds = append(cmds, cmd)
		m.loadRemote, cmd = m.loadRemote.Update(msg)
		cmds = append(cmds, cmd)
		m.loadPack, cmd = m.loadPack.Update(msg)
		cmds = append(cmds, cmd)
		m.loadDl, cmd = m.loadDl.Update(msg)
		cmds = append(cmds, cmd)
	case setVersionMsg:
		m.version = msg.version
		if m.appState == 2 {
			m.appState = 3
		}
		go m.bg.fetchPackageMeta(m.version.ID)
	case setPackageMsg:
		m.loadPack, cmd = m.loadPack.Update(msg.state)
		cmds = append(cmds, cmd)
		m.pack = msg.pack
		m.dl = make([]*DownloadInfo, 0)
		if m.pack.Downloads.Client != nil {
			m.dl = append(m.dl, &DownloadInfo{"client", "Client", m.pack.Downloads.Client})
		}
		if m.pack.Downloads.ClientMappings != nil {
			m.dl = append(m.dl, &DownloadInfo{"client-mappings", "Client Mappings", m.pack.Downloads.ClientMappings})
		}
		if m.pack.Downloads.Server != nil {
			m.dl = append(m.dl, &DownloadInfo{"server", "Server", m.pack.Downloads.Server})
		}
		if m.pack.Downloads.ServerMappings != nil {
			m.dl = append(m.dl, &DownloadInfo{"server-mappings", "Server Mappings", m.pack.Downloads.ServerMappings})
		}

		// Next screen
		if m.loadPack.state == loadingStateDone && (m.appState == 2 || m.appState == 3) {
			m.appState = 4
		}
	case errMsg:
		m.err = msg.err
	}
	switch m.appState {
	case 2:
		m.vSelect, cmd = m.vSelect.Update(msg)
		cmds = append(cmds, cmd)
	case 4:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.pack != nil && m.dl != nil && len(m.dl) > 0 {
					m.appState = 5
					m.loadDl.state = loadingStateWait
					if m.dlChoice == len(m.dl) {
						go m.bg.asyncDownloadAll(m.pack.ID, m.dl)
					} else {
						go m.bg.asyncDownloadJar(m.pack.ID, *m.dl[m.dlChoice].data)
					}
				}
			case tea.KeyUp:
				m.dlChoice--
				if m.dlChoice < 0 {
					m.dlChoice = 0
				}
			case tea.KeyDown:
				m.dlChoice++
				if m.dlChoice > len(m.dl) {
					m.dlChoice = len(m.dl)
				}
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() (s string) {
	s = "\n"
	if m.err != nil {
		s += fmt.Sprintf("Oh no:\n%s\n\n", m.err)
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
		if m.appState == 2 {
			s += "\n\n"
			switch m.vSelect.action {
			case 0:
				s += "  " + footerText.Render("up/down: select") + footerDot + footerText.Render("enter: choose") + footerDot + footerText.Render("esc: quit")
			case 1:
				s += "  " + footerText.Render("enter: choose") + footerDot + footerText.Render("esc: quit")
			}
		}
	}
	if m.appState > 2 {
		s += "\n\n"
		s += m.loadPack.View()
	}
	if m.appState > 3 {
		s += "\n\n"
		if m.pack == nil {
			s += "  Failed to load version package"
		} else {
			version := m.pack
			s += fmt.Sprintf("  ID: %s\n", version.ID)
			s += fmt.Sprintf("  Type: %s\n", version.Type)
			s += fmt.Sprintf("  Time: %s\n", version.Time)
			s += fmt.Sprintf("  Release Time: %s\n", version.ReleaseTime)
			s += fmt.Sprintf("  Main Class: %s\n", version.MainClass)
			s += fmt.Sprintf("  Minimum Launcher Version: %d\n", version.MinimumLauncherVersion)
			s += fmt.Sprintf("  Compliance Level: %d\n", version.ComplianceLevel)

			s += "\n"
			s += fmt.Sprintf("  Downloads:")
			t := uint64(0)
			for i := range m.dl {
				a := " "
				if i == m.dlChoice {
					a = arrowDlIcon
				}
				t += uint64(m.dl[i].data.Size)
				s += fmt.Sprintf("\n  %s %s (%s)", a, m.dl[i].name, humanize.Bytes(uint64(m.dl[i].data.Size)))
			}
			a := " "
			if len(m.dl) == m.dlChoice {
				a = arrowDlIcon
			}
			s += fmt.Sprintf("\n  %s Download All (%s)", a, humanize.Bytes(t))
			if m.appState == 4 {
				s += "\n\n"
				s += "  " + footerText.Render("up/down: select") + footerDot + footerText.Render("enter: choose") + footerDot + footerText.Render("esc: quit")
			}
		}
	}
	if m.appState > 4 {
		s += "\n\n"
		s += m.loadDl.View()
		if m.appState == 5 {
			s += "\n\n"
			s += "  " + footerText.Render("esc: quit")
		}
	}
	s += "\n\n"
	return
}
