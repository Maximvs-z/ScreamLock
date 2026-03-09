//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "ScreamLock Config is only available on Windows.")
	os.Exit(1)
}
