@echo off

:: run the go run command passing the path to the main.go file, with the working directory set to the bin folder. pass in the arguments

set __ARGON_DEBUG__=true
go run ./src %*