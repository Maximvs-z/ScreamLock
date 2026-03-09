//go:build windows
// +build windows

// Package warning provides audio feedback before lock: tone and optional TTS.
package warning

import (
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var (
	kernel32     = syscall.NewLazyDLL("kernel32.dll")
	procBeep     = kernel32.NewProc("Beep")
	voiceMu      sync.Mutex
	voiceInUse   bool
	voiceRelease time.Time
)

const (
	// Tone: 800 Hz, 300 ms
	beepFreqHz   = 800
	beepDuration = 300
)

// PlayTone plays a short warning beep (Windows Beep API).
func PlayTone() {
	procBeep.Call(uintptr(beepFreqHz), uintptr(beepDuration))
}

// SpeakAsync says the text using Windows TTS (PowerShell System.Speech.Synthesis).
// It starts the speech and returns immediately. Only one instance runs at a time (guarded by mutex).
func SpeakAsync(text string) {
	voiceMu.Lock()
	defer voiceMu.Unlock()
	// Avoid starting a new speech if one was started very recently
	if voiceInUse && time.Since(voiceRelease) < 2*time.Second {
		return
	}
	voiceInUse = true
	voiceRelease = time.Now()
	go func() {
		defer func() {
			voiceMu.Lock()
			voiceInUse = false
			voiceRelease = time.Now()
			voiceMu.Unlock()
		}()
		// PowerShell with System.Speech.Synthesis (works offline)
		script := `Add-Type -AssemblyName System.Speech; $s = New-Object System.Speech.Synthesis.SpeechSynthesizer; $s.Speak('` + escapePSString(text) + `')`
		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
		_ = cmd.Run()
	}()
}

func escapePSString(s string) string {
	// Escape single quotes for PowerShell: ' -> ''
	var out []rune
	for _, r := range s {
		if r == '\'' {
			out = append(out, '\'', '\'')
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}

// RunSequence runs the three-stage sequence: tone, optional voice, wait, then calls lockFn.
// If enableVoice is false, only plays tone then waits ~0.3s then calls lockFn.
// If enableVoice is true: tone, start voice, wait 1s from voice start, call lockFn.
func RunSequence(enableVoice bool, lockFn func()) {
	PlayTone()
	if enableVoice {
		SpeakAsync("Please lower your voice.")
		time.Sleep(1 * time.Second)
	} else {
		time.Sleep(300 * time.Millisecond)
	}
	if lockFn != nil {
		lockFn()
	}
}

