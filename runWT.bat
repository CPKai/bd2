@echo off

:: 設置字符編碼為 UTF-8
chcp 65001 >nul

:: cd至bat檔存放的資料夾
cd /d "%~dp0"

:: 檢查是否已具有管理員權限
net session >nul 2>&1
if not %errorlevel%==0 (
    echo 以管理員模式重新啟動...
    powershell -Command "Start-Process -Verb RunAs -FilePath '%~f0'"
    exit /b
)

:: 將 autoBD2 還原回 exe 副檔名
set target_exe_name="autoBD2.ex"
set new_exe_name="autoBD2.exe"

:: 檢查目標檔案是否存在
if exist %new_exe_name% (
    echo %new_exe_name% 已存在，跳過更名動作。
) else if exist %target_exe_name% (
    ren %target_exe_name% %new_exe_name%
    echo 檔案已成功更名為 %new_exe_name%
) else (
    echo 檔案 %target_exe_name% 不存在，無法更名。
    exit
)

set echoStr1="請確認上方檔案清單是否有autoBD2執行檔，若有代表所在目錄正確，若沒有，請將此bat檔移至autoBD2的資料夾下再次執行。"
set echoStr2="目錄正確後，執行命令「.\autoBD2.exe」即可運行程式"

:: 啟動 Windows Terminal
start wt -p "Windows PowerShell" powershell.exe -NoExit -Command "& {cd '%~dp0'\; ls\; echo '%echoStr1%'\; echo '%echoStr2%'}"
exit