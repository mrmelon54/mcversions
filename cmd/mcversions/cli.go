package main

import (
	"code.mrmelon54.xyz/sean/go-mcversions/tui"
	"flag"
)

type cliFlags struct {
	rawMode    bool
	jsonOutput bool
	pattern    string
	listAction bool
	dlAction   bool
	dlClient   bool
	dlServer   bool
}

func main() {
	var f cliFlags
	flag.BoolVar(&f.rawMode, "r", false, "Raw mode - without interactive mode")
	flag.BoolVar(&f.jsonOutput, "j", false, "Outputs raw json data")
	flag.StringVar(&f.pattern, "v", "", "Specify a version of pattern to match")
	flag.BoolVar(&f.listAction, "list", false, "List action")
	flag.BoolVar(&f.dlAction, "dl", false, "Download action")
	flag.BoolVar(&f.dlClient, "client", false, "Interact with client")
	flag.BoolVar(&f.dlServer, "server", false, "Interact with server")
	flag.Parse()

	if f.rawMode {
		rawMode(f)
		return
	}
	if f.jsonOutput {
		jsonMode(f)
		return
	}

	tui.Launch()
}
