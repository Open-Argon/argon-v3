@echo off
set GOOS=js
set GOARCH=wasm
go build -o wasm/bin/argon.wasm ./src