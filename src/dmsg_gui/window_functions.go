package dmsg_gui

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
)

//ServerPageCountMax holds the max page give the current server's index
var ServerPageCountMax int

//DownloadBrowserIndex stores the current page of the download browsers
var DownloadBrowserIndex int = 1

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
	}

	fmt.Println("Downloading Server Index...")
	dmsggetWrapper(serverPublicKey, IndexDownloadLoc, "index", "index."+serverPublicKey, false)
	LoadServerIndex(serverPublicKey)
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

func RenderServerBrowser2() {
	currendDirPtr := &rootDir
	bufferHeight := 7
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
	ServerPageCountMax = len(CurrentServerIndex) / terminalHeight
	pageRemainder := len(CurrentServerIndex) % terminalHeight
	if pageRemainder > 0 {
		ServerPageCountMax++
	}
	// Avoid 1/0 pages
	if ServerPageCountMax == 0 {
		ServerPageCountMax = 1
	}
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex, ServerPageCountMax)

	divider := ""
	for i := 0; i < terminalWidth; i++ {
		divider += "="
	}
	ClearScreen()
	fmt.Println(divider)
	fmt.Println("SERVER DOWNLOAD INDEX")
	fmt.Println(divider)
	currendDirPtr = currendDirPtr.subDirs["Cowboy_Bebop"]
	renderDirectories(currendDirPtr, terminalWidth)
	renderFiles(currendDirPtr, terminalWidth)
	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<< B  |  N >>")
}

func renderDirectories(dirPtr *Directory, terminalWidth int) {
	//sort sub dirs A-Z
	var subDirKeys []string
	for key := range dirPtr.subDirs {
		subDirKeys = append(subDirKeys, key)
	}
	sort.Strings(subDirKeys)

	//Render directories
	renderIndex := 1
	if len(subDirKeys) > 0 {
		for _, key := range subDirKeys {
			itemIndex := renderIndex + (terminalWidth*DownloadBrowserIndex - 1) - terminalWidth + 1
			if itemIndex-1 < len(dirPtr.subDirs) {
				fileBuffer := ""
				tmpEntry := fmt.Sprintf("%d) %s%s", itemIndex, key, "Directory")
				bufferToAdd := terminalWidth - len(tmpEntry)
				for i := 0; i < bufferToAdd; i++ {
					fileBuffer += "-"
				}
				listEntry := fmt.Sprintf("%d) %s%s%s", itemIndex, key, fileBuffer, "Directory")
				fmt.Println(listEntry)
			} else {
				fmt.Println("-")
			}
			renderIndex++
		}
	} else {
		//fmt.Println("[Empty server index/Could not fetch list]")
	}
}
func renderFiles(dirPtr *Directory, terminalWidth int) {
	//sort sub dirs A-Z
	var fileNames []string
	for key := range dirPtr.files {
		fileNames = append(fileNames, key)
	}
	sort.Strings(fileNames)

	//Render directories
	itemIndex := 1
	if len(fileNames) > 0 {
		for _, key := range fileNames {
			if true {
				fileSize := dirPtr.files[key]
				tmpEntry := fmt.Sprintf("%d) %s%d", itemIndex, key, fileSize)
				bufferToAdd := terminalWidth - len(tmpEntry)
				fileBuffer := ""
				for i := 0; i < bufferToAdd; i++ {
					fileBuffer += "-"
				}
				listEntry := fmt.Sprintf("%d) %s%s%d", itemIndex, key, fileBuffer, fileSize)
				fmt.Println(listEntry)
			} else {
				fmt.Println("-")
			}
			itemIndex++
		}
	} else {
		fmt.Println("[Empty Dir]")
	}
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
