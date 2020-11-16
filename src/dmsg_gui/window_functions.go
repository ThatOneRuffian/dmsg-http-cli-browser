package dmsg_gui

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
)

//ServerPageCountMax holds the max page give the current server's index
var ServerPageCountMax int

//SavedServers stores server cache - initalized on loadCache
var SavedServers map[int][2]string

//CurrentServerIndex will store the parsed server index values
var CurrentServerIndex map[int][2]string

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func refreshServerIndex(serverPublicKey string, clearCache bool) {
	ClearScreen()
	DownloadBrowserIndex = 1
	if clearCache {
		ClearServerIndexFile(serverPublicKey)
		fmt.Println("Downloading Server Index...")
		dmsggetWrapper(serverPublicKey, IndexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
	if LoadServerIndex(serverPublicKey) {

	} else {
		ClearServerIndexFile(serverPublicKey)
		fmt.Println("Downloading Server Index...")
		dmsggetWrapper(serverPublicKey, IndexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
}

func renderServerBrowser() {
	pageStatus := fmt.Sprintf("page (%d / %d)", 1, 20)
	divider := "----------------------"
	ClearScreen()

	fmt.Println(divider)
	fmt.Println("DMSG HTTP SERVER LIST")
	fmt.Println(divider)

	for i := 0; i < len(SavedServers); i++ {
		listEntry := fmt.Sprintf("%d) %s", i+1, SavedServers[i][0])
		fmt.Println(listEntry)
	}

	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<< B  |  N >>")
}

func renderServerBrowser2() map[int]map[string]bool {

	bufferHeight := 7
	terminalHeightAvailable, err := sttyWrapperGetTerminalHeight()
	if err != nil {
		terminalHeightAvailable = 10 //default on error

	} else {
		terminalHeightAvailable -= bufferHeight

	}
	terminalWidth, err := sttyWrapperGetTerminalWidth()
	if err != nil {
		fmt.Println(err)
	}
	dirNumberOfItems := (len(navPtr.subDirs) + len(navPtr.files))
	ServerPageCountMax = dirNumberOfItems / terminalHeightAvailable
	pageRemainder := dirNumberOfItems % terminalHeightAvailable

	if pageRemainder > 0 {
		ServerPageCountMax++
	}

	// Avoid 1/0 pages
	if ServerPageCountMax == 0 {
		ServerPageCountMax = 1
	}
	//Create header divider of appropriate length
	divider := ""
	for i := 0; i < terminalWidth; i++ {
		divider += "="
	}

	//Render logic
	ClearScreen()
	fmt.Println(divider)
	menuTitle := "SERVER DOWNLOAD INDEX"
	currentDir := getPresentWorkingDirectory()
	tmpTitle := fmt.Sprintf("%s%s", menuTitle, currentDir)
	titleBufferLength := terminalWidth - len(tmpTitle)
	titleBuffer := ""

	for i := 0; i < titleBufferLength; i++ {
		titleBuffer = titleBuffer + " "
	}
	fmt.Println(fmt.Sprintf("%s%s%s", menuTitle, titleBuffer, currentDir))
	fmt.Println(divider)

	dirMetaData := getCurrentDirMetaData()
	renderMetaData(dirMetaData, terminalHeightAvailable)
	fmt.Println(divider)
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex, ServerPageCountMax)
	fmt.Println(pageStatus)
	fmt.Println("<< B  |  N >>")
	return dirMetaData
}
func renderMetaData(directoryMetaData map[int]map[string]bool, terminalHeightAvailable int) {
	verticalHeightBuffer := terminalHeightAvailable
	for index := 1; index < len(directoryMetaData) && index <= terminalHeightAvailable; index++ {
		for key := range directoryMetaData[index] {
			if directoryMetaData[index][key] {
				lineEntry := fmt.Sprintf("%d) %s/", index, key)
				fmt.Println(lineEntry)
			} else {
				lineEntry := fmt.Sprintf("%d) %s", index, key)
				fmt.Println(lineEntry)
			}
			verticalHeightBuffer--

		}

	}

	for ; verticalHeightBuffer > 0; verticalHeightBuffer-- {
		fmt.Println("-")
	}
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

func renderServerIndexBrowser() {
	bufferHeight := 7
	DownloadListLength, err := sttyWrapperGetTerminalHeight()
	if err != nil {
		DownloadListLength = 10 //default on error

	} else {
		DownloadListLength = DownloadListLength - bufferHeight

	}
	ServerPageCountMax = len(CurrentServerIndex) / DownloadListLength
	pageRemainder := len(CurrentServerIndex) % DownloadListLength
	if pageRemainder > 0 {
		ServerPageCountMax++
	}
	// Avoid 1/0 pages
	if ServerPageCountMax == 0 {
		ServerPageCountMax = 1
	}
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex, ServerPageCountMax)
	terminalWidth, err := sttyWrapperGetTerminalWidth()
	if err != nil {
		fmt.Println(err)
	}
	divider := ""
	for i := 0; i < terminalWidth; i++ {
		divider += "="
	}
	ClearScreen()
	fmt.Println(divider)
	fmt.Println("SERVER DOWNLOAD INDEX")
	fmt.Println(divider)
	renderIndex := 1
	if len(CurrentServerIndex) > 0 {
		for ; renderIndex <= DownloadListLength; renderIndex++ {
			itemIndex := renderIndex + (DownloadListLength*DownloadBrowserIndex - 1) - DownloadListLength + 1
			if itemIndex-1 < len(CurrentServerIndex) {
				fileBuffer := ""
				tmpEntry := fmt.Sprintf("%d) %s%s", itemIndex, CurrentServerIndex[itemIndex-1][0], fmt.Sprintf(CurrentServerIndex[itemIndex-1][1]))
				bufferToAdd := terminalWidth - len(tmpEntry)
				for i := 0; i < bufferToAdd; i++ {
					fileBuffer += "-"
				}
				listEntry := fmt.Sprintf("%d) %s%s%s", itemIndex, CurrentServerIndex[itemIndex-1][0], fileBuffer, fmt.Sprintf(CurrentServerIndex[itemIndex-1][1]))

				fmt.Println(listEntry)
			} else {
				fmt.Println("-")
			}
		}
	} else {
		fmt.Println("[Empty server index/Could not fetch list]")
	}

	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<< B  |  N >>")
}
