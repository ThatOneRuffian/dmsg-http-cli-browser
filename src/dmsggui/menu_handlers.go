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
		if mainMenuBrowserIndex > 0 {
			mainMenuBrowserIndex--
		}
	case "F":
		mainMenuBrowserIndex = 0
	case "L":
		mainMenuBrowserIndex = mainMenuPageCountMax - 1
	case "N":
		if mainMenuBrowserIndex < mainMenuPageCountMax-1 {

			mainMenuBrowserIndex++
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
			mainMenuBrowserIndex = pageNumber - 1
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
		if downloadBrowserIndex > 0 {
			downloadBrowserIndex--
		}

	case "F":
		downloadBrowserIndex = 0
	case "G":
		fmt.Print("Enter page number:")
		consoleReader := bufio.NewReader(os.Stdin)
		userInput, _ := consoleReader.ReadString('\n')
		userInput = strings.ToUpper(removeNewline(userInput))
		pageNumber, err := strconv.Atoi(userInput)

		if err != nil {

		} else if pageNumber > 0 && pageNumber-1 < serverPageCountMax {
			downloadBrowserIndex = pageNumber - 1
		}

	case "L":
		downloadBrowserIndex = serverPageCountMax - 1

	case "N":
		if downloadBrowserIndex < serverPageCountMax-1 {

			downloadBrowserIndex++
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
					downloadBrowserIndex = 0
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
