package main

import (
	"code.mrmelon54.xyz/sean/go-mcversions/tui"
	"flag"
	"fmt"
)

type cliFlags struct {
	jsonOutput bool
	pattern    string
	listAction bool
	infoAction bool
	dlAction   bool
	dlClient   bool
	dlServer   bool
}

func main() {
	var f cliFlags
	flag.BoolVar(&f.jsonOutput, "j", false, "Outputs raw json data")
	flag.StringVar(&f.pattern, "v", "", "Specify a version of pattern to match")
	flag.BoolVar(&f.listAction, "list", false, "List action")
	flag.BoolVar(&f.infoAction, "info", false, "Info action")
	flag.BoolVar(&f.dlAction, "dl", false, "Download action")
	flag.BoolVar(&f.dlClient, "client", false, "Interact with client")
	flag.BoolVar(&f.dlServer, "server", false, "Interact with server")
	flag.Parse()

	if f != (cliFlags{}) {
		if !singleBool([]bool{f.listAction, f.infoAction, f.dlAction}) {
			fmt.Println("Only one action can be set at a time!")
			return
		}

		if f.jsonOutput {
			jsonMode(f)
			return
		}
		rawMode(f)
		return
	}

	tui.Launch()
}

func singleBool(arr []bool) bool {
	count := 0
	for _, i := range arr {
		if i {
			count++
		}
	}
	return count == 1
}
