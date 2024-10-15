@echo off
go build -trimpath -ldflags="-s -w" -tags WINDOWS -o bin/argon.exe ./src