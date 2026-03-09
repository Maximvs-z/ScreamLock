// ScreamLock monitors microphone input level and locks the Windows session when
// the level exceeds a configured threshold. Intended for parental use to discourage
// loud behavior (e.g. screaming) during screen time. Runs as a background process
// with no visible window; configure via config file and -list-devices.
//
// Build for Windows (no console): go build -ldflags "-H windowsgui" -o screamlock.exe .
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/screamlock/screamlock/config"
	"github.com/screamlock/screamlock/internal/audio"
	"github.com/screamlock/screamlock/internal/lock"
	"github.com/screamlock/screamlock/internal/logger"
)

func main() {
	listDevices := flag.Bool("list-devices", false, "List capture devices to a file and open the config folder")
	flag.Parse()

	cfg, configPath, err := config.Load()
	if err != nil {
		// No logger yet; try to init after we have config dir
		dir, _ := config.ConfigDir()
		_ = logger.Init(dir)
		logger.Errorf("Load config: %v", err)
		os.Exit(1)
	}

	configDir := filepath.Dir(configPath)
	if err := logger.Init(configDir); err != nil {
		os.Exit(1)
	}
	defer logger.Close()

	if *listDevices {
		runListDevices(configDir)
		return
	}

	runMonitor(cfg)
}

func runListDevices(configDir string) {
	// CoInitialize is required for COM (WASAPI) on the main thread
	if err := oleCoInitialize(); err != nil {
		logger.Errorf("COM init: %v", err)
		return
	}
	defer oleCoUninitialize()

	devices, err := audio.ListCaptureDevices()
	if err != nil {
		logger.Errorf("List devices: %v", err)
		return
	}

	path := filepath.Join(configDir, "devices.txt")
	f, err := os.Create(path)
	if err != nil {
		logger.Errorf("Create devices.txt: %v", err)
		return
	}
	fmt.Fprintln(f, "Available microphone (capture) devices. Copy the ID into config.json \"device_id\" if you want a specific device, or leave device_id empty for default.")
	fmt.Fprintln(f, "")
	for _, d := range devices {
		fmt.Fprintf(f, "ID:   %s\nName: %s\n\n", d.ID, d.Name)
	}
	if len(devices) == 0 {
		fmt.Fprintln(f, "(No capture devices found.)")
	}
	f.Close()

	// Open folder in Explorer so the parent can see devices.txt and config.json
	explorerPath := "explorer.exe"
	_ = exec.Command(explorerPath, "/select,"+path).Run()
	logger.Infof("Device list written to %s", path)
}

func runMonitor(cfg config.Config) {
	if err := oleCoInitialize(); err != nil {
		logger.Errorf("COM init: %v", err)
		return
	}
	defer oleCoUninitialize()

	reader, err := audio.OpenPeakReader(cfg.DeviceID)
	if err != nil || reader == nil {
		logger.Errorf("Open microphone: %v (check device_id in config or run with -list-devices)", err)
		return
	}
	defer reader.Close()

	thresholdLinear := audio.LinearFromDB(cfg.ThresholdDB)
	interval := time.Duration(cfg.CheckIntervalSeconds) * time.Second
	if interval < time.Second {
		interval = time.Second
	}

	logger.Infof("ScreamLock monitoring (threshold %.2f dB = %.4f linear); interval %v", cfg.ThresholdDB, thresholdLinear, interval)

	for {
		peak, err := reader.Peak()
		if err != nil {
			logger.Errorf("Peak read: %v", err)
			time.Sleep(interval)
			continue
		}
		if peak > thresholdLinear {
			logger.Infof("Peak %.4f above threshold %.4f — locking workstation", peak, thresholdLinear)
			if lock.LockWorkStation() {
				logger.Infof("Workstation locked")
			} else {
				logger.Errorf("LockWorkStation failed")
			}
			// Brief cooldown after lock to avoid immediate re-lock
			time.Sleep(5 * time.Second)
		}
		time.Sleep(interval)
	}
}
