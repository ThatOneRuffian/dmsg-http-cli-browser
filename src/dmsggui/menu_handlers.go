package dmsggui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var RetryAttemptsUserInput string

var DownloadLocationUserInput string

// Handler for the server picker, main menu
func ServerListMainMenu() string {
	serverPublicKey := ""
	consoleInput := bufio.NewReader(os.Stdin)
ServerMenu:
	currentDirFilter = ""
	renderServerBrowser()

	fmt.Print("(A to Add server, D to Delete a server, G to Goto page, V to View Downloads, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(stripIllegalChars(userChoice))
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
		userInput = strings.ToUpper(stripIllegalChars(userInput))
		pageNumber, err := strconv.Atoi(userInput)

		if err != nil {

		} else if pageNumber > 0 && pageNumber-1 < mainMenuPageCountMax {
			mainMenuBrowserIndex = pageNumber - 1
		}
	case "V":
		downloadQueuePageHandler()
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

// Handler for the server download page
func ServerIndexMenuHandler(serverPublicKey string) {
	navPtr = &rootDir
	assembleFileStructure(serverPublicKey)

ServerIndexMenu:
	directoryMetaData := renderServerDownloadList()
SearchLoop:
	consoleInput := bufio.NewReader(os.Stdin)
	fmt.Print("(R to Refresh Server Index, E to Exit Server File Browser, G to Goto page, S to Search Dir, V to View Downloads, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(stripIllegalChars(userChoice))
	switch userChoice {
	case "Q":
		ClearScreen()
		os.Exit(1)
	case "E":
		//clear server root dir
		initRootDir()
		goto ExitLoop
	case "B":
		if downloadBrowserIndex > 0 {
			downloadBrowserIndex--
		}

	case "F":
		resetDownLoadPageIndex()
	case "G":
		fmt.Print("Enter page number:")
		consoleReader := bufio.NewReader(os.Stdin)
		userInput, _ := consoleReader.ReadString('\n')
		userInput = strings.ToUpper(stripIllegalChars(userInput))
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
		initRootDir()
		navPtr = &rootDir
		currentDirFilter = ""
		resetDownLoadPageIndex()
		refreshServerIndex(serverPublicKey, true)
	case "S":
		fmt.Print("Search current directory for the following substring (X to clear current filter): ")
		consoleInput := bufio.NewReader(os.Stdin)
		inputQuery, _ := consoleInput.ReadString('\n')
		inputQuery = stripIllegalChars(inputQuery)
		//do not reset page number when there is no input
		if len(inputQuery) > 0 {
			currentDirFilter = inputQuery
			resetDownLoadPageIndex()
		}
		renderServerDownloadList()
		goto SearchLoop
	case "V":
		downloadQueuePageHandler()
	case "X":
		currentDirFilter = ""
		resetDownLoadPageIndex()

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
					resetDownLoadPageIndex()
					currentDirFilter = ""
				} else {
					// if object is a directory then navigate into
					if directoryMetaData[userInputVar][index] {
						currentDirFilter = ""
						navPtr = navPtr.subDirs[index]
					} else {
						// runs once
						for fileName := range directoryMetaData[userInputVar] {
							// download file
							ClearScreen()
							initMuxDownload(serverPublicKey, MainDownloadsLoc, fileName)
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

// Handler for download queue page
func downloadQueuePageHandler() {
DownloadQueueMenu:
	renderDownloadQueuePage()
SearchLoop:
	consoleInput := bufio.NewReader(os.Stdin)
	fmt.Print("(E to Exit download queue, G to Goto page, S to Search download queue, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(stripIllegalChars(userChoice))
	switch userChoice {
	case "Q":
		ClearScreen()
		os.Exit(1)
	case "E":
		goto ExitLoop

	case "B":

	case "F":
		resetDownLoadPageIndex()
	case "G":
		fmt.Print("Enter page number:")
		consoleReader := bufio.NewReader(os.Stdin)
		userInput, _ := consoleReader.ReadString('\n')
		userInput = strings.ToUpper(stripIllegalChars(userInput))
		pageNumber, err := strconv.Atoi(userInput)

		if err != nil {

		} else if pageNumber > 0 && pageNumber-1 < serverPageCountMax {
			downloadBrowserIndex = pageNumber - 1
		}

	case "L":

	case "N":

	case "S":
		fmt.Print("Search current directory for the following substring (X to clear current filter): ")
		consoleInput := bufio.NewReader(os.Stdin)
		inputQuery, _ := consoleInput.ReadString('\n')
		inputQuery = stripIllegalChars(inputQuery)
		//do not reset page number when there is no input
		if len(inputQuery) > 0 {
			currentDirFilter = inputQuery
			//reset download queue page number
		}
		//renderServerDownloadList()
		goto SearchLoop

	case "X":
		currentDirFilter = ""
		//reset download queue page number

	default:
		userInputVar, err := strconv.Atoi(userChoice)
		if err != nil {
			break
		}
		if userInputVar >= 1 && userInputVar <= len("test") {
			// attach to stdout?
		}

	}
	goto DownloadQueueMenu
ExitLoop:
}
