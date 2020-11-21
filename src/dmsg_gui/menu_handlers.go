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
		deleteServerWizard()
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
	navPtr = &rootDir
	assembleFileStructure(serverPublicKey)

ServerIndexMenu:
	directoryMetaData := renderServerDownloadList()
	consoleInput := bufio.NewReader(os.Stdin)
	fmt.Print("(R to Refresh Server Index, E to Exit Server File Browser, G to Goto page, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(removeNewline(userChoice))
	switch userChoice {
	case "Q":
		ClearScreen()
		os.Exit(1)
	case "E":
		//clear server root dir
		rootDir = directory{
			files:     make(map[string]float64),
			parentDir: nil,
			dirName:   "/",
			subDirs:   make(map[string]*directory),
		}
		goto ExitLoop
	case "B":
		if DownloadBrowserIndex > 0 {
			DownloadBrowserIndex--
		}

		goto ServerIndexMenu
	case "F":
		DownloadBrowserIndex = 0
		goto ServerIndexMenu
	case "G":
		fmt.Print("Enter page number:")
		consoleReader := bufio.NewReader(os.Stdin)
		pageNumber, _ := consoleReader.ReadString('\n')
		pageNumber = strings.ToUpper(removeNewline(pageNumber))
		pageNumber2, err := strconv.Atoi(pageNumber)

		if err != nil {

		} else if pageNumber2 > 0 && pageNumber2-1 < ServerPageCountMax {
			DownloadBrowserIndex = pageNumber2 - 1
		}

		goto ServerIndexMenu

	case "L":
		DownloadBrowserIndex = ServerPageCountMax - 1
		goto ServerIndexMenu

	case "N":
		if DownloadBrowserIndex < ServerPageCountMax-1 {

			DownloadBrowserIndex++
		}

		goto ServerIndexMenu
	case "R":
		refreshServerIndex(serverPublicKey, true)
		goto ServerIndexMenu

	default:
		userInputVar, err := strconv.Atoi(userChoice)
		if err != nil {
			break
		}
		if userInputVar >= 1 && userInputVar <= len(directoryMetaData) {

			//determine if the key is a dir or a file
			for index := range directoryMetaData[userInputVar] {
				// navigate up a directory
				if index == ".." && navPtr.parentDir != nil {
					navPtr = navPtr.parentDir
				} else {
					// if object is a directory then navigate into
					if directoryMetaData[userInputVar][index] {
						navPtr = navPtr.subDirs[index]
					} else {
						// runs once
						for fileName := range directoryMetaData[userInputVar] {
							// download file
							ClearScreen()
							dmsggetWrapper(serverPublicKey, MainDownloadsLoc, getPresentWorkingDirectory()+fileName, "", true)
						}
					}
				}

			}

		} else {
			break
		}
	}
	goto ServerIndexMenu

ExitLoop:
}
