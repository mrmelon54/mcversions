package main

import (
	"code.mrmelon54.xyz/sean/go-mcversions"
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func rawMode(f cliFlags) {
	switch {
	case f.listAction:
		listAction(f)
	case f.dlAction:
		dlAction(f)
	}
}

func listAction(f cliFlags) {
	reg := "^" + strings.ReplaceAll(regexp.QuoteMeta(os.Args[2]), "\\*", ".*?") + "$"

	mcv, err := mcversions.NewMCVersions()
	if err != nil {
		fmt.Printf("Failed to load Minecraft versions: %s\n", err)
		return
	}
	fmt.Printf("Minecraft versions list:\n")
	versions, err := mcv.ListVersions()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	for i := 0; i < len(versions); i++ {
		matched, _ := regexp.MatchString(reg, versions[i].ID)
		if f.pattern == "all" || versions[i].Type == f.pattern || versions[i].ID == f.pattern || matched {
			fmt.Printf(" - %s %s\n", versions[i].Type, versions[i].ID)
		}
	}
}

func dlAction(f cliFlags) {
	var version *structure.APIVersionData
	var err error
	switch f.pattern {
	case "release":
		version, err = mcversions.LatestRelease()
		if err != nil {
			fmt.Println("Failed to get latest release metadata")
			return
		}
	case "snapshot":
		version, err = mcversions.LatestSnapshot()
		if err != nil {
			fmt.Println("Failed to get latest snapshot metadata")
			return
		}
	default:
		version, err = mcversions.Version(f.pattern)
		if err != nil {
			fmt.Printf("Failed to get version: '%s'\n", f.pattern)
		}
	}

	meta, err := mcversions.NewPistonMeta(version.URL)
	if err != nil {
		fmt.Printf("Failed to get piston meta: '%s' (%s)\n", f.pattern, version.ID)
	}

	switch {
	case f.dlClient:
		downloadJar(version.ID, *meta.GetClient())
	case f.dlServer:
		downloadJar(version.ID, *meta.GetServer())
	}
}

func downloadJar(id string, dd structure.APIDownloadData) int64 {
	filename := id + "-" + path.Base(dd.URL)
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		fmt.Printf("Error: file already exists\n")
		return 0
	}
	out, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file\n")
		return 0
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)
	resp, err := http.Get(dd.URL)
	if err != nil {
		fmt.Printf("Error starting download\n")
		return 0
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	h := sha1.New()

	// Connect 'out' and 'h' as a single writer
	w := io.MultiWriter(out, h)

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		fmt.Printf("Error during download\n")
		return 0
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	if n == dd.Size {
		fmt.Printf("Download size matches\n")
	} else {
		fmt.Printf("Incorrect download size\n")
		return 0
	}

	sha1str := h.Sum(nil)
	if hex.EncodeToString(sha1str) == dd.Sha1 {
		fmt.Printf("Sha1 hashes match so the download is probably safe\n")
	} else {
		fmt.Printf("Sha1 hashes don't match... deleting it for your safety\n")
		err = os.Remove(filename)
		if err != nil {
			fmt.Println("Failed to remove the unsafe file:", err)
			return 0
		}
		return 0
	}
	return n
}

/*

	// Details options
	if len(os.Args) == 2 {
		if os.Args[1] == "list" {
			fmt.Printf("Usage 'mcversions list <all/release/snapshot/old_alpha/old_beta/pattern>'\n")
			return
		}

		mcv, err := mcversions.NewMCVersions()
		if err != nil {
			fmt.Printf("Failed to load Minecraft versions\n")
			return
		}
		versionid := os.Args[1]
		if os.Args[1] == "release" {
			versionid = mcv.GetLatestRelease()
		} else if os.Args[1] == "snapshot" {
			versionid = mcv.GetLatestSnapshot()
		}
		v, err := mcv.GetVersion(versionid)
		if err != nil {
			fmt.Printf("Failed to get version information\n")
			return
		}
		fmt.Printf("ID: %s\n", v.GetID())
		fmt.Printf("Type: %s\n", v.GetType())
		fmt.Printf("Release time: %s\n", v.GetReleaseTime())
		fmt.Printf("Client:\n")
		fmt.Printf(" - URL: %s\n", v.GetClient().URL)
		fmt.Printf(" - Sha1: %s\n", v.GetClient().Sha1)
		fmt.Printf(" - Size: %v\n", v.GetClient().Size)
		fmt.Printf("Server:\n")
		fmt.Printf(" - URL: %s\n", v.GetServer().URL)
		fmt.Printf(" - Sha1: %s\n", v.GetServer().Sha1)
		fmt.Printf(" - Size: %v\n", v.GetServer().Size)
		return
	}

	// Help options
	if len(os.Args) == 1 {
		fmt.Printf("mcversions list <all/release/snapshot/old_alpha/old_beta/pattern> - List all versions of the specified type\n")
		fmt.Printf("mcversions <version id/release/snapshot> - Get details about the version\n")
		fmt.Printf("mcversions <version id/release/snapshot> <client/server> - Download the client/server jar\n")
		return
	}
}
*/
