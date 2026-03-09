# Building ScreamLock

## Prerequisites

- [Go 1.21 or later](https://go.dev/dl/)

## Build from repository root

### Windows (Command Prompt or PowerShell)

```bat
cd path\to\ScreamLock
go build -ldflags "-H windowsgui" -o build\screamlock.exe .\cmd\screamlock
```

Or run the script:

```bat
build\build.bat
```

### macOS / Linux (cross-compile for Windows)

```bash
cd /path/to/ScreamLock
GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui" -o build/screamlock.exe ./cmd/screamlock
```

The `-H windowsgui` flag makes the executable a Windows GUI application so that **no console window** appears when it runs. This is required for the background “stealth” behavior.

## Output

- **Path:** `build/screamlock.exe`
- **Single executable:** no separate DLLs or runtime; safe to copy to another Windows PC.

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
