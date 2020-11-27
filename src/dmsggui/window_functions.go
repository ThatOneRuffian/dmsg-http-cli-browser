package dmsggui

import (
	"fmt"
	"os"
	"os/exec"
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

const defaultTerminalWidth = 85

const defaultTerminalHeight = 30

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func refreshServerIndex(serverPublicKey string, clearCache bool) {
	ClearScreen()
	downloadBrowserIndex = 0
	if clearCache {
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
		terminalHeightAvailable = defaultTerminalHeight //default on error
		terminalWidth = defaultTerminalWidth

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
	currentDir := getPresentWorkingDirectory()
	tmpTitle := fmt.Sprintf("%s%s", menuTitle, currentDir)
	titleBufferLength := terminalWidth - len(tmpTitle)
	for i := 0; i < titleBufferLength; i++ {
		titleBuffer = titleBuffer + " "
	}
	pageStatus := fmt.Sprintf("page (%d / %d)", mainMenuBrowserIndex+1, mainMenuPageCountMax)

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
	terminalHeightAvailable, heightError := sttyWrapperGetTerminalHeight()
	terminalWidth, widthError := sttyWrapperGetTerminalWidth()

	if heightError != nil || widthError != nil {
		fmt.Println("Error fetching terminal dimensions")
		fmt.Println(heightError)
		fmt.Println(widthError)
		terminalHeightAvailable = defaultTerminalHeight //default on error
		terminalWidth = defaultTerminalWidth

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
	pageStatus := fmt.Sprintf("page (%d / %d)", downloadBrowserIndex+1, serverPageCountMax)

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

	for index := 1 + downloadBrowserIndex*terminalHeightAvailable; index <= len(directoryMetaData); index++ {
		for entryName := range directoryMetaData[index] {
			// if entry is a directory
			if directoryMetaData[index][entryName] {
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
			} else {
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
				tmpLineEntry := fmt.Sprintf("%d) %s  %.2f %s", index, entryName, fileSize, fileSizeUnits)
				horizontalFill := ""
				for i := terminalWidthAvailable - len(tmpLineEntry); i > 0; i-- {
					horizontalFill += "-"
				}

				//draw line
				lineEntry := fmt.Sprintf("%d) %s %s %.2f %s", index, entryName, horizontalFill, fileSize, fileSizeUnits)
				fmt.Println(lineEntry)
			}
			verticalHeightBuffer--
			if verticalHeightBuffer == 0 {
				goto END
			}

		}

	}
	//vertical buffer
	for ; verticalHeightBuffer > 0; verticalHeightBuffer-- {
		fmt.Println("-")
	}
END:
}

func truncateStringTo(stringToTruncate string, sizeToFit int, terminalWidthAvailable int) string {
	overflowAmount := terminalWidthAvailable - sizeToFit
	truncateBuffer := overflowAmount / 2 //amount to keep on beginning and end of string
	entryNameLength := len(stringToTruncate)
	if entryNameLength > overflowAmount {
		//add fill here

		if overflowAmount-truncateBuffer*2 >= 0 {
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
