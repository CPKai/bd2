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
	debugSwitch    = false
	stopKey        = "f1"
	imgFormat      = ".png"
	prRatio_width  float64
	prRatio_height float64
)

func main() {

	// 啟動監聽 stopKey 的執行緒
	go exitEvent()
	fmt.Println("*******************************************************")
	fmt.Println("終止程序請按按鍵 " + stopKey)
	fmt.Println("*******************************************************")
	fmt.Println("查看設定，請下autoBD2.exe info")
	fmt.Printf("測試程式截圖與計分邏輯，請下autoBD2.exe test\n\n")

	// 加載設定檔
	configMap := load_map_str_float64("config.txt")
	setScreenPhysicResolutionRatio(configMap)

	if debugSwitch {
		fmt.Printf("args[]:%v\n", os.Args)
		fmt.Printf("configMap:%v\n", configMap)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			activeTargetProcess("BrownDust")
			fmt.Printf("運行測試\n")
			test(configMap)
		case "info":
			fmt.Printf("列出設定內容\n")
			settingInfo(configMap)
		case "debug":
			debugSwitch = true
			start_infinite_gacha(configMap)
		default:
			for i := 3; i > 0; i-- {
				fmt.Printf("腳本將在%d秒後開始運行...\n", i)
				time.Sleep(time.Duration(1000) * time.Millisecond)
			}
			activeTargetProcess("BrownDust")
			// 主邏輯運行
			start_infinite_gacha(configMap)
		}
	} else {
		for i := 3; i > 0; i-- {
			fmt.Printf("腳本將在%d秒後開始運行...\n", i)
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}
		activeTargetProcess("BrownDust")
		// 主邏輯運行
		start_infinite_gacha(configMap)
	}
}

func exitEvent() {

	// 設定按鍵事件
	exitEvent := gohook.AddEvents(stopKey)

	if exitEvent {
		// 確認退出
		fmt.Println("\n檢測到退出指令，程序即將結束...")
		pid := robotgo.GetPid()
		robotgo.ActivePid(pid)
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

		findImage(imgMap["reroll"]+imgFormat, -1, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]), configMap["tolerance_s"], true, 100, prRatio_width, prRatio_height)
		findImage(imgMap["confirm"]+imgFormat, -1, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]), configMap["tolerance_s"], true, 100, prRatio_width, prRatio_height)

		for {
			findImage(imgMap["skip1"]+imgFormat, 2, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]), configMap["tolerance_s"], true, 100, prRatio_width, prRatio_height)
			findImage(imgMap["skip2"]+imgFormat, 2, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]), configMap["tolerance_s"], true, 100, prRatio_width, prRatio_height)
			findImage(imgMap["skip3"]+imgFormat, 2, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]), configMap["tolerance_s"], true, 100, prRatio_width, prRatio_height)

			tempPos = findImage(imgMap["reroll"]+imgFormat, 2, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]), configMap["tolerance_s"], false, 100, prRatio_width, prRatio_height)
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
		printStr := "step: " + imgPath + " || scanTime: " + strconv.Itoa(scanTime) + " / " + strconv.Itoa(scanMaxTime)
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
	bitmap_screen := robotgo.CaptureScreen(0, 0, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]))

	for character, score := range scoreMap {
		if (character != "目標分數") && (score > 0) {
			var tempPosArr []robotgo.Point

			// 刷1次求快速
			if character == "5星角色" {
				tempPosArr = bitmap.FindAll(bitmap.Open(imgMap[character]+imgFormat), bitmap_screen, configMap["tolerance_s"])
			} else {
				tempPosArr = bitmap.FindAll(bitmap.Open(imgMap[character]+imgFormat), bitmap_screen, configMap["tolerance"])
			}

			current_core += len(tempPosArr) * int(score)
			fmt.Printf("當前計分對象[%s]，獲得分數[%d]，當前分數[%d]\n", character, len(tempPosArr)*int(score), current_core)
			if debugSwitch {
				fmt.Printf("tempPosArr:[%v] | len(tempPosArr):%d\n", tempPosArr, len(tempPosArr))
			}
		}
	}
	robotgo.FreeBitmap(bitmap_screen)

	if current_core >= int(target_score) {
		fmt.Printf("當前總分 %d 大於等於目標分數 %d ，結束抽取。\n", current_core, int(target_score))
		fmt.Println()
		fmt.Println("終止程式")
		os.Exit(0)

	} else if current_core >= 300 {
		if debugSwitch {
			bitmap_screen := robotgo.CaptureScreen(0, 0, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]))
			img := robotgo.ToImage(bitmap_screen)
			imgName := getNextImageFileName("save")
			imgo.Save("save/"+imgName+"-"+strconv.Itoa(current_core)+".png", img)
			robotgo.FreeBitmap(bitmap_screen)
		}
		fmt.Printf("當前總分 %d 小於目標分數 %d ，繼續下一輪。\n", current_core, int(target_score))
	} else {
		fmt.Printf("當前總分 %d 小於目標分數 %d ，繼續下一輪。\n", current_core, int(target_score))
	}
}

func load_map_str_float64(csvPath string) map[string]float64 {

	// open csv
	csvFile, err := os.Open(csvPath)
	errHandler(err, "Open CSV["+csvPath+"] failed.")
	defer csvFile.Close()

	// load csv content
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	errHandler(err, "ErrCode 003002. Read CSV file failed.")

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
		fmt.Println("Activated window for PID:", pid)
	}
}

func test(configMap map[string]float64) {

	fmt.Println("呼叫計分func測試")
	imgMap := load_map_str_str("imgPathMap.txt")
	calculateScore(configMap, imgMap)

	fmt.Println("進行螢幕截圖......")
	sbit := robotgo.CaptureScreen(0, 0, int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]))
	fmt.Println("進行螢幕截圖的儲存(於執行檔目錄)，檔名：test.png")
	img := robotgo.ToImage(sbit)
	imgo.Save("test.png", img)

	robotgo.FreeBitmap(sbit)

	pid := robotgo.GetPid()
	robotgo.ActivePid(pid)
}

func settingInfo(configMap map[string]float64) {

	// 加載設定檔
	imgMap := load_map_str_str("imgPathMap.txt")
	scoreMap := load_map_str_float64("scoreMap.txt")

	fmt.Printf("螢幕解析度-寬:%d | 長:%d\n", int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"]))
	fmt.Println("螢幕解析度設定錯誤的話，請至「config.txt」中修正ScreenResolutionWidth(寬)與ScreenResolutionHeight(高)的值")

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
}

// 取得 save 資料夾中下一個圖片檔名
func getNextImageFileName(folderPath string) string {

	// 取得資料夾中的檔案
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("無法讀取資料夾 %s: %v", folderPath, err)
	}

	fileNumb := len(files)
	// 生成下一個圖片檔名
	return strconv.Itoa(fileNumb + 1)
}

// 設定螢幕物理寬高與解析度寬高的比率
func setScreenPhysicResolutionRatio(configMap map[string]float64) {

	pW, pH := robotgo.GetScreenSize()
	rW, rH := int(configMap["ScreenResolutionWidth"]), int(configMap["ScreenResolutionHeight"])

	prRatio_width = float64(pW) / float64(rW)
	prRatio_height = float64(pH) / float64(rH)
	if debugSwitch {
		fmt.Printf("configMap:%v\nconfigMap[\"ScreenResolutionWidth\"]:%f\n", configMap, configMap["ScreenResolutionWidth"])
		fmt.Printf("physical width:%d | resolution width:%d | widthRatio:%f\n", pW, rW, prRatio_width)
		fmt.Printf("physical height:%d | resolution height:%d | heightRatio:%f\n", pH, rH, prRatio_height)
	}
}
