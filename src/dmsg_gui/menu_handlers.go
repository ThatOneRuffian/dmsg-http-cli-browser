package dmsg_gui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ServerListMainMenu() string {
	serverPublicKey := ""
	consoleInput := bufio.NewReader(os.Stdin)
ServerMenu:

	renderServerBrowser()

	fmt.Print("(Press A to Add server, D to Delete a server, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(removeNewline(userChoice))
	switch userChoice {
	case "Q":
		ClearScreen()
		os.Exit(1)
	case "A":
		ClearScreen()
		addServer()
		LoadCache()
	case "B":
		//TODO
	case "N":
		//TODO
	case "D":
		DeleteServerWizard()
	default:
		userInt, err := strconv.Atoi(userChoice)
		if err != nil {
			break
		}
		if userInt >= 1 && userInt <= len(SavedServers) {
			serverPublicKey = SavedServers[userInt-1][1]
			refreshServerIndex(serverPublicKey, false)
			goto ExitLoop
		} else {
			break
		}
	}
	goto ServerMenu
ExitLoop:
	return serverPublicKey
}

func ServerIndexMenuHandler(serverPublicKey string) {
	consoleInput := bufio.NewReader(os.Stdin)
ServerIndexMenu:

	renderServerIndexBrowser()

	fmt.Print("(R to Refresh Server Index, E to Exit Server File Browser, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(removeNewline(userChoice))
	switch userChoice {
	case "Q":
		ClearScreen()
		os.Exit(1)
	case "E":
		goto ExitLoop
	case "B":
		if DownloadBrowserIndex > 1 {
			DownloadBrowserIndex--
		}

		goto ServerIndexMenu
	case "N":
		if DownloadBrowserIndex < ServerPageCountMax {

			DownloadBrowserIndex++
		}

		goto ServerIndexMenu
	case "R":
		refreshServerIndex(serverPublicKey, true)
		goto ServerIndexMenu

	default:
		userInt, err := strconv.Atoi(userChoice)
		if err != nil {
			break
		}
		if userInt >= 1 && userInt <= len(CurrentServerIndex) {
			filenameDownload := CurrentServerIndex[userInt-1][0]
			// download file
			ClearScreen()

			dmsggetWrapper(serverPublicKey, MainDownloadsLoc, filenameDownload, "", true)
		} else {
			break
		}
	}
	goto ServerIndexMenu

ExitLoop:
}
