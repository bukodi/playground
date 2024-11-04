package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	printVersion()

}

// printVersion prints the application version
func printVersion() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Unable to determine version information.")
		return
	}

	if buildInfo.Main.Version != "" {
		fmt.Printf("Version: %s\n", buildInfo.Main.Version)
	} else {
		fmt.Println("Version: unknown")
	}
}
