//go:build windows
// +build windows

// Package lock provides Windows session lock via LockWorkStation.
// We use the Windows API directly so no external process (rundll32) is needed.
package lock

import (
	"syscall"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procLockWorkStation = user32.NewProc("LockWorkStation")
)

// LockWorkStation locks the Windows session (same as Win+L or Ctrl+Alt+Del -> Lock).
// Must be called from a process on the interactive desktop.
func LockWorkStation() bool {
	ret, _, _ := procLockWorkStation.Call()
	return ret != 0
}
