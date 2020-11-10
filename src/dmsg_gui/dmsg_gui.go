package dmsg_gui

import (
	"fmt"
	"os"
	"os/exec"
)

//DownloadBrowserIndex stores the current page of the download browsers
var DownloadBrowserIndex int = 1

//IndexDownloadLoc is where the active server's index is downloaded
var IndexDownloadLoc string = "/tmp/"

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

func RenderServerBrowser() {
	pageStatus := fmt.Sprintf("page (%d / %d)", 1, 20)
	divider := "----------------------"
	ClearScreen()

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

func RenderServerIndexBrowser() {
	bufferHeight := 7
	DownloadListLength, err := SttyWrapperGetTerminalHeight()
	if err != nil {
		DownloadListLength = 10 //default on error

	} else {
		DownloadListLength = DownloadListLength - bufferHeight

	}
	pageCountMax := len(CurrentServerIndex) / DownloadListLength
	pageRemainder := len(CurrentServerIndex) % DownloadListLength
	if pageRemainder > 0 {
		pageCountMax++
	}
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex, pageCountMax)
	divider := "------------------------------------------"
	ClearScreen()
	fmt.Println(divider)
	fmt.Println("SERVER DOWNLOAD INDEX")
	fmt.Println(divider)
	renderIndex := 1
	for ; renderIndex <= DownloadListLength; renderIndex++ {
		itemIndex := renderIndex + (DownloadListLength*DownloadBrowserIndex - 1) - DownloadListLength + 1
		if itemIndex-1 < len(CurrentServerIndex) {
			listEntry := fmt.Sprintf("%d) %s\t\t\t%s", itemIndex, CurrentServerIndex[itemIndex-1][0], fmt.Sprintf(CurrentServerIndex[itemIndex-1][1]))
			fmt.Println(listEntry)
		}
	}

	fmt.Println(divider)
	fmt.Println(pageStatus)
	fmt.Println("<< B  |  N >>")
}
