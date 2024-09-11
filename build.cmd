@echo off
cd %~dp0src
set GOOS=windows
set GOARCH=386
go-winres make --in "..\winres\winres.json"
echo Compiling x86 (DEBUG)...
go build -o "..\build\Debug\RA3.exe" launcher
echo Compiling x86 (RELEASE)...
go build -ldflags "-w -s -H windowsgui" -o "..\build\Release\RA3.exe" launcher