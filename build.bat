@echo off
set GOARCH=amd64
set GOOS=windows
set dateStr=%date:~5,2%%date:~8,2%%date:~11,2%
@REM garble -literals -seed=random build -ldflags="-s -w" -o autoBD2-%dateStr%.exe
go build -o autoBD2.exe
@REM go build -o autoBD2-%dateStr%.exe