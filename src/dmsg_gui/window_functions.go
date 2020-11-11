package dmsg_gui

import (
	"fmt"
	"os"
	"os/exec"
)

// ServerPageCountMax holds the max page give the current server's index
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

func RefreshServerIndex(serverPublicKey string, clearCache bool) {
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
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex, ServerPageCountMax)
	terminalWidth, err := sttyWrapperGetTerminalWidth()
	if err != nil {
		fmt.Println(err)
	}
	divider := ""
	for i := 0; i < terminalWidth; i++ {
		divider += "-"
	}
	//divider := "------------------------------------------"
	ClearScreen()
	fmt.Println(divider)
	fmt.Println("SERVER DOWNLOAD INDEX")
	fmt.Println(divider)
	renderIndex := 1
	if len(CurrentServerIndex) > 0 {
		for ; renderIndex <= DownloadListLength; renderIndex++ {
			itemIndex := renderIndex + (DownloadListLength*DownloadBrowserIndex - 1) - DownloadListLength + 1
			if itemIndex-1 < len(CurrentServerIndex) {
				listEntry := fmt.Sprintf("%d) %s\t\t\t%s", itemIndex, CurrentServerIndex[itemIndex-1][0], fmt.Sprintf(CurrentServerIndex[itemIndex-1][1]))
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
