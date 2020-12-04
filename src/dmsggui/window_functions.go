package dmsggui

import (
	"fmt"
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

func refreshServerIndex(serverPublicKey string, clearCache bool) {
	ClearScreen()
	resetDownLoadPageIndex()
	if clearCache {
		fmt.Println("Downloading Server Index...")
		clearServerIndexFile(serverPublicKey)
		dmsggetWrapper(serverPublicKey, indexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
	if loadServerIndex(serverPublicKey) {

	} else {
		ClearScreen()
		fmt.Println("Downloading Server Index...")
		clearServerIndexFile(serverPublicKey)
		dmsggetWrapper(serverPublicKey, indexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
	assembleFileStructure(serverPublicKey)
}

func renderServerBrowser() {
	bufferHeight := 7 //lines consumed by menu elements
	dirNumberOfItems := len(SavedServers)
	terminalHeightAvailable, terminalWidthAvailable := getTerminalDims(bufferHeight)
	mainMenuPageCountMax = dirNumberOfItems / terminalHeightAvailable
	pageRemainder := dirNumberOfItems % terminalHeightAvailable

	// add additional page to fit remaining line items
	if pageRemainder > 0 {
		mainMenuPageCountMax++
	}

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

func renderServerDownloadList() map[int]map[string]bool {

	bufferHeight := 7 //lines consumed by menu elements
	dirNumberOfItems := len(navPtr.subDirs) + len(navPtr.files)
	terminalHeightAvailable, terminalWidthAvailable := getTerminalDims(bufferHeight)

	serverPageCountMax = dirNumberOfItems / terminalHeightAvailable
	pageRemainder := dirNumberOfItems % terminalHeightAvailable

	// add additional page to fit remaining line items
	if pageRemainder > 0 {
		serverPageCountMax++
	}

	// Avoid 1/0 pages
	if serverPageCountMax == 0 {
		serverPageCountMax = 1
	}

	//Create header divider of appropriate length
	divider := ""
	for i := 0; i < terminalWidthAvailable; i++ {
		divider += "="
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
	if len(currentDirFilter) == 0 {
		pageStatus = fmt.Sprintf("page (%d / %d)", downloadBrowserIndex+1, serverPageCountMax)
		currentFilterStringStatus = divider
	} else {
		serverPageCountMax = (len(dirMetaData) - 1) / terminalHeightAvailable
		pageRemainder := dirNumberOfItems % terminalHeightAvailable

		// add additional page to fit remaining line items
		if pageRemainder > 0 {
			serverPageCountMax++
		}

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
		if resultCount > 1 {
			results += "s"
		}
		currentFilterInfo := fmt.Sprintf(" Current Filter (X to clear): \"%s\" | (%d %s) ", currentDirFilter, resultCount, results)
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
	fmt.Println(divider)
	renderMetaData(dirMetaData, terminalHeightAvailable, terminalWidthAvailable)
	fmt.Println(currentFilterStringStatus)
	fmt.Println(pageStatus)
	fmt.Println("<<F <B | N> L>>")

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

func isMetaDataSorted(directoryMetaData map[int]map[string]bool) bool {
	metaDataLength := len(directoryMetaData)
	isSorted := true
	for index := range directoryMetaData {
		if index > metaDataLength {
			isSorted = false
		}
	}
	return isSorted
}

func drawDirEntry(entryName string, terminalWidthAvailable int, index int) {
	tmpBasicLineEntry := fmt.Sprintf("%d)  / Directory", index)
	tmpLineEntry := tmpBasicLineEntry + entryName
	horizontalFill := ""
	//detect and format dirs with names that will overflow current terminal width

	if len(tmpLineEntry) >= terminalWidthAvailable {
		entryName = truncateStringTo(entryName, len(tmpBasicLineEntry), terminalWidthAvailable)
	}
	for i := terminalWidthAvailable - len(tmpLineEntry); i > 0; i-- {
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

	if len(tmpBasicLineEntry)+len(entryName) >= terminalWidthAvailable {
		entryName = truncateStringTo(entryName, len(tmpBasicLineEntry), terminalWidthAvailable)
	}

	for i := terminalWidthAvailable - len(entryName+tmpBasicLineEntry); i > 0; i-- {
		horizontalFill += "-"
	}

	//draw line
	lineEntry := fmt.Sprintf("%d) %s %s %.2f %s", index, entryName, horizontalFill, fileSize, fileSizeUnits)
	fmt.Println(lineEntry)
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

func getCurrentDirMetaData() map[int]map[string]bool {
	var subDirKeys []string
	var fileNames []string
	returnValue := make(map[int]map[string]bool)
	swapDir := make(map[string]bool)
	lengthOfFilterList := len(currentDirFilter)
	//dump dir/file keys in current dir and sort A-Z
	if navPtr != &rootDir {
		//add key for directory back
		subDirKeys = append(subDirKeys, "..")
	}

	for key := range navPtr.subDirs {
		//append subdir keys
		subDirKeys = append(subDirKeys, key)
	}

	for key := range navPtr.files {
		//append file keys
		fileNames = append(fileNames, key)
	}

	sort.Strings(fileNames)
	sort.Strings(subDirKeys)

	//merge metadata
	for key, value := range subDirKeys {
		if lengthOfFilterList > 0 {
			if strings.Contains(strings.ToUpper(value), strings.ToUpper(string(currentDirFilter))) || strings.Contains(strings.ToUpper(value), "..") {
				swapDir[value] = true
				returnValue[key+1] = swapDir
				swapDir = make(map[string]bool)
			}
		} else if lengthOfFilterList == 0 {
			swapDir[value] = true
			returnValue[key+1] = swapDir
			swapDir = make(map[string]bool)
		}
	}

	for key, value := range fileNames {
		if strings.Contains(strings.ToUpper(value), strings.ToUpper(string(currentDirFilter))) {
			swapDir[value] = false
			returnValue[key+1+len(subDirKeys)] = swapDir
			swapDir = make(map[string]bool)
		} else if lengthOfFilterList == 0 {
			swapDir[value] = false
			returnValue[key+1+len(subDirKeys)] = swapDir
			swapDir = make(map[string]bool)
		}
	}
	return returnValue
}
