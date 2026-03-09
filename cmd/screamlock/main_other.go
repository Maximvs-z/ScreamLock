//go:build !windows
// +build !windows

package main

func oleCoInitialize() error {
	return nil
}

func oleCoUninitialize() {}
