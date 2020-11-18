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
	fmt.Print("(R to Refresh Server Index, E to Exit Server File Browser, Q to quit): ")
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
	case "N":
		if DownloadBrowserIndex < ServerPageCountMax-1 {

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
		if userInt >= 1 && userInt <= len(directoryMetaData) {

			//determine if the key is a dir or a file
			for index := range directoryMetaData[userInt] {

				if index == ".." && navPtr.parentDir != nil {
					navPtr = navPtr.parentDir
				} else {
					if directoryMetaData[userInt][index] {

						navPtr = navPtr.subDirs[index]
					} else {
						for fileName := range directoryMetaData[userInt] {
							// download file
							fmt.Println("File", getPresentWorkingDirectory()+fileName)
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
