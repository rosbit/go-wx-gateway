/**
 * main process
 * Usage: wx-gateway[ -v]
 * Rosbit Xu
 */
package main

import (
	"os"
	"fmt"
	"wx-gateway/conf"
)

// variables set via go build -ldflags
var (
	buildTime string
	osInfo    string
	goInfo    string
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		ShowInfo("name",       os.Args[0])
		ShowInfo("build time", buildTime)
		ShowInfo("os name",    osInfo)
		ShowInfo("compiler",   goInfo)
		return
	}

	if err := gwconf.CheckGlobalConf(); err != nil {
		fmt.Printf("Failed to check conf: %v\n", err)
		os.Exit(3)
		return
	}
	gwconf.DumpConf()

	if err := StartWxGateway(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(4)
	}
	os.Exit(0)
}

func ShowInfo(prompt, info string) {
	if info != "" {
		fmt.Printf("%10s: %s\n", prompt, info)
	}
}
