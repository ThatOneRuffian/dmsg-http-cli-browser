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
	currendDirPtr := navPtr
	bufferHeight := 8
	terminalHeight, err := sttyWrapperGetTerminalHeight()
	if err != nil {
		terminalHeight = 10 //default on error

	} else {
		terminalHeight = terminalHeight - bufferHeight

	}
	terminalWidth, err := sttyWrapperGetTerminalWidth()
	if err != nil {
		fmt.Println(err)
	}
	dirNumberOfItems := (len(navPtr.subDirs) + len(navPtr.files))
	ServerPageCountMax = dirNumberOfItems / terminalHeight
	pageRemainder := dirNumberOfItems % terminalHeight
	if pageRemainder > 0 {
		ServerPageCountMax++
	}
	// Avoid 1/0 pages
	if ServerPageCountMax == 0 {
		ServerPageCountMax = 1
	}

	divider := ""
	for i := 0; i < terminalWidth; i++ {
		divider += "="
	}
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
	dirIndexMetaData, indexStartValue, remainingHeight := renderDirectories(currendDirPtr, terminalWidth, terminalHeight)
	fileIndexMetaData, remainingHeight := renderFiles(currendDirPtr, terminalWidth, remainingHeight, indexStartValue)
	//fmt.Println("test", terminalHeight, remainingHeight)

	//merge metadata
	metaData := dirIndexMetaData
	for key, value := range fileIndexMetaData {
		metaData[key] = value
	}
	metaDataLength := len(metaData)
	if metaDataLength > 0 {
		for i := 0; i <= remainingHeight; i++ {
			fmt.Println("-")
		}
	} else {
		fmt.Println("[Empty]")
		for i := 0; i < remainingHeight-metaDataLength; i++ {
			fmt.Println("-")
		}
	}

	fmt.Println(divider)
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex, ServerPageCountMax)
	fmt.Println(pageStatus)
	fmt.Println("<< B  |  N >>")
	return metaData
}

func renderDirectories(dirPtr *Directory, terminalWidth int, verticalHeightBuffer int) (map[int]map[string]bool, int, int) {
	//sort sub dirs A-Z
	var subDirKeys []string
	if dirPtr != &rootDir {
		subDirKeys = append(subDirKeys, "..")
	}

	for key := range dirPtr.subDirs {
		subDirKeys = append(subDirKeys, key)
	}
	sort.Strings(subDirKeys)

	//Render directories
	renderIndex := 1
	verticalHeightSwap := verticalHeightBuffer
	if len(subDirKeys) > 0 {
		for i := 0; i < len(subDirKeys) && i <= verticalHeightBuffer; i++ {
			itemIndex := renderIndex
			if itemIndex-1 < len(subDirKeys) {
				fileBuffer := ""
				tmpEntry := fmt.Sprintf("%d) %s/%s", itemIndex, subDirKeys[i], "Directory")
				bufferToAdd := terminalWidth - len(tmpEntry)
				for i := 0; i < bufferToAdd; i++ {
					fileBuffer += "-"
				}
				listEntry := fmt.Sprintf("%d) %s/%s%s", itemIndex, subDirKeys[i], fileBuffer, "Directory")
				fmt.Println(listEntry)
			}
			renderIndex++
			verticalHeightSwap--
		}
	} else {
		//fmt.Println("[Empty server index/Could not fetch list]")
	}
	verticalHeightBuffer = verticalHeightSwap
	returnValue := make(map[int]map[string]bool)
	swapDir := make(map[string]bool)

	for key, value := range subDirKeys {
		swapDir[value] = true
		returnValue[key+1] = swapDir
		swapDir = make(map[string]bool)
	}

	return returnValue, len(returnValue) + 1, verticalHeightBuffer
}

func renderFiles(dirPtr *Directory, terminalWidth int, remainingHeight int, indexStartValue int) (map[int]map[string]bool, int) {
	//sort sub dirs A-Z
	var fileNames []string
	for key := range dirPtr.files {
		fileNames = append(fileNames, key)
	}
	sort.Strings(fileNames)

	//Render files

	itemIndex := indexStartValue
	swapHeight := remainingHeight
	if len(fileNames) > 0 {
		for i := 0; i < len(fileNames) && i <= remainingHeight; i++ {
			//fill horizontal space
			fileSize := dirPtr.files[fileNames[i]]
			tmpEntry := fmt.Sprintf("%d) %s%d", itemIndex, fileNames[i], fileSize)
			bufferToAdd := terminalWidth - len(tmpEntry)
			fileBuffer := ""
			for z := 0; z < bufferToAdd; z++ {
				fileBuffer += "-"
			}
			listEntry := fmt.Sprintf("%d) %s%s%d", itemIndex, fileNames[i], fileBuffer, fileSize)
			fmt.Println(listEntry)

			itemIndex++
			swapHeight--
		}
	} else {
		//fmt.Println("[No files in this dir]")
	}
	remainingHeight = swapHeight
	//format return value
	returnValue := make(map[int]map[string]bool)
	swapDir := make(map[string]bool)

	for key, value := range fileNames {
		swapDir[value] = false
		returnValue[indexStartValue+key] = swapDir
		swapDir = make(map[string]bool)
	}
	return returnValue, remainingHeight
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
