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
	"strings"
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
	// If not running as admin, re-launch with UAC (Yes/password) and exit
	if !isAdmin() {
		runElevated()
		os.Exit(0)
	}

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

	configExe := filepath.Join(installDir, "screamlock-config.exe")
	_ = exec.Command(configExe).Start()

	showInfo("ScreamLock has been installed to:\n" + installDir + "\n\nscreamlock-config.exe has been opened so you can choose your microphone and settings.\n\nA task was added to run ScreamLock when you log on. Use the config app to change the microphone or disable autostart.")
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

func isAdmin() bool {
	f, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func runElevated() {
	exe, err := os.Executable()
	if err != nil {
		showError("Could not get executable path: " + err.Error())
		return
	}
	cwd, _ := os.Getwd()
	var args string
	if len(os.Args) > 1 {
		args = strings.Join(os.Args[1:], " ")
	}
	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecuteW := shell32.NewProc("ShellExecuteW")
	const SW_SHOWNORMAL = 1
	r1, _, _ := shellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(mustUTF16("runas"))),
		uintptr(unsafe.Pointer(mustUTF16(exe))),
		uintptr(unsafe.Pointer(mustUTF16(args))),
		uintptr(unsafe.Pointer(mustUTF16(cwd))),
		SW_SHOWNORMAL,
	)
	if r1 <= 32 {
		showError("This installer needs Administrator rights. Please allow when Windows asks, or right-click the installer and choose \"Run as administrator\".")
	}
}
