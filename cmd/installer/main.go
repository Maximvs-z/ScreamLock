// ScreamLock Installer: copies screamlock.exe and screamlock-config.exe to
// Program Files\ScreamLock. Run as Administrator. Build with embedded exes:
//   copy build\screamlock.exe build\screamlock-config.exe cmd\installer\files\
//   go build -o build\ScreamLock-Setup.exe .\cmd\installer
//
//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	installDirName = "ScreamLock"
	taskName      = "ScreamLock"
)

//go:embed files/screamlock.exe files/screamlock-config.exe
var embedded embed.FS

func main() {
	programFiles := os.Getenv("ProgramFiles")
	if programFiles == "" {
		programFiles = `C:\Program Files`
	}
	installDir := filepath.Join(programFiles, installDirName)

	if err := os.MkdirAll(installDir, 0755); err != nil {
		showError(fmt.Sprintf("Cannot create folder:\n%s\n\nError: %v\n\nPlease run as Administrator (right-click → Run as administrator).", installDir, err))
		os.Exit(1)
	}

	// Extract embedded exes
	files := []string{"files/screamlock.exe", "files/screamlock-config.exe"}
	for _, name := range files {
		data, err := embedded.ReadFile(name)
		if err != nil {
			showError(fmt.Sprintf("Missing embedded file %s: %v", name, err))
			os.Exit(1)
		}
		base := filepath.Base(name)
		dest := filepath.Join(installDir, base)
		if err := os.WriteFile(dest, data, 0755); err != nil {
			showError(fmt.Sprintf("Cannot write %s:\n%v\n\nPlease run as Administrator.", dest, err))
			os.Exit(1)
		}
	}

	// Optional: create Task Scheduler entry for screamlock.exe at logon
	screamlockExe := filepath.Join(installDir, "screamlock.exe")
	_ = createLogonTask(screamlockExe)

	showInfo("ScreamLock has been installed to:\n" + installDir + "\n\nscreamlock.exe — runs in the background (no window).\nscreamlock-config.exe — choose microphone and settings.\n\nA task was added to run ScreamLock when you log on. Open screamlock-config.exe to change the microphone or disable autostart.")
}

var (
	user32    = syscall.NewLazyDLL("user32.dll")
	messageBox = user32.NewProc("MessageBoxW")
)

const (
	mbOK             = 0x00000000
	mbIconError      = 0x00000010
	mbIconInformation = 0x00000040
)

func showError(msg string) {
	messageBox.Call(0, uintptr(unsafe.Pointer(mustUTF16(msg))), uintptr(unsafe.Pointer(mustUTF16("ScreamLock Installer"))), mbOK|mbIconError)
}

func showInfo(msg string) {
	messageBox.Call(0, uintptr(unsafe.Pointer(mustUTF16(msg))), uintptr(unsafe.Pointer(mustUTF16("ScreamLock Installer"))), mbOK|mbIconInformation)
}

func mustUTF16(s string) *uint16 {
	u, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return u
}

func createLogonTask(exePath string) error {
	cmd := exec.Command("schtasks", "/Create", "/TN", taskName, "/TR", exePath, "/SC", "ONLOGON", "/F")
	return cmd.Run()
}
