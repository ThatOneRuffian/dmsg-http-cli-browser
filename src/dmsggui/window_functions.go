package dmsggui

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

//mainMenuPageCountMax holds the max page of the server list
var mainMenuPageCountMax int

//serverPageCountMax holds the max page given the current server's index and terminal height
var serverPageCountMax int

var downloadQueuePageMax int

//SavedServers stores server cache - initalized on loadCache
var SavedServers map[int][2]string

//currentServerIndexContents will store the parsed server index values
var currentServerIndexContents map[int][2]string

//current nav filter for search function
var currentDirFilter string

func ClearScreen() {
	cmd := exec.Command("")
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		cmd = exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd = exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// main server list

func renderServerBrowser() {
	bufferHeight := 7 //lines consumed by menu elements
	dirNumberOfItems := len(SavedServers)
	terminalHeightAvailable, terminalWidthAvailable := getTerminalDims(bufferHeight)
	mainMenuPageCountMax = calcNumberOfPages(dirNumberOfItems, terminalHeightAvailable)

	// Avoid 1/0 pages
	if mainMenuPageCountMax == 0 {
		mainMenuPageCountMax = 1
	}

	//Create header divider of appropriate length
	divider := ""
	for i := 0; i < terminalWidthAvailable; i++ {
		divider += "="
	}

	//Render variables
	titleBuffer := ""
	menuTitle := "SERVER DOWNLOAD INDEX"
	currentDir := getPresentWorkingDirectory()
	tmpTitle := fmt.Sprintf("%s%s", menuTitle, currentDir)
	titleBufferLength := terminalWidthAvailable - len(tmpTitle)
	for i := 0; i < titleBufferLength; i++ {
		titleBuffer = titleBuffer + " "
	}
	pageStatus := fmt.Sprintf("page (%d / %d)", mainMenuBrowserIndex+1, mainMenuPageCountMax)

	ClearScreen()
	fmt.Println(divider)
	fmt.Println("DMSG-HTTP SERVER LIST")
	fmt.Println(divider)

	renderHomeMenuServerList(terminalHeightAvailable, terminalWidthAvailable)

	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<<F <B | N> L>>")
}

func renderHomeMenuServerList(terminalHeightAvailable int, terminalWidthAvailable int) {
	verticalHeightBuffer := terminalHeightAvailable
	var sortedIndex []int
	for index := range SavedServers {
		sortedIndex = append(sortedIndex, index)
	}

	sort.Ints(sortedIndex)

	for index := range sortedIndex {
		indexOffset := index + (terminalHeightAvailable * mainMenuBrowserIndex)

		serverFriendlyName := SavedServers[indexOffset][0]

		tmpLineEntry := fmt.Sprintf("%d) %s ", indexOffset+1, serverFriendlyName)
		horizontalFill := ""
		for i := terminalWidthAvailable - len(tmpLineEntry); i > 0; i-- {
			horizontalFill += "-"
		}
		lineEntry := fmt.Sprintf("%d) %s %s", indexOffset+1, serverFriendlyName, horizontalFill)

		fmt.Println(lineEntry)

		verticalHeightBuffer--
		if verticalHeightBuffer == 0 {
			goto END
		}

	}

	for ; verticalHeightBuffer > 0; verticalHeightBuffer-- {
		fmt.Println("-")
	}
END:
}

func truncateStringTo(stringToTruncate string, rawMenuLength int, terminalWidthAvailable int) string {
	screenWidthLeft := terminalWidthAvailable - rawMenuLength
	truncateBuffer := screenWidthLeft / 2 //amount to keep on beginning and end of string
	entryNameLength := len(stringToTruncate)
	if entryNameLength > screenWidthLeft {
		//truncate string here

		if screenWidthLeft-truncateBuffer*2 >= 0 {
			stringBuffer := "~~~~"
			tmpEntryNameStart := stringToTruncate[0 : truncateBuffer-len(stringBuffer)/2]
			tmpEntryNameEnd := stringToTruncate[entryNameLength-truncateBuffer+len(stringBuffer)/2:]
			stringToTruncate = fmt.Sprintf("%s%s%s", tmpEntryNameStart, stringBuffer, tmpEntryNameEnd)
		} else {
			tmpStartString := stringToTruncate[:int(len(stringToTruncate)/2)-3]
			tmpEndString := stringToTruncate[int(len(stringToTruncate)/2):]
			stringToTruncate = tmpStartString + "~~~" + tmpEndString
		}
	}
	return stringToTruncate
}

// download queue

func renderDownloadQueuePage() {
	ClearScreen()
	bufferHeight := 7 //lines consumed by menu elements
	numberOfQueueItems := len(downloadQueue)
	terminalHeightAvailable, terminalWidthAvailable := getTerminalDims(bufferHeight)
	downloadQueuePageMax = calcNumberOfPages(numberOfQueueItems, terminalHeightAvailable)

	// Avoid 1/0 pages
	if downloadQueuePageMax == 0 {
		downloadQueuePageMax = 1
	}

	//Create header divider of appropriate length
	divider := ""
	for i := 0; i < terminalWidthAvailable; i++ {
		divider += "="
	}

	//Render variables
	titleBuffer := ""
	pageStatus := fmt.Sprintf("page (%d / %d)", downloadQueueIndex+1, downloadQueuePageMax)
	menuTitle := "DOWNLOAD QUEUE"
	tmpTitle := fmt.Sprintf("%s Download Progress", menuTitle)
	titleBufferLength := terminalWidthAvailable - len(tmpTitle)

	if titleBufferLength < 0 {

	} else {
		// fill empty space
		for i := 0; i < titleBufferLength; i++ {
			titleBuffer = titleBuffer + " "
		}
	}

	menuHeader := fmt.Sprintf("%s%s Download Progress", menuTitle, titleBuffer)

	//Render download queue
	fmt.Println(divider)
	fmt.Println(menuHeader)
	fmt.Println(divider)
	renderDownloadQueueMetaData(terminalHeightAvailable, terminalWidthAvailable)
	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<<F <B | N> L>>")
}

func renderDownloadQueueMetaData(terminalHeightAvailable int, terminalWidthAvailable int) {
	verticalHeightBuffer := terminalHeightAvailable
	horizontalFill := ""
	// sort index
	indexes := []int{}
	for index := range downloadQueue {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)

	for _, index := range indexes {
		indexOffset := index + downloadQueueIndex*terminalHeightAvailable
		if _, ok := downloadQueue[indexOffset]; ok {
			fileSize := downloadQueue[indexOffset].fileSize
			fileName := downloadQueue[indexOffset].fileName

			currentFileSize := getDownloadFileSize(fileName)
			downloadPercentage := (currentFileSize / downloadQueue[indexOffset].fileSize) * 100
			// check if download is actually 100% and non-zero
			if (downloadPercentage == 100) && (currentFileSize >= fileSize) && *downloadQueue[indexOffset].downloadStatus || downloadQueue[indexOffset].fileSize == 0 {
				markAsDone := false
				downloadPercentage = 100
				*downloadQueue[indexOffset].downloadStatus = markAsDone
			} else if downloadPercentage == 100 && *downloadQueue[indexOffset].downloadStatus {
				downloadPercentage = 99
			}
			tmpBasicLineEntry := fmt.Sprintf("%d)   (%.0f/%.0f)  %.0f%%", indexOffset+1, currentFileSize, fileSize, downloadPercentage)
			tmpLineEntry := tmpBasicLineEntry + fileName
			spaceToFill := terminalWidthAvailable - len(tmpLineEntry)
			if len(tmpLineEntry) >= terminalWidthAvailable {
				fileName = truncateStringTo(fileName, len(tmpBasicLineEntry), terminalWidthAvailable)
			}
			for i := 0; i < spaceToFill; i++ {
				horizontalFill += "-"
			}

			lineEntry := fmt.Sprintf("%d) %v %s (%.0f/%.0f)  %.0f%%", indexOffset+1, fileName, horizontalFill, currentFileSize, fileSize, downloadPercentage)
			fmt.Println(lineEntry)
			horizontalFill = ""
			verticalHeightBuffer--
		}

		if verticalHeightBuffer == 0 {
			break
		}
	}

	//vertical buffer
	for ; verticalHeightBuffer > 0; verticalHeightBuffer-- {
		fmt.Println("-")
	}
}

// server download index

func renderServerDownloadList() map[int]map[string]bool {
	bufferHeight := 7 //lines consumed by menu elements
	notificationBar := ""
	dirNumberOfItems := len(navPtr.subDirs) + len(navPtr.files)
	terminalHeightAvailable, terminalWidthAvailable := getTerminalDims(bufferHeight)
	serverPageCountMax = calcNumberOfPages(dirNumberOfItems, terminalHeightAvailable)

	// Avoid 1/0 pages
	if serverPageCountMax == 0 {
		serverPageCountMax = 1
	}

	//Create header divider of appropriate length
	divider := ""
	for i := 0; i < terminalWidthAvailable; i++ {
		divider += "="
		notificationBar += "="
	}

	//Render variables
	titleBuffer := ""
	pageStatus := ""
	currentFilterStringStatus := ""
	menuTitle := "SERVER DOWNLOAD INDEX"
	dirMetaData := getCurrentDirMetaData()
	currentDir := getPresentWorkingDirectory()
	tmpTitle := fmt.Sprintf("%s%s", menuTitle, currentDir)
	titleBufferLength := terminalWidthAvailable - len(tmpTitle)
	truncateIndex := 0
	truncateBuffer := ""
	presentWorkingDirTitle := ""

	if titleBufferLength*-1 >= len(currentDir) {
		//check case if space to be buffered is longer than len(currentDir) (overflow), do not render PWD
		truncateIndex = 0
		truncateBuffer = ""
		presentWorkingDirTitle = ""
	} else {

		if titleBufferLength < 0 {
			//check if title overflow
			truncateIndex = titleBufferLength * -1
			truncateBuffer = " ~~~"
		} else {
			// fill empty space
			for i := 0; i < titleBufferLength; i++ {
				titleBuffer = titleBuffer + " "
			}
		}

		if truncateIndex+len(truncateBuffer) < len(currentDir) {
			presentWorkingDirTitle = truncateBuffer + currentDir[truncateIndex+len(truncateBuffer):]

		} else {
			presentWorkingDirTitle = ""
		}
	}

	menuHeader := fmt.Sprintf("%s%s%s", menuTitle, titleBuffer, presentWorkingDirTitle)
	// format header to show download notifications
	if downloadNotification != "" {
		notificationBar = "= " + downloadNotification + " ="
		tmpTitleLength := len(notificationBar)
		numberOfIterations := (terminalWidthAvailable - tmpTitleLength) / 2

		// check if title overflow
		if terminalWidthAvailable-len(notificationBar) < 0 {
			notificationBar = truncateStringTo(notificationBar, 0, terminalWidthAvailable)
		} else {
			// add notification bar padding
			for i := 0; i < (terminalWidthAvailable-tmpTitleLength)/2; i++ {
				notificationBar = "=" + notificationBar + "="
			}
		}
		// if notification bar requires extra padding
		if terminalWidthAvailable-(numberOfIterations*2+tmpTitleLength) > 0 {
			notificationBar += "="
		}
	}

	// format menu to show current search string
	if len(currentDirFilter) == 0 {
		pageStatus = fmt.Sprintf("page (%d / %d)", downloadBrowserIndex+1, serverPageCountMax)
		currentFilterStringStatus = divider
	} else {
		serverPageCountMax = calcNumberOfPages((len(dirMetaData) - 1), terminalHeightAvailable)

		// Avoid 1/0 pages
		if serverPageCountMax == 0 {
			serverPageCountMax = 1
		}

		pageStatus = fmt.Sprintf("page (%d / %d)", downloadBrowserIndex+1, serverPageCountMax)
		results := "result"
		resultCount := len(dirMetaData) - 1
		if navPtr == &rootDir {
			resultCount++
		}
		if resultCount != 1 {
			results += "s"
		}
		currentFilterInfo := fmt.Sprintf(" Current Filter (X to clear): \"%s\" | (%d %s) ", currentDirFilter, resultCount, results)
		baseFilterInfo := fmt.Sprintf(" Current Filter (X to clear): \"\" | (%d %s) ", resultCount, results)

		// check if screen overflow and truncate displayed search string if need be
		if terminalWidthAvailable-len(currentFilterInfo) < 0 {
			paddingSize := 2
			currentDirFilterTmp := currentDirFilter
			currentDirFilterTmp = truncateStringTo(currentDirFilterTmp, len(baseFilterInfo)+paddingSize, terminalWidthAvailable)
			currentFilterInfo = fmt.Sprintf(" Current Filter (X to clear): \"%s\" | (%d %s) ", currentDirFilterTmp, resultCount, results)
			currentFilterInfo = "=" + currentFilterInfo + "="
		}
		divider := ""
		for i := 0; i < (terminalWidthAvailable-len(currentFilterInfo))/2; i++ {
			divider += "="
		}
		currentFilterStringStatus = divider + currentFilterInfo + divider
	}

	//Render download menu
	ClearScreen()
	fmt.Println(divider)
	fmt.Println(menuHeader)
	fmt.Println(notificationBar)
	renderMetaData(dirMetaData, terminalHeightAvailable, terminalWidthAvailable)
	fmt.Println(currentFilterStringStatus)
	fmt.Println(pageStatus)
	fmt.Println("<<F <B | N> L>>")

	downloadNotification = ""

	return dirMetaData
}

func renderMetaData(directoryMetaData map[int]map[string]bool, terminalHeightAvailable int, terminalWidthAvailable int) {
	verticalHeightBuffer := terminalHeightAvailable
	if len(directoryMetaData) == 0 {
		fmt.Println("[EMPTY SERVER DMSG-HTTP-SERVER RESPONSE OR NO INDEX FILE FOUND OR NO SEARCH RESULTS IN ROOT DIR]")
	} else {

		if isMetaDataSorted(directoryMetaData) { // if the meta data is sorted
			for index := 1 + downloadBrowserIndex*terminalHeightAvailable; index <= len(directoryMetaData); index++ {

				for entryName := range directoryMetaData[index] {
					// if entry is a directory
					if directoryMetaData[index][entryName] {
						drawDirEntry(entryName, terminalWidthAvailable, index)

					} else {
						drawFileEntry(entryName, terminalWidthAvailable, index)
					}

					verticalHeightBuffer--
					if verticalHeightBuffer == 0 {
						goto END
					}

				}
			}
		} else { // if the meta data is unsorted
			metaKeys := []int{}
			for key := range directoryMetaData {
				metaKeys = append(metaKeys, key)
			}
			sort.Ints(metaKeys)
			for index := range metaKeys {
				indexOffset := index + downloadBrowserIndex*terminalHeightAvailable

				if indexOffset < len(metaKeys) {
					filteredIndexValue := metaKeys[indexOffset]
					for fileName, isDir := range directoryMetaData[filteredIndexValue] {
						if isDir {
							drawDirEntry(fileName, terminalWidthAvailable, filteredIndexValue)
						} else {
							drawFileEntry(fileName, terminalWidthAvailable, filteredIndexValue)
						}
						verticalHeightBuffer--
						if verticalHeightBuffer == 0 {
							goto END
						}
					}
				}

			}

		}

	}

	//vertical buffer
	for ; verticalHeightBuffer > 0; verticalHeightBuffer-- {
		fmt.Println("-")
	}
END:
}

func drawDirEntry(entryName string, terminalWidthAvailable int, index int) {
	tmpBasicLineEntry := fmt.Sprintf("%d)  / Directory", index)
	tmpLineEntry := tmpBasicLineEntry + entryName
	horizontalFill := ""
	//detect and format dirs with names that will overflow current terminal width

	if len(tmpLineEntry) > terminalWidthAvailable {
		entryName = truncateStringTo(entryName, len(tmpBasicLineEntry), terminalWidthAvailable)
	}
	for i := terminalWidthAvailable - len(tmpBasicLineEntry) - len(entryName); i > 0; i-- {
		horizontalFill += "-"
	}

	lineEntry := fmt.Sprintf("%d) %s/ %s Directory", index, entryName, horizontalFill)

	fmt.Println(lineEntry)
}

func drawFileEntry(entryName string, terminalWidthAvailable int, index int) {
	//if entry is a file
	fileSize := navPtr.files[entryName]
	fileSizeUnits := ""

	// format file sizes for human readability
	if fileSize > 1e9 {
		fileSize /= 1e9
		fileSizeUnits = "GB"
	} else if fileSize > 1e6 {
		fileSize /= 1e6
		fileSizeUnits = "MB"
	} else if fileSize > 1e3 {
		fileSize /= 1e3
		fileSizeUnits = "KB"
	} else {
		fileSizeUnits = "B"
	}

	//determine fill amount required
	entryName = strings.ReplaceAll(entryName, "–", "-") //replace em dash with regular dash em dash doesn't render correctly
	tmpBasicLineEntry := fmt.Sprintf("%d)   %.2f %s", index, fileSize, fileSizeUnits)
	horizontalFill := ""

	if len(tmpBasicLineEntry)+len(entryName) > terminalWidthAvailable {
		entryName = truncateStringTo(entryName, len(tmpBasicLineEntry), terminalWidthAvailable)
	}

	for i := terminalWidthAvailable - len(entryName+tmpBasicLineEntry); i > 0; i-- {
		horizontalFill += "-"
	}

	//draw line
	lineEntry := fmt.Sprintf("%d) %s %s %.2f %s", index, entryName, horizontalFill, fileSize, fileSizeUnits)
	fmt.Println(lineEntry)
}

func calcNumberOfPages(numberOfLineItems int, LineItemsAvailable int) int {
	return int(math.Ceil(float64(numberOfLineItems) / float64(LineItemsAvailable)))
}
