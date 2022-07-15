package tui

import (
	"code.mrmelon54.xyz/sean/go-mcversions"
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	tea "github.com/charmbracelet/bubbletea"
)

type Background struct {
	program *tea.Program
	mcv     *mcversions.MCVersions
}

func NewBackground() *Background {
	return new(Background)
}

func (b *Background) getMCV() bool {
	if b.mcv == nil {
		var err error
		b.mcv, err = mcversions.NewMCVersions()
		if err != nil {
			b.program.Send(errMsg(err))
			return false
		}
	}
	return true
}

func (b *Background) fetchLocalData() {
	if !b.getMCV() {
		return
	}

	err := b.mcv.Load()
	if err != nil {
		if err == mcversions.ErrCacheMissing || err == mcversions.ErrCacheExpired {
			b.program.Send(loadingStateLocalMsg{loadingStateFail, err})
			return
		}
		b.program.Send(loadingStateLocalMsg{loadingStateFail, err})
		b.program.Send(errMsg(err))
		return
	}
	b.program.Send(loadingStateLocalMsg{loadingStateDone, err})
}

func (b *Background) fetchRemoteData() {
	if !b.getMCV() {
		return
	}

	err := b.mcv.Fetch()
	if err != nil {
		if err == mcversions.ErrCacheMissing || err == mcversions.ErrCacheExpired {
			b.program.Send(loadingStateRemoteMsg{loadingStateFail, err})
			return
		}
		b.program.Send(loadingStateRemoteMsg{loadingStateFail, err})
		b.program.Send(errMsg(err))
		return
	}
	b.program.Send(loadingStateRemoteMsg{loadingStateDone, err})
}

func (b *Background) getLatestRelease() *structure.PistonMetaVersionData {
	if !b.getMCV() {
		return nil
	}

	release, err := b.mcv.LatestRelease()
	if err != nil {
		return nil
	}
	return release
}

func (b *Background) getLatestSnapshot() *structure.PistonMetaVersionData {
	if !b.getMCV() {
		return nil
	}

	snapshot, err := b.mcv.LatestSnapshot()
	if err != nil {
		return nil
	}
	return snapshot
}

func (b *Background) getAllVersions() []*structure.PistonMetaVersionData {
	if !b.getMCV() {
		return nil
	}

	v, err := b.mcv.ListVersions()
	if err != nil {
		return nil
	}
	return v
}
