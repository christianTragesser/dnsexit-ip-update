package main

import (
	"fmt"
	"os"

	"github.com/christiantragesser/dnsexit-ip-update/dnsexit"
)

var version = "dev-build"

func main() {
	if len(os.Args) < 2 {
		dnsexit.CLI()
	} else {
		action := os.Args[1]
		switch action {
		case "version", "-v":
			fmt.Printf("DNSExit dynamic IP client version: %s\n", version)

		default:
			dnsexit.CLI()
		}
	}
}
