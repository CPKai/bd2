# 使用流程-241226版本
1. 確認config.txt中，「螢幕解析度-寬」和「螢幕解析度-高」與主螢幕解析度的寬長相同
2. 確認scoreMap.txt中，「目標分數」與「各角色分數」按照自己預期的設定
3. 確認imgPathMap中，各個角色的判定圖片路徑正確(預設img資料夾下圖片是1920x1080、遊戲設定FHD環境下截取)
4. 確認遊戲以「全螢幕」、「FHD圖像」並在「主螢幕」運行，並停在「特別無限抽抽樂」抽取結果的畫面
5. 以系統管理員權限開啟終端機，運行執行檔
6. 達到「目標分數」或按下「F1」後，程式會停止，可切換到終端機查看log
7. 範例影片: https://youtu.be/ebD1dYi-Jfw

# 使用流程-250101後版本
1. 下載「autoBD2-日期.rar」壓縮包並解壓縮，取得資料夾「autoBD2」
2. 於autoBD2資料夾中新增文字文件，將「[runWT.bat內容連結](https://github.com/CPKai/bd2/blob/main/runWT.bat)」中，文本內容複製，並貼進文件中，然後儲存檔案
3. 修改剛剛新增的文件名稱，取名「runWT.bat」，並確認檔案類型轉為「Windows批次檔案」，而非文字檔(.txt)  
不知如何修改副檔名，可以Google「Win11 變更副檔名」，網上有許多相關教學
4. 雙擊「runWT.bat」後，會以管理員權限開啟「終端機」的PowerShell
5. 在上一步開啟的PowerShell中執行「.\autoBD2.exe」命令，會列出可運行的模式及其代碼
6. 執行「.\autoBD2.exe 欲運行模式代碼」即可運行該模式

# 使用上注意項目
- 需要管理員權限啟動終端機，用終端機去運行執行檔，這樣可以在終端機上看log(直接以系統管理員開exe的話，程式中止會直接關閉含有log的小黑窗)
- 螢幕截圖判斷是截主螢幕，遊戲請在主螢幕用全螢幕模式運行
- 有複數螢幕的狀況，請將主螢幕置於最左，使0,0的位置對齊主螢幕左上
- 遊戲圖像設定全螢幕、FHD。FPS可30可60，60判定成功頻率會高些
- 程式停止按鍵是F1，操作影片最後就是按F1後才切回終端機查看
- img資料夾內有放使用到的判斷圖片，範例是使用1920x1080解析度的圖片，若是其他解析度需要自行截圖使用
- 螢幕解析度非1920x1080的人，可以在config.txt中調整數值
- scoreMap.txt中有各角色的分數設定，與及格分數
- 計分方式以預設的範例為例，及格分是500，有截圖的角色每隻100分，5星100分。當抽到1隻b海時，等同5星(+100)、b海(+100)=200分。分數可自行調整。
- 要確認是否成功判定，可以在刷下一輪時按下f1中止程式，跳出到終端機查看最後一輪的分數計算，檢查計分結果。多螢幕的話直接把終端機放在其他的螢幕即可即時確認計分
- 要新增角色的話，需要新增判定圖至img資料夾下(不是"img/1920x1080"或"img/2560x1440")，再在imgPathMap.txt與scoreMap.txt新增角色資訊即可(可以參考既存的資料)
- 請勿移動鼠標，避免程式操作鼠標點擊時未點擊在正確位置

# 常見問題
Q: 有沒有辦法保持「終端機」覆蓋在遊戲上，一邊刷抽抽樂一邊看記錄?  
A: 在「終端機」的設定中，「外觀」設定裡，將「最上層顯示」設定開啟即可，但要注意不要遮到程式判斷角色或按鍵的位置

Q: 將config.txt中的解析度調成2K後，無法正常判斷  
A: 確認img資料夾中的圖片是否有更新，因原本是1920x1080解析度下截取的圖片，無法適用2K解析度的比對

Q: 「reroll」正確判斷，但鼠標點擊位置在畫面左下，且卡在「confirm」不斷scan  
A: 通常是雙螢幕且主螢幕在右側的狀況，將左螢幕換為主螢幕運行遊戲，應可正常運行

Q: 主螢幕在左側，卻卡在「confirm」不斷scan  
A: 很可能是「reroll」抓到位置要進行點擊時，有操作到鼠標或切換畫面，導致點擊失敗，確認的對話框沒跳出而卡在「confirm」，手動點擊再抽一次後應會繼續運行

Q: 卡在「confirm」不斷scan，也沒發現「reroll」判定成功時鼠標有移動  
A: 很可能是程式沒有系統管理員權限，無法操作鼠標。請以系統管理員權限開啟「終端機」來運行autoBD2

# 📜 免責聲明 (Disclaimer)
1. 本專案僅供個人學習與娛樂用途  
本專案旨在提供技術範例與學習資源，作者不對任何因使用此專案而產生的後果負責，包括但不限於遊戲帳號風險、服務條款違反或其他法律問題。
2. 使用者需自行承擔風險  
使用本專案時可能會涉及遊戲開發商的使用條款或相關規範，請使用者自行判斷是否適當。因使用本專案導致的任何問題，作者概不負責。
3. 不提供維護與支援  
本專案為作者的非商業性自發作品，無法保證未來維護更新或任何技術支援。