@echo off
go build -trimpath -ldflags="-s -w" -o bin/argon.exe ./src