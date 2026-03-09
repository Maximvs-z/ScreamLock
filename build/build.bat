@echo off
REM Build ScreamLock for Windows (no console window).
REM Run from repository root: build\build.bat
REM Requires Go installed: https://go.dev/dl/

cd /d "%~dp0.."
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-H windowsgui" -o build\screamlock.exe .\cmd\screamlock
if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
echo Built: %CD%\build\screamlock.exe
