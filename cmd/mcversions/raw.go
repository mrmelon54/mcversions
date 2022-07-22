package main

import (
	"code.mrmelon54.xyz/sean/go-mcversions"
	"code.mrmelon54.xyz/sean/go-mcversions/structure"
	"code.mrmelon54.xyz/sean/go-mcversions/utils"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"os"
)

func rawMode(f cliFlags) {
	switch {
	case f.listAction != "":
		listAction(f)
	case f.infoAction != "":
		infoAction(f)
	case f.dlAction != "":
		dlAction(f)
	}
}

func closeForJson(f cliFlags) {
	if f.jsonOutput {
		os.Exit(1)
	}
}

func listAction(f cliFlags) {
	if f.listAction == "" {
		closeForJson(f)
		fmt.Println("Set a pattern to find matching versions")
		fmt.Println("  mcversions -list ~1.18")
		fmt.Println("  mcversions -list ~1.16.3")
		return
	}

	con, err := semver.NewConstraint(f.listAction)
	if err != nil {
		closeForJson(f)
		fmt.Printf("Invalid constraint string: %s\n", f.listAction)
		return
	}

	mcv, err := mcversions.NewMCVersions()
	if err != nil {
		closeForJson(f)
		fmt.Printf("Failed to load Minecraft versions: %s\n", err)
		return
	}
	err = mcv.Grab()
	if err != nil {
		closeForJson(f)
		fmt.Printf("Failed to get version metadata: %s\n", err)
		return
	}
	if !f.jsonOutput {
		fmt.Printf("Minecraft versions list:\n")
	}
	versions, err := mcv.ListVersions()
	if err != nil {
		closeForJson(f)
		fmt.Printf("Error: %s\n", err)
		return
	}
	var a []*structure.PistonMetaVersionData
	for i := 0; i < len(versions); i++ {
		if f.listAction == "all" || versions[i].Type == f.listAction || structure.PistonMetaIdCheckConstraints(versions[i].ID, con) {
			a = append(a, versions[i])
		}
	}
	if f.jsonOutput {
		err = json.NewEncoder(os.Stdout).Encode(a)
		if err != nil {
			closeForJson(f)
		}
	} else {
		for _, i := range a {
			fmt.Printf(" - %s %s\n", i.Type, i.ID)
		}
	}
}

func infoAction(f cliFlags) {
	if f.infoAction == "" {
		closeForJson(f)
		fmt.Println("Set a version ID")
		fmt.Println("  mcversions -info 1.18.2")
		fmt.Println("  mcversions -info 1.16.5")
		return
	}

	mcv, err := mcversions.NewMCVersions()
	if err != nil {
		closeForJson(f)
		fmt.Printf("Failed to load Minecraft versions: %s\n", err)
		return
	}
	err = mcv.Grab()
	if err != nil {
		closeForJson(f)
		fmt.Printf("Failed to get version metadata: %s\n", err)
		return
	}
	fmt.Printf("Minecraft version info:\n")
	ver, err := structure.NewPistonMetaId(f.infoAction)
	if err != nil {
		closeForJson(f)
		fmt.Printf("Invalid version code: %s\n", err)
		return
	}
	version, err := mcv.GetVersion(ver)
	if err != nil {
		closeForJson(f)
		fmt.Printf("Error: %s\n", err)
		return
	}
	if f.jsonOutput {
		err = json.NewEncoder(os.Stdout).Encode(version)
		if err != nil {
			closeForJson(f)
		}
	} else {
		fmt.Printf("  ID: %s\n", version.ID)
		fmt.Printf("  Type: %s\n", version.Type)
		fmt.Printf("  URL: %s\n", version.URL)
		fmt.Printf("  Time: %s\n", version.Time)
		fmt.Printf("  Release Time: %s\n", version.ReleaseTime)
		fmt.Printf("  SHA1: %s\n", version.Sha1)
		fmt.Printf("  Compliance Level: %d\n", version.ComplianceLevel)
		fmt.Printf("\nDownload:\n")
		fmt.Printf("  Client: mcversions -dl %s -client\n", version.ID)
		fmt.Printf("  Client Mappings: mcversions -dl %s -client-mappings\n", version.ID)
		fmt.Printf("  Server: mcversions -dl %s -server\n", version.ID)
		fmt.Printf("  Server Mappings: mcversions -dl %s -server-mappings\n", version.ID)
	}
}

func dlAction(f cliFlags) {
	var err error
	var ver *structure.PistonMetaId
	switch f.dlAction {
	case "release":
		version, err := mcversions.LatestRelease()
		if err != nil {
			closeForJson(f)
			fmt.Println("Failed to get latest release metadata")
			return
		}
		ver = version.ID
	case "snapshot":
		version, err := mcversions.LatestSnapshot()
		if err != nil {
			closeForJson(f)
			fmt.Println("Failed to get latest snapshot metadata")
			return
		}
		ver = version.ID
	default:
		ver, err = structure.NewPistonMetaId(f.dlAction)
		if err != nil {
			closeForJson(f)
			fmt.Printf("Invalid version code: %s\n", err)
			return
		}
	}

	if ver == nil {
		closeForJson(f)
		fmt.Printf("Failed to load version data for %s\n", f.dlAction)
		return
	}

	meta, err := mcversions.VersionPackage(ver)
	if err != nil {
		closeForJson(f)
		fmt.Printf("Failed to get piston meta: %s\n", f.dlAction)
		return
	}

	if meta == nil {
		closeForJson(f)
		fmt.Printf("Failed to load download data for %s\n", f.dlAction)
		return
	}

	if f.jsonOutput {
		_ = json.NewEncoder(os.Stdout).Encode(meta)
		closeForJson(f)
	}

	switch {
	case f.dlClient:
		_, err = utils.DownloadJar(meta.ID, *meta.Downloads.Client)
	case f.dlClientMappings:
		_, err = utils.DownloadJar(meta.ID, *meta.Downloads.ClientMappings)
	case f.dlServer:
		_, err = utils.DownloadJar(meta.ID, *meta.Downloads.Server)
	case f.dlServerMappings:
		_, err = utils.DownloadJar(meta.ID, *meta.Downloads.ServerMappings)
	default:
		err = fmt.Errorf("no download options selected: -client, -client-mappings, -server, -server-mappings")
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
