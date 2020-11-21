package dmsggui

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

	fmt.Print("(Press A to Add server, D to Delete a server, G to Goto page, Q to quit): ")
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
		if DownloadBrowserIndex > 0 {
			DownloadBrowserIndex--
		}
	case "F":
		DownloadBrowserIndex = 0
	case "L":
		DownloadBrowserIndex = mainMenuPageCountMax - 1
	case "N":
		if DownloadBrowserIndex < mainMenuPageCountMax-1 {

			DownloadBrowserIndex++
		}
	case "D":
		deleteServerWizard()
	case "G":
		fmt.Print("Enter page number:")
		consoleReader := bufio.NewReader(os.Stdin)
		userInput, _ := consoleReader.ReadString('\n')
		userInput = strings.ToUpper(removeNewline(userInput))
		pageNumber, err := strconv.Atoi(userInput)

		if err != nil {

		} else if pageNumber > 0 && pageNumber-1 < mainMenuPageCountMax {
			DownloadBrowserIndex = pageNumber - 1
		}
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

	case "F":
		DownloadBrowserIndex = 0
	case "G":
		fmt.Print("Enter page number:")
		consoleReader := bufio.NewReader(os.Stdin)
		pageNumber, _ := consoleReader.ReadString('\n')
		pageNumber = strings.ToUpper(removeNewline(pageNumber))
		pageNumber2, err := strconv.Atoi(pageNumber)

		if err != nil {

		} else if pageNumber2 > 0 && pageNumber2-1 < serverPageCountMax {
			DownloadBrowserIndex = pageNumber2 - 1
		}

	case "L":
		DownloadBrowserIndex = serverPageCountMax - 1

	case "N":
		if DownloadBrowserIndex < serverPageCountMax-1 {

			DownloadBrowserIndex++
		}

	case "R":
		refreshServerIndex(serverPublicKey, true)
		assembleFileStructure(serverPublicKey)

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
					DownloadBrowserIndex = 0
				} else {
					// if object is a directory then navigate into
					if directoryMetaData[userInputVar][index] {
						navPtr = navPtr.subDirs[index]
					} else {
						// runs once
						for fileName := range directoryMetaData[userInputVar] {
							// download file
							ClearScreen()
							dmsggetWrapper(serverPublicKey, mainDownloadsLoc, getPresentWorkingDirectory()+fileName, "", true)
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
