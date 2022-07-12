package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"code.mrmelon54.xyz/sean/go-mcversions"
)

func main() {
	var interactive bool
	var jsonOutput bool
	var listType string

	flag.BoolVar(&interactive, "i", false, "Run in interactive mode")
	flag.BoolVar(&jsonOutput, "j", false, "Outputs raw json data")
	flag.StringVar(&listType, "l", "", "Specifies the versions to list")

	if !jsonOutput {
		fmt.Printf("Minecraft Version Lookup -- by MrMelon\n")
	}

	// Download options
	if len(os.Args) == 3 {
		// List versions
		if os.Args[1] == "list" {
			reg := "^" + strings.ReplaceAll(regexp.QuoteMeta(os.Args[2]), "\\*", ".*?") + "$"

			mcv, err := mcversions.NewMCVersions()
			if err != nil {
				fmt.Printf("Failed to load Minecraft versions: %s\n", err)
				return
			}
			fmt.Printf("Minecraft versions list:\n")
			versions := mcv.List()
			for i := 0; i < len(versions); i++ {
				matched, _ := regexp.MatchString(reg, versions[i].ID)
				if os.Args[2] == "all" || versions[i].Type == os.Args[2] || matched {
					fmt.Printf(" - %s %s\n", versions[i].Type, versions[i].ID)
				}
			}
			return
		}

		if os.Args[2] != "client" && os.Args[2] != "server" {
			fmt.Printf("Invalid option for download: use 'client' or 'server'.\n")
			return
		}
		mcv, err := mcversions.NewMCVersions()
		if err != nil {
			fmt.Printf("Failed to load Minecraft versions\n")
			return
		}
		versionId := os.Args[1]
		if os.Args[1] == "release" {
			versionId = mcv.GetLatestRelease()
		} else if os.Args[1] == "snapshot" {
			versionId = mcv.GetLatestSnapshot()
		}
		v, err := mcv.Get(versionId)
		if err != nil {
			fmt.Printf("Failed to get version information\n")
			return
		}
		var d mcversions.APIDownloadData
		if os.Args[2] == "client" {
			d = v.GetClient()
		} else {
			d = v.GetServer()
		}
		fmt.Printf("Downloading Minecraft %s:\n", os.Args[2])
		fmt.Printf(" - URL: %s\n", d.URL)
		fmt.Printf(" - Sha1: %s\n", d.Sha1)
		fmt.Printf(" - Size: %v\n", d.Size)
		outputsize := downloadJar(versionId, d)
		if outputsize != 0 {
			fmt.Printf("File saved :)\n")
		}
		return
	}

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
		v, err := mcv.Get(versionid)
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

func downloadJar(id string, dd mcversions.APIDownloadData) int64 {
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
