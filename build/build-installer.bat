@echo off
REM Build the single-file installer (ScreamLock-Setup.exe) that installs to Program Files.
REM Run from repository root: build\build-installer.bat
REM Prerequisite: run build\build.bat first, or this script builds the exes then the installer.

cd /d "%~dp0.."
set GOOS=windows
set GOARCH=amd64

REM Build main exes if not present
if not exist build\screamlock.exe (
  go build -ldflags "-H windowsgui" -o build\screamlock.exe .\cmd\screamlock
  if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
)
if not exist build\screamlock-config.exe (
  go build -o build\screamlock-config.exe .\cmd\screamlock-config
  if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%
)

REM Copy exes into installer embed folder
copy /Y build\screamlock.exe cmd\installer\files\
copy /Y build\screamlock-config.exe cmd\installer\files\

REM Build installer (embeds the two exes)
go build -o build\ScreamLock-Setup.exe .\cmd\installer
if %ERRORLEVEL% neq 0 exit /b %ERRORLEVEL%

echo Built: build\ScreamLock-Setup.exe
echo Run as Administrator to install to Program Files\ScreamLock.
