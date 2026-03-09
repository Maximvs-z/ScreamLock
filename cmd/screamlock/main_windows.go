//go:build windows
// +build windows

package main

import (
	"github.com/go-ole/go-ole"
)

func oleCoInitialize() error {
	return ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
}

func oleCoUninitialize() {
	ole.CoUninitialize()
}
