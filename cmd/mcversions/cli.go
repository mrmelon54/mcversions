package main

import (
	"flag"
	"fmt"
	"github.com/MrMelon54/mcversions/v3/cmd/tui"
)

type cliFlags struct {
	jsonOutput       bool
	listAction       string
	infoAction       string
	dlAction         string
	dlClient         bool
	dlClientMappings bool
	dlServer         bool
	dlServerMappings bool
}

func main() {
	var f cliFlags
	flag.BoolVar(&f.jsonOutput, "j", false, "Outputs raw json data")
	flag.StringVar(&f.listAction, "list", "", "List action: -list ~1.18")
	flag.StringVar(&f.infoAction, "info", "", "Info action: -info 1.16.5")
	flag.StringVar(&f.dlAction, "dl", "", "Download action: -dl 1.12.2 [-client] [-server]")
	flag.BoolVar(&f.dlClient, "client", false, "Download client: -dl 1.10 -client")
	flag.BoolVar(&f.dlClientMappings, "client-mappings", false, "Download client: -dl 1.10 -client-mappings")
	flag.BoolVar(&f.dlServer, "server", false, "Download server: -dl 1.10 -server")
	flag.BoolVar(&f.dlServerMappings, "server-mappings", false, "Download server: -dl 1.10 -server-mappings")
	flag.Parse()

	if f != (cliFlags{}) {
		if !singleString([]string{f.listAction, f.infoAction, f.dlAction}) {
			fmt.Println("Only one action can be set at a time!")
			return
		}

		rawMode(f)
		return
	}

	tui.Launch()
}

func singleString(arr []string) bool {
	count := 0
	for _, i := range arr {
		if i != "" {
			count++
		}
	}
	return count == 1
}
