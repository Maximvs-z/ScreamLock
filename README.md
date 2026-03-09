# ScreamLock

ScreamLock monitors the microphone input level and **locks the Windows session** when the level exceeds a configurable threshold. It is intended to run continuously in the background—for example, while a child is using the computer. When the microphone picks up loud noise (e.g. screaming during a game), the screen locks and the user must log back in, providing a consistent consequence without parental intervention.

*While it may be used to address certain behavioral issues related to loud noises during screen time, the author is not a psychologist and cannot guarantee the effectiveness or lack of harm of this tool. Use at your own risk. You are welcome to modify and adapt it to your needs.*

---

## What It Does

- **Runs in the background** with no visible window, so it is not easily discovered or closed by the child.
- **Samples the microphone** at a set interval (default once per second).
- When the level goes above your threshold, it runs a **three-stage warning** (see below), then **locks the PC** (same as Win+L).
- **Survives reboots** when you set it to start at login (e.g. via Task Scheduler).
- **Lets you choose which microphone** to use via **ScreamLock Config** (small GUI) or the config file.
- **Logs to a file** instead of showing errors on screen, so you can troubleshoot without exposing the program.

The program is a **single Windows executable** (`.exe`). No installer or .NET runtime is required. You place the exe in a folder, configure it once, and optionally register it to run at startup.

### When the threshold is exceeded

When the configured sound threshold is exceeded, the application starts a warning sequence:

1. **Warning tone** — A short system beep plays immediately.
2. **Grace period** — For about 1.5 seconds after the tone, the microphone level is checked repeatedly. **If the level drops below the threshold during this time, the sequence is cancelled** and the app returns to normal monitoring (no voice, no lock). If the level stays at or above the threshold for the full 1.5 seconds, the sequence continues.
3. **Voice** — The app says “Please lower your voice.” (Windows text-to-speech; can be turned off in config).
4. **Lock** — About 1 second after the voice starts, the Windows session is locked.

So if the user lowers their voice right after the beep, the lock is avoided. After a lock, monitoring pauses for a **cooldown** (default 15 seconds, configurable) to avoid repeated triggers.

---

## Why It Exists

The original idea was a small script that would “inconvenience” the user (e.g. a child) by locking the session when the microphone level got too high—for example, when screaming during a game. This repository evolved from the **PowerShell script** (`ScreamLock.ps1`) into a **redesigned, maintainable application**: same goal, but implemented as a single Windows executable that is easier to install, configure, and run reliably.

---

## Quick Start (For Parents)

**Single-file installer (recommended):**  
Download **ScreamLock-Setup.exe** from the [Releases](https://github.com/Maximvs-z/ScreamLock/releases) page and run it. Windows will show an administrator prompt (UAC); choose **Yes** (or enter your password). The installer copies the apps to `C:\Program Files\ScreamLock`, adds a task to run ScreamLock at logon, and then **opens the config app** so you can choose your microphone and settings.

**Alternative — wizard only:** Run **ScreamLock Setup** (`screamlock-setup.exe`) from a folder where you’ve placed the exes. It’s a short wizard that asks about autostart and opens the microphone dialog with a level meter.

**Manual setup:**

1. **Get the programs**  
   From [Releases](https://github.com/Maximvs-z/ScreamLock/releases): use **ScreamLock-Setup.exe** (installer) or download `screamlock.exe`, `screamlock-config.exe`, and optionally `screamlock-setup.exe`.
2. **Put them in a folder**  
   e.g. `C:\Programs\ScreamLock`. (Or use the installer to install to Program Files.)
3. **Choose the microphone**  
   Run **screamlock-config.exe** (or use the installer). Pick the microphone, set sensitivity (dB), and in the installer you get a **live level meter** to test without locking.
4. **Run at startup**  
   In the installer, tick “Yes, run ScreamLock when I log on”, or in **screamlock-config.exe** click **“Run at Windows startup”**.

Full installation and configuration details: **[docs/INSTALL.md](docs/INSTALL.md)**.

---

## How It Runs

- **Normal run:** Double-clicking or starting `screamlock.exe` (or having Task Scheduler start it) runs the monitor. **No window appears**; it runs in the background.
- **Config folder:** Config and log files are stored under your user profile, typically:  
  `%APPDATA%\ScreamLock`  
  (e.g. `C:\Users\YourName\AppData\Roaming\ScreamLock`).
- **Config file:** `config.json` — device ID, threshold (dB), check interval, **cooldown_seconds** (pause after a lock), and **enable_voice_warning** (if `false`, only the tone plays before lock; if `true`, tone + spoken message + lock).
- **Log file:** `screamlock.log` — startup messages and errors. Use this to confirm it’s running or to troubleshoot.
- **ScreamLock-Setup.exe** — Single-file installer (from Releases). Prompts for Administrator (UAC), installs to `C:\Program Files\ScreamLock`, sets up logon task, then opens the config app.  
- **screamlock-setup.exe** — Wizard (autostart question → Finish) then microphone dialog with **live level meter** (use when exes are already in a folder).  
- **screamlock-config.exe** — Pick microphone, sensitivity, and **Run at Windows startup**. Saves to the same config.

---

## Troubleshooting

| Issue | What to do |
|--------|------------|
| **“Open microphone” error in log** | Run `screamlock.exe -list-devices` and check that the `device_id` in `config.json` matches an ID in `devices.txt`, or set `device_id` to `""` to use the default microphone. |
| **No log file / not sure if it’s running** | Check `%APPDATA%\ScreamLock\screamlock.log`. If the program started, you should see a line like “ScreamLock monitoring (threshold -50.00 dB …)”. |
| **Microphone was unplugged or changed** | ScreamLock will try to fall back to the default device. If problems persist, run `-list-devices` again and update `config.json` with the correct `device_id`. |
| **Locks too often or not at all** | Adjust `threshold_db` in `config.json`: more negative (e.g. `-60`) = less sensitive; less negative (e.g. `-40`) = more sensitive. Then restart the program (or log off and on if using Task Scheduler). |
| **Stopping the program** | Disable or delete the Task Scheduler task that runs ScreamLock, or use Task Manager to end the `screamlock.exe` process. |

---

## Build Instructions (For Developers)

- **Requirement:** Go 1.21 or later.
- **Produce a single Windows executable (no console):**
  ```bash
  GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui" -o build/screamlock.exe ./cmd/screamlock
  ```
  The `-H windowsgui` linker flag is what hides the console window.
- **From Windows:**  
  See [build/README.md](build/README.md) for `build.bat` and local testing.

The project layout:

- `cmd/screamlock/` — main monitor entrypoint
- `cmd/screamlock-config/` — small GUI to choose microphone and sensitivity (Windows only)
- `cmd/screamlock-setup/` — “next, next, next” wizard and microphone dialog with live level meter (Windows only)
- `cmd/installer/` — single-file installer that installs to Program Files (Windows only)
- `config/` — config load/save (JSON)
- `internal/audio/` — Windows capture device enumeration and peak level (WASAPI via go-wca)
- `internal/lock/` — Windows LockWorkStation
- `internal/logger/` — file logging
- `docs/` — installation and build docs
- `config.example.json` — example config

---

## Creating a release (for maintainers)

To publish a release so others can download **ScreamLock-Setup.exe**:

1. Create and push a version tag (e.g. `v1.0.0`):
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
2. The [release workflow](.github/workflows/release.yml) runs on GitHub Actions: it builds the installer and creates a **Release** with **ScreamLock-Setup.exe** attached.
3. Users download the installer from the [Releases](https://github.com/Maximvs-z/ScreamLock/releases) page; when they run it, Windows prompts for Administrator (UAC) and the installer then opens the config app.

---

## Pushing this repo to GitHub

1. If you fork the repo, update the Releases links in this README and in `docs/INSTALL.md` to your fork.
2. **Optional:** In `go.mod`, set the module path to match your repo so `go get` works from the new URL.
3. From the project root (the folder containing `go.mod`):
   ```bash
   git add .
   git commit -m "Redesign ScreamLock as Go app: single exe, WASAPI mic monitoring, config, Task Scheduler docs"
   git push origin main
   ```

---

## License and Disclaimer

You may use, modify, and distribute this software. The author is not responsible for any consequences of its use. This is not professional psychological or parenting advice; use at your own risk.
