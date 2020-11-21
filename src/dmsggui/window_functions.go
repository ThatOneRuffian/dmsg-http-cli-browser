package dmsggui

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
)

//mainMenuPageCountMax holds the max page of the server list
var mainMenuPageCountMax int

//serverPageCountMax holds the max page given the current server's index and terminal height
var serverPageCountMax int

//SavedServers stores server cache - initalized on loadCache
var SavedServers map[int][2]string

//currentServerIndexContents will store the parsed server index values
var currentServerIndexContents map[int][2]string

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func refreshServerIndex(serverPublicKey string, clearCache bool) {
	ClearScreen()
	DownloadBrowserIndex = 0
	if clearCache {
		clearServerIndexFile(serverPublicKey)
		fmt.Println("Downloading Server Index...")
		dmsggetWrapper(serverPublicKey, indexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
	if loadServerIndex(serverPublicKey) {

	} else {
		clearServerIndexFile(serverPublicKey)
		fmt.Println("Downloading Server Index...")
		dmsggetWrapper(serverPublicKey, indexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
}

func renderServerBrowser() {
	bufferHeight := 7 //lines consumed by menu elements
	dirNumberOfItems := len(SavedServers)
	terminalHeightAvailable, heightError := sttyWrapperGetTerminalHeight()
	terminalWidth, widthError := sttyWrapperGetTerminalWidth()

	if heightError != nil || widthError != nil {
		fmt.Println("Error fetching terminal dimensions")
		fmt.Println(heightError)
		fmt.Println(widthError)
		terminalHeightAvailable = 10 //default on error
		terminalWidth = 20

	} else {
		terminalHeightAvailable -= bufferHeight
	}

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
	for i := 0; i < terminalWidth; i++ {
		divider += "="
	}

	//Render variables
	titleBuffer := ""
	menuTitle := "SERVER DOWNLOAD INDEX"
	//dirMetaData := getCurrentDirMetaData()
	currentDir := getPresentWorkingDirectory()
	tmpTitle := fmt.Sprintf("%s%s", menuTitle, currentDir)
	titleBufferLength := terminalWidth - len(tmpTitle)
	for i := 0; i < titleBufferLength; i++ {
		titleBuffer = titleBuffer + " "
	}
	//menuHeader := fmt.Sprintf("%s%s%s", menuTitle, titleBuffer, currentDir)
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex+1, mainMenuPageCountMax)

	ClearScreen()
	fmt.Println(divider)
	fmt.Println("DMSG-HTTP SERVER LIST")
	fmt.Println(divider)
	renderHomeMenuServerList(terminalHeightAvailable, terminalWidth)

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
		indexOffset := index + (terminalHeightAvailable * DownloadBrowserIndex)
		tmpLineEntry := fmt.Sprintf("%d) %s ", indexOffset+1, SavedServers[indexOffset][0])
		horizontalFill := ""
		for i := terminalWidthAvailable - len(tmpLineEntry); i > 0; i-- {
			horizontalFill += "-"
		}
		lineEntry := fmt.Sprintf("%d) %s %s", indexOffset+1, SavedServers[indexOffset][0], horizontalFill)

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
	terminalHeightAvailable, heightError := sttyWrapperGetTerminalHeight()
	terminalWidth, widthError := sttyWrapperGetTerminalWidth()

	if heightError != nil || widthError != nil {
		fmt.Println("Error fetching terminal dimensions")
		fmt.Println(heightError)
		fmt.Println(widthError)
		terminalHeightAvailable = 10 //default on error
		terminalWidth = 20

	} else {
		terminalHeightAvailable -= bufferHeight
	}

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
	for i := 0; i < terminalWidth; i++ {
		divider += "="
	}

	//Render variables
	titleBuffer := ""
	menuTitle := "SERVER DOWNLOAD INDEX"
	dirMetaData := getCurrentDirMetaData()
	currentDir := getPresentWorkingDirectory()
	tmpTitle := fmt.Sprintf("%s%s", menuTitle, currentDir)
	titleBufferLength := terminalWidth - len(tmpTitle)
	for i := 0; i < titleBufferLength; i++ {
		titleBuffer = titleBuffer + " "
	}
	menuHeader := fmt.Sprintf("%s%s%s", menuTitle, titleBuffer, currentDir)
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex+1, serverPageCountMax)

	//Render download menu
	ClearScreen()
	fmt.Println(divider)
	fmt.Println(menuHeader)
	fmt.Println(divider)
	renderMetaData(dirMetaData, terminalHeightAvailable, terminalWidth)
	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<<F <B | N> L>>")

	return dirMetaData
}

func renderMetaData(directoryMetaData map[int]map[string]bool, terminalHeightAvailable int, terminalWidthAvailable int) {
	verticalHeightBuffer := terminalHeightAvailable

	for index := 1 + DownloadBrowserIndex*terminalHeightAvailable; index <= len(directoryMetaData); index++ {
		for key := range directoryMetaData[index] {
			// if entry is a directory
			if directoryMetaData[index][key] {
				tmpLineEntry := fmt.Sprintf("%d) %s / Directory", index, key)
				horizontalFill := ""
				for i := terminalWidthAvailable - len(tmpLineEntry); i > 0; i-- {
					horizontalFill += "-"
				}
				lineEntry := fmt.Sprintf("%d) %s/ %s Directory", index, key, horizontalFill)

				fmt.Println(lineEntry)
			} else {
				tmpLineEntry := fmt.Sprintf("%d) %s  %.2f MB", index, key, navPtr.files[key]/1e6)
				horizontalFill := ""
				for i := terminalWidthAvailable - len(tmpLineEntry); i > 0; i-- {
					horizontalFill += "-"
				}
				lineEntry := fmt.Sprintf("%d) %s %s %.2f MB", index, key, horizontalFill, navPtr.files[key]/1e6)

				fmt.Println(lineEntry)
			}
			verticalHeightBuffer--
			if verticalHeightBuffer == 0 {
				goto END
			}

		}

	}
	for ; verticalHeightBuffer > 0; verticalHeightBuffer-- {
		fmt.Println("-")
	}
END:
}

func getCurrentDirMetaData() map[int]map[string]bool {
	var subDirKeys []string
	var fileNames []string

	returnValue := make(map[int]map[string]bool)
	swapDir := make(map[string]bool)

	//dump dir names in current dir
	if navPtr != &rootDir {
		subDirKeys = append(subDirKeys, "..")
	}

	for key := range navPtr.subDirs {
		subDirKeys = append(subDirKeys, key)
	}

	//dump files names in current dir
	for key := range navPtr.files {
		fileNames = append(fileNames, key)
	}

	sort.Strings(fileNames)
	sort.Strings(subDirKeys)

	//merge metadata
	for key, value := range subDirKeys {
		swapDir[value] = true
		returnValue[key+1] = swapDir
		swapDir = make(map[string]bool)
	}
	for key, value := range fileNames {
		swapDir[value] = false
		returnValue[key+1+len(subDirKeys)] = swapDir
		swapDir = make(map[string]bool)
	}
	return returnValue
}
