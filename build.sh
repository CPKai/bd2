#!/bin/bash

# 設置 GOARCH 和 GOOS 環境變數
GOARCH=amd64
GOOS=windows

# 取得當前日期 (格式為 YYMMDD)
dateStr=$(date +%y%m%d)

# 生成執行檔案 (autoBD2.ex)
# 如需加密使用 garble，請取消註解相應行
# garble -literals -seed=random build -ldflags="-s -w" -o "autoBD2-${dateStr}.exe"
go build -o "autoBD2.ex"

# 如果需要產生帶日期的檔案名稱版本，請取消註解以下行
# go build -o "autoBD2-${dateStr}.exe"