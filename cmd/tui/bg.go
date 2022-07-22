package tui

import (
	"code.mrmelon54.xyz/sean/go-mcversions"
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"code.mrmelon54.xyz/sean/go-mcversions/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type Background struct {
	program       *tea.Program
	mcv           *mcversions.MCVersions
	loadedPackage *structure.PistonMetaPackage
}

func NewBackground() *Background {
	return new(Background)
}

func (b *Background) getMCV() bool {
	if b.mcv == nil {
		var err error
		b.mcv, err = mcversions.NewMCVersions()
		if err != nil {
			b.program.Send(errMsg{err})
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
		b.program.Send(loadingStateLocalMsg{loadingStateFail, err})
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
		b.program.Send(loadingStateRemoteMsg{loadingStateFail, err})
		return
	}
	b.program.Send(loadingStateRemoteMsg{loadingStateDone, err})
}

func (b *Background) fetchPackageMeta(id *structure.PistonMetaId) {
	if !b.getMCV() {
		return
	}

	pack, err := b.mcv.GetVersionPackage(id)
	if err != nil {
		b.program.Send(setPackageMsg{loadingStateMsg{loadingStateFail, err}, nil})
		return
	}
	b.loadedPackage = pack
	b.program.Send(setPackageMsg{loadingStateMsg{loadingStateDone, nil}, pack})
}

func (b *Background) asyncDownloadJar(id *structure.PistonMetaId, data structure.PistonMetaPackageDownloadsData) {
	_, err := utils.DownloadJar(id, data)
	if err != nil {
		b.program.Send(loadingDownloadMsg{loadingStateMsg{loadingStateFail, err}, 1, 1})
	} else {
		b.program.Send(loadingDownloadMsg{loadingStateMsg{loadingStateDone, nil}, 1, 1})
	}
}

func (b *Background) asyncDownloadAll(id *structure.PistonMetaId, dl []*DownloadInfo) {
	t := len(dl)
	for i := range dl {
		b.program.Send(loadingDownloadMsg{loadingStateMsg{loadingStateWait, nil}, i + 1, t})
		_, err := utils.DownloadJar(id, *dl[i].data)
		if err != nil {
			b.program.Send(loadingDownloadMsg{loadingStateMsg{loadingStateFail, err}, i + 1, t})
			return
		}
	}
	b.program.Send(loadingDownloadMsg{loadingStateMsg{loadingStateDone, nil}, t, t})
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

func (b *Background) getLoadedPackage() *structure.PistonMetaPackage {
	return b.loadedPackage
}
