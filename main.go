package main

import (
	"os"

	"github.com/christiantragesser/dnsexit-ip-update/dnsexit"
)

func main() {
	// if CLI arguments are provided
	if len(os.Args) > 1 {
		dnsexit.CLIArgs()
	}
}
