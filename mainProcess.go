package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	robotgo "github.com/go-vgo/robotgo"
	gohook "github.com/robotn/gohook"
	bitmap "github.com/vcaesar/bitmap"
	imgo "github.com/vcaesar/imgo"
)

var (
	version        = "20250109006"
	debugSwitch    = false
	stopKey        = "f1"
	imgFormat      = ".png"
	prRatio_width  float64
	prRatio_height float64
	modeList       = map[string]string{
		"iGacha":    "無限抽抽樂-自動抽取",
		"info":      "查看「無限抽抽樂-自動抽取」的評分設定",
		"test":      "測試「autoBD2螢幕截圖」與「無樂抽抽樂-自動抽取」的評分邏輯",
		"lastnight": "未日之書-無限重複挑戰",
	}
)

func main() {

	// 啟動監聽 stopKey 的執行緒
	go exitEvent()
	fmt.Println("*******************************************************")
	fmt.Println("Version:" + version)
	fmt.Println("終止程序請按按鍵 " + stopKey)
	fmt.Println("*******************************************************")

	if len(os.Args) > 2 {
		if os.Args[2] == "debug" {
			debugSwitch = true
		}
	}

	// 加載設定檔
	configMap := load_map_str_float64("config.txt")
	if debugSwitch {
		fmt.Printf("args[]:%v\n", os.Args)
		fmt.Printf("configMap:%v\n", configMap)
	}
	setScreenPhysicResolutionRatio(configMap)

	fmt.Printf("\n\n")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			fmt.Printf("運行模式：%s\n", os.Args[1])
			activeTargetProcess("BrownDust")
			test(configMap)
		case "info":
			fmt.Printf("運行模式：%s\n", os.Args[1])
			settingInfo(configMap)
		case "lastnight":
			fmt.Printf("運行模式：%s\n", os.Args[1])
			activeTargetProcess("BrownDust")
			startLastNight(configMap)
		case "iGacha":
			fmt.Printf("運行模式：%s\n", os.Args[1])
			for i := 3; i > 0; i-- {
				fmt.Printf("腳本將在%d秒後開始運行...\n", i)
				time.Sleep(time.Duration(1000) * time.Millisecond)
			}
			activeTargetProcess("BrownDust")
			// 主邏輯運行
			start_infinite_gacha(configMap)
		default:
			fmt.Printf("查無你輸入的模式：%s，請確認是否輸入正確，模式說明如下。\n", os.Args[1])
			list_mode()
		}
	} else {
		fmt.Printf("未檢測到欲運行的模式，請確認欲運行模式為何。模式說明如下：\n\n")
		list_mode()
	}
}

func exitEvent() {

	// 設定按鍵事件
	exitEvent := gohook.AddEvents(stopKey)

	if exitEvent {
		// 確認退出
		fmt.Println("\n檢測到退出指令，程序即將結束...")
		os.Exit(0) // 結束程式
	}
}

func start_infinite_gacha(configMap map[string]float64) {

	// 加載設定檔
	imgMap := load_map_str_str("imgPathMap.txt")

	var tempPos []int
	round := 1

	for {

		fmt.Printf("\nCurrent Round: %d\n", round)
		time.Sleep(200 * time.Millisecond)

		findImage(imgMap["再抽一次"]+imgFormat, -1, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, int(configMap["睡眠參數(毫秒)"]), prRatio_width, prRatio_height)
		findImage(imgMap["確認"]+imgFormat, -1, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, int(configMap["睡眠參數(毫秒)"]), prRatio_width, prRatio_height)

		for {
			findImage(imgMap["skip1"]+imgFormat, 2, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, int(configMap["睡眠參數(毫秒)"]), prRatio_width, prRatio_height)
			findImage(imgMap["skip2"]+imgFormat, 2, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, int(configMap["睡眠參數(毫秒)"]), prRatio_width, prRatio_height)
			findImage(imgMap["skip3"]+imgFormat, 2, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, int(configMap["睡眠參數(毫秒)"]), prRatio_width, prRatio_height)

			tempPos = findImage(imgMap["再抽一次"]+imgFormat, 2, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], false, int(configMap["睡眠參數(毫秒)"]), prRatio_width, prRatio_height)
			if tempPos[0] > 0 {
				break
			}
		}

		fmt.Println("計分階段")
		calculateScore(configMap, imgMap)

		round++
	}
}

// 找到指定圖片(imgPath)在螢幕擷圖中的位置，並決定是否點擊(clickImg)
func findImage(imgPath string, scanMaxTime int, resolutionWidth int, resolutionHeight int, tolerance float64, clickImg bool, intervalTime int, prRatio_width float64, prRatio_height float64) []int {

	resultPos := []int{-1, -1}
	scanTime := 0

	for {
		scanTime++

		// 印出當前訊息
		printStr := "步驟: 比對「" + imgPath + "」圖片... || scanTime: " + strconv.Itoa(scanTime) + " / " + strconv.Itoa(scanMaxTime)
		fmt.Print(printStr)

		// 取得螢幕擷圖
		bitmapScreen := robotgo.CaptureScreen(0, 0, resolutionWidth, resolutionHeight)

		// 從螢幕擷圖中尋找目標圖片
		fx, fy := bitmap.FindPic(imgPath, bitmapScreen, tolerance)

		robotgo.FreeBitmap(bitmapScreen)

		// 找到對應的圖片在螢幕擷圖中的x,y
		if (fx != -1) && (fy != -1) {
			resultPos = []int{fx, fy}
			cx := int(math.Ceil(float64(fx) * prRatio_width))
			cy := int(math.Ceil(float64(fy) * prRatio_height))
			if debugSwitch {
				fmt.Printf("click_x:%d | click_y: %d\n", cx, cy)
			}
			if clickImg {
				time.Sleep(time.Duration(intervalTime) * time.Millisecond)
				robotgo.MoveClick(cx, cy, "left", false)
			}
			fmt.Printf("\n")
			return resultPos
		}

		// scanMaxTime >= 0 則啟用scan上限機制，達到上限後自動break
		if (scanTime >= scanMaxTime) && (scanMaxTime >= 0) {
			// 使用 \r 使光標回到行首，\033[0K 清除光標後內容
			fmt.Printf("\r\033[0K")
			break
		} else {
			time.Sleep(time.Duration(intervalTime) * time.Millisecond)
			// 使用 \r 使光標回到行首，\033[0K 清除光標後內容
			fmt.Printf("\r\033[0K")
		}
	}

	return resultPos
}

func calculateScore(configMap map[string]float64, imgMap map[string]string) {
	scoreMap := load_map_str_float64("scoreMap.txt")
	target_score := int(scoreMap["目標分數"])
	current_core := 0
	bitmap_screen := robotgo.CaptureScreen(0, 0, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]))

	for character, score := range scoreMap {
		if character != "目標分數" && character != "截圖分數" {
			var tempPosArr []robotgo.Point

			// 刷1次求快速
			if character == "5星角色" {
				tempPosArr = bitmap.FindAll(bitmap.Open(imgMap[character]+imgFormat), bitmap_screen, configMap["容忍值_小"])
			} else {
				tempPosArr = bitmap.FindAll(bitmap.Open(imgMap[character]+imgFormat), bitmap_screen, configMap["容忍值_角色"])
			}

			current_core += len(tempPosArr) * int(score)
			fmt.Printf("發現[%d]個[%s]，獲得分數[%d]，累計總分[%d]\n", len(tempPosArr), character, len(tempPosArr)*int(score), current_core)
			if debugSwitch {
				fmt.Printf("tempPosArr:「%v」\n", tempPosArr)
			}
		}
	}

	if current_core >= int(target_score) {
		fmt.Printf("當前總分 %d 大於等於目標分數 %d ，結束抽取。\n", current_core, int(target_score))
		fmt.Println()
		fmt.Println("終止程式")
		os.Exit(0)

	} else if current_core >= int(scoreMap["截圖分數"]) {
		if debugSwitch {
			img := robotgo.ToImage(bitmap_screen)
			imgName := getNextImageFileName("save")
			imgo.Save("save/"+imgName+"-"+strconv.Itoa(current_core)+".png", img)
		}
		fmt.Printf("當前總分 %d 小於目標分數 %d ，繼續下一輪。\n", current_core, int(target_score))
	} else {
		fmt.Printf("當前總分 %d 小於目標分數 %d ，繼續下一輪。\n", current_core, int(target_score))
	}

	robotgo.FreeBitmap(bitmap_screen)
}

func load_map_str_float64(csvPath string) map[string]float64 {

	// open csv
	csvFile, err := os.Open(csvPath)
	errHandler(err, "Open CSV["+csvPath+"] failed.")
	defer csvFile.Close()

	// load csv content
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	errHandler(err, "ErrCode 003001. Read CSV file failed.")

	// create map
	rowCount := len(csvLines)
	dataMap := make(map[string]float64, rowCount-1)

	for _, line := range csvLines {
		tempFloat, err := strconv.ParseFloat(line[1], 64)
		errHandler(err, "ErrCode 003003. Convert string to float error.")
		dataMap[line[0]] = tempFloat
	}

	return dataMap
}

func load_map_str_str(csvPath string) map[string]string {

	// open csv
	csvFile, err := os.Open(csvPath)
	errHandler(err, "Open CSV["+csvPath+"] failed.")
	defer csvFile.Close()

	// load csv content
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	errHandler(err, "ErrCode 003002. Read CSV file failed.")

	// create map
	rowCount := len(csvLines)
	dataMap := make(map[string]string, rowCount-1)

	for _, line := range csvLines {
		dataMap[line[0]] = line[1]
	}

	return dataMap
}

func errHandler(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

func activeTargetProcess(processName string) {
	pids, _ := robotgo.FindIds(processName) // 查找含有processName的pids

	// 把所有含有關鍵字的process都active
	for _, pid := range pids {
		robotgo.ActivePid(pid)
		fmt.Println("喚起BrownDust視窗，PID:", pid)
	}
}

func test(configMap map[string]float64) {

	fmt.Println("呼叫計分func測試")
	imgMap := load_map_str_str("imgPathMap.txt")
	calculateScore(configMap, imgMap)

	fmt.Println("進行螢幕截圖......")
	sbit := robotgo.CaptureScreen(0, 0, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]))
	fmt.Println("進行螢幕截圖的儲存(於執行檔目錄)，檔名：test.png")
	img := robotgo.ToImage(sbit)
	imgo.Save("test.png", img)

	robotgo.FreeBitmap(sbit)
}

func settingInfo(configMap map[string]float64) {

	// 加載設定檔
	imgMap := load_map_str_str("imgPathMap.txt")
	scoreMap := load_map_str_float64("scoreMap.txt")

	fmt.Println("config設定如下(小數點後多餘的0可忽略)")
	for k, v := range configMap {
		fmt.Printf("%s:%f\n", k, v)
	}
	fmt.Printf("螢幕解析度設定錯誤的話，請至「config.txt」中修正「螢幕解析度-寬」與「螢幕解析度-高」的值\n")
	fmt.Printf("img資料夾中預設圖片取自環境1920x1080解析度、FHD圖像，若本機環境不符記得更換圖片\n\n")

	fmt.Printf("目標分數:%d\n", int(scoreMap["目標分數"]))
	fmt.Printf("以下是各角色設定的分數與判斷圖片路徑\n")
	for character, score := range scoreMap {
		if character != "目標分數" {
			var imgPath string
			if imgMap[character] != "" {
				imgPath = imgMap[character] + imgFormat
			} else {
				imgPath = "未找到" + character + "的圖片路徑"
			}
			fmt.Printf("角色：%s | 分數:%d | 判定圖片路徑:%s\n", character, int(score), imgPath)
		}
	}
	fmt.Printf("目標分數或角色分數設定錯誤的話，請至「scoreMap.txt」中修正")
}

// 取得 save 資料夾中下一個圖片檔名
func getNextImageFileName(folderPath string) string {

	// 檢查資料夾是否存在
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 如果資料夾不存在，則建立資料夾
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			log.Fatalf("無法建立資料夾 %s: %v", folderPath, err)
		}
	}

	// 取得資料夾中的檔案
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("無法讀取資料夾 %s: %v", folderPath, err)
	}

	// 計算檔案數量並生成下一個圖片檔名
	fileNumb := len(files)
	return strconv.Itoa(fileNumb + 1)
}

// 設定螢幕物理寬高與解析度寬高的比率
func setScreenPhysicResolutionRatio(configMap map[string]float64) {

	pW, pH := robotgo.GetScreenSize()
	rW, rH := int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"])

	prRatio_width = float64(pW) / float64(rW)
	prRatio_height = float64(pH) / float64(rH)
	if debugSwitch {
		fmt.Printf("physical width:%d | resolution width:%d | widthRatio:%f\n", pW, rW, prRatio_width)
		fmt.Printf("physical height:%d | resolution height:%d | heightRatio:%f\n", pH, rH, prRatio_height)
	}
}

func startLastNight(configMap map[string]float64) {

	// 加載設定檔
	imgMap := load_map_str_str("imgPathMap.txt")
	round := 1

	for {
		fmt.Printf("\nCurrent Round: %d\n", round)
		time.Sleep(200 * time.Millisecond)
		findImage(imgMap["重新挑戰"]+imgFormat, -1, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, 100, prRatio_width, prRatio_height)
		findImage(imgMap["確認"]+imgFormat, -1, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, 100, prRatio_width, prRatio_height)
		findImage(imgMap["skip4"]+imgFormat, -1, int(configMap["螢幕解析度-寬"]), int(configMap["螢幕解析度-高"]), configMap["容忍值_小"], true, 100, prRatio_width, prRatio_height)
		round++
	}
}

func list_mode() {
	for k, v := range modeList {
		fmt.Printf("模式關鍵字:%s  |  模式說明:%s\n", k, v)
	}

	fmt.Printf("\n運行指定模式方法：\n")
	fmt.Printf("1.以「系統管理員」權限啟動「終端機」\n")
	fmt.Printf("2.「cd」至「autoBD2」資料夾目錄下\n")
	fmt.Printf("3.運行命令「.\\autoBD2.exe 模式關鍵字」即可運行指定模式\n")
	fmt.Printf("3-1.範例: 「.\\autoBD2.exe info」\n\n")
}
