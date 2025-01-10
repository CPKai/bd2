@echo off
set GOARCH=amd64
set GOOS=windows
set dateStr=%date:~5,2%%date:~8,2%%date:~11,2%
@REM garble -literals -seed=random build -ldflags="-s -w" -o autoBD2-%dateStr%.exe
go build -o autoBD2.ex
@REM go build -o autoBD2-%dateStr%.exe

del autoBD2.exe
ren autoBD2.ex autoBD2.exe