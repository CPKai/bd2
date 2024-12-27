@echo off

:: 檢查是否已具有管理員權限
net session >nul 2>&1
if not %errorlevel%==0 (
    echo 以管理員模式重新啟動...
    powershell -Command "Start-Process -Verb RunAs -FilePath '%~f0'"
    exit /b
)

set echoStr1="請確認上方檔案清單是否有autoBD2執行檔，若有代表所在目錄正確，若沒有，請將「runWT.bat」移至autoBD2的資料夾下再次執行。"
set echoStr2="目錄正確後，執行命令「.\autoBD2.exe」即可運行程式"

:: 啟動 Windows Terminal
start wt -p "Windows PowerShell" powershell.exe -NoExit -Command "& {cd '%~dp0'\; ls\; echo '%echoStr1%'\; echo '%echoStr2%'}"
exit