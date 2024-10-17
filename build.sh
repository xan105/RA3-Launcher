#!/bin/sh

cd "$(dirname "$0")/src"
export GOOS=windows
export GOARCH=386
${GOPATH:-$HOME/go}/bin/go-winres make --in "../winres/winres.json"
echo "Compiling x86 (DEBUG)..."
go build -o "../build/Debug/RA3.exe" launcher
echo "Compiling x86 (RELEASE)..."
go build -ldflags "-w -s -H windowsgui" -o "../build/Release/RA3.exe" launcher

