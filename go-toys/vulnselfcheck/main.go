package main

import (
	"context"
	"fmt"
	"golang.org/x/vuln/scan"
	"os"
)

func main() {
	selfCheck()
	fmt.Printf("Hello, playground! Binary: %s", os.Args[0])

}

func selfCheck() {
	ctx := context.Background()
	cmd := scan.Command(ctx, "-mode", "binary", "-show", "verbose", os.Args[0])
	err := cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}
	switch err := err.(type) {
	case nil:
	case interface{ ExitCode() int }:
		os.Exit(err.ExitCode())
	default:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
