@echo off
REM Build ScreamLock and ScreamLock Config for Windows.
REM Run from repository root: build\build.bat
REM Requires Go installed: https://go.dev/dl/

cd /d "%~dp0.."
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-H windowsgui" -o build\screamlock.exe .\cmd\screamlock
if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
go build -o build\screamlock-config.exe .\cmd\screamlock-config
if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
echo Built: build\screamlock.exe and build\screamlock-config.exe
