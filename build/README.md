# Building ScreamLock

## Prerequisites

- [Go 1.21 or later](https://go.dev/dl/)

## Build from repository root

### Windows (Command Prompt or PowerShell)

**Main app (background monitor):**
```bat
cd path\to\ScreamLock
go build -ldflags "-H windowsgui" -o build\screamlock.exe .\cmd\screamlock
```

**Config app (GUI to choose microphone):**
```bat
go build -o build\screamlock-config.exe .\cmd\screamlock-config
```

**Setup installer (wizard + microphone dialog with level meter):**
```bat
go build -o build\screamlock-setup.exe .\cmd\screamlock-setup
```

Or run the script to build all three:

```bat
build\build.bat
```

**Single-file installer (for Program Files):**  
Builds an exe that installs to `C:\Program Files\ScreamLock`. Run **as Administrator** when distributing.

```bat
build\build-installer.bat
```

This produces **build\ScreamLock-Setup.exe**. It embeds `screamlock.exe` and `screamlock-config.exe`; when run, it extracts them to Program Files and adds a logon task.

### macOS / Linux (cross-compile for Windows)

```bash
cd /path/to/ScreamLock
GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui" -o build/screamlock.exe ./cmd/screamlock
GOOS=windows GOARCH=amd64 go build -o build/screamlock-config.exe ./cmd/screamlock-config
GOOS=windows GOARCH=amd64 go build -o build/screamlock-setup.exe ./cmd/screamlock-setup
```

To build the single-file installer (ScreamLock-Setup.exe), copy the two exes into `cmd/installer/files/` then build the installer:

```bash
cp build/screamlock.exe build/screamlock-config.exe cmd/installer/files/
GOOS=windows GOARCH=amd64 go build -o build/ScreamLock-Setup.exe ./cmd/installer
```

The `-H windowsgui` flag makes the executable a Windows GUI application so that **no console window** appears when it runs. This is required for the background “stealth” behavior.

## Output

- **screamlock.exe** — main monitor (no window; use Task Scheduler to run at logon).
- **screamlock-config.exe** — small GUI to choose microphone and sensitivity.
- **screamlock-setup.exe** — “next, next, next” wizard: asks about autostart, then opens the microphone dialog with a **live level meter**.
- **ScreamLock-Setup.exe** — Single-file installer (run `build-installer.bat` or the copy+build steps above). Installs to Program Files; run as Administrator.
- Single executables: no separate DLLs or runtime; safe to copy to another Windows PC.

## Running locally (development)

To test with a console (so you can see log output in the terminal):

```bat
go run .\cmd\screamlock
```

To test list-devices:

```bat
go run .\cmd\screamlock -list-devices
```

Note: `go run` will show a console. The built `screamlock.exe` with `-H windowsgui` will not.
