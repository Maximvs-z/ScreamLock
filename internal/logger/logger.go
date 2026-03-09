// Package logger provides file-based logging so errors are not shown on screen.
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var logFile *os.File

// Init initializes the logger to write to a file in the config directory.
// Safe to call multiple times; subsequent calls are no-ops if already initialized.
func Init(configDir string) error {
	if logFile != nil {
		return nil
	}
	path := filepath.Join(configDir, "screamlock.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logFile = f
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime)
	return nil
}

// Close closes the log file. Call before exit if desired.
func Close() {
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}
}

// Errorf logs an error message.
func Errorf(format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args...)
}

// Infof logs an info message.
func Infof(format string, args ...interface{}) {
	log.Printf("INFO: "+format, args...)
}

// Fatal logs and exits. Use only when the process cannot continue.
func Fatal(v ...interface{}) {
	log.Print(append([]interface{}{"FATAL: "}, v...)...)
	fmt.Fprint(logFile, "FATAL: ")
	fmt.Fprintln(logFile, v...)
	Close()
	os.Exit(1)
}
