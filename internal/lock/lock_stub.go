//go:build !windows
// +build !windows

package lock

// LockWorkStation is only implemented on Windows.
func LockWorkStation() bool {
	return false
}
