package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//SavedServers stores server cache - initalized on loadCache
var SavedServers map[int][2]string

//CurrentServerIndex will store the parsed server index values
var CurrentServerIndex map[int][2]string

//IndexDownloadLoc is where the active server's index is downloaded
var IndexDownloadLoc string = "/tmp/"

// MainDownloadsLoc is the location where downloads are stored
var MainDownloadsLoc string

//DownloadListLength is how many entries are displayed on-screen at a time
var DownloadListLength int = 20

//DownloadBrowserIndex stores the current page of the download browsers
var DownloadBrowserIndex int = 1

func main() {
	clearScreen()
	initDownloadsFolder()
	// if config not found then run the first launch wizard
	if !loadCache() {
		firstRunWizard()
	}

	for true {
		userChoice := menuHandler()
		serverIndexMenuHandler(userChoice)
	}
}

// =========== User Interface ===========
func menuHandler() string {
	serverPublicKey := ""
	consoleInput := bufio.NewReader(os.Stdin)
ServerMenu:

	renderServerBrowser()

	fmt.Print("(Press A to Add server, D to Delete a server, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(removeNewline(userChoice))
	switch userChoice {
	case "Q":
		os.Exit(1)
	case "A":
		clearScreen()
		addServer()
		loadCache()
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

func serverIndexMenuHandler(serverPublicKey string) {
	consoleInput := bufio.NewReader(os.Stdin)
ServerIndexMenu:

	renderServerIndexBrowser()

	fmt.Print("(R to Refresh Server Index, E to Exit Server File Browser, Q to quit): ")
	userChoice, _ := consoleInput.ReadString('\n')
	userChoice = strings.ToUpper(removeNewline(userChoice))
	switch userChoice {
	case "Q":
		os.Exit(1)
	case "E":
		goto ExitLoop
	case "B":
		if DownloadBrowserIndex > 1 {
			DownloadBrowserIndex--
		}

		goto ServerIndexMenu
	case "N":
		if DownloadBrowserIndex*DownloadListLength < len(CurrentServerIndex) {
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
			clearScreen()

			dmsggetWrapper(serverPublicKey, MainDownloadsLoc, filenameDownload, "", true)
		} else {
			break
		}
	}
	goto ServerIndexMenu

ExitLoop:
}

func renderServerBrowser() {
	pageStatus := fmt.Sprintf("page (%d / %d)", 1, 20)
	divider := "----------------------"
	clearScreen()

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

	pageCountMax := len(CurrentServerIndex) / DownloadListLength
	pageRemainder := len(CurrentServerIndex) % DownloadListLength
	if pageRemainder > 0 {
		pageCountMax++
	}
	pageStatus := fmt.Sprintf("page (%d / %d)", DownloadBrowserIndex, pageCountMax)
	divider := "------------------------------------------"
	clearScreen()
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
func initDownloadsFolder() {
	tmpString, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error initializing downloads location")
	}
	MainDownloadsLoc = tmpString + "/Downloads"
}
func refreshServerIndex(serverPublicKey string, clearCache bool) {
	clearScreen()

	if clearCache {
		clearServerIndexFile(serverPublicKey)
	}

	fmt.Println("Downloading Server Index...")
	dmsggetWrapper(serverPublicKey, IndexDownloadLoc, "index", "index."+serverPublicKey, false)
	loadServerIndex(serverPublicKey)
}

func generateConfigAbsPath() string {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	configPath := "/.config/dmsg-http-browser.config"

	return homeDir + configPath
}

func generateServerIndexAbsPath(serverPublicKey string) string {
	indexPath := "/tmp/index." + serverPublicKey

	return indexPath
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func dmsggetWrapper(publicKey string, downloadLoc string, file string, alternateFileName string, stdOutput bool) bool {
	downloadInfo := ""

	if strings.Contains(file, "/") {
		fileName := strings.Split(file, "/")

		fileName[0] = fileName[len(fileName)-1]
		alternateFileName = fileName[0]
		downloadInfo = fmt.Sprintf("Downloading %s to %s", string(alternateFileName), downloadLoc)

	} else {
		downloadInfo = fmt.Sprintf("Downloading %s to %s", file, downloadLoc)

	}
	if downloadLoc == MainDownloadsLoc {
		clearFile(downloadLoc + "/" + alternateFileName)
		clearFile(downloadLoc + "/" + file)
	}

	fmt.Println(alternateFileName)
	fetchString := fmt.Sprintf("dmsg://%s:80/%s", publicKey, file)
	returnValue := true
	stdOutLoc := os.Stdout
	fmt.Println(downloadInfo)
	if !stdOutput {
		nullFile, err := os.OpenFile("/dev/null", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening /dev/null for writing")
		}
		defer nullFile.Close()
		stdOutLoc = nullFile
	}
	dmsggetPath, err := exec.LookPath("dmsgget")
	if err != nil {
		fmt.Println(err)
	}

	dmsggetCmd := &exec.Cmd{
		Path:   dmsggetPath,
		Args:   []string{dmsggetPath, "-O", downloadLoc + "/" + alternateFileName, fetchString},
		Stdout: stdOutLoc,
		Stderr: os.Stderr,
	}
	if err := dmsggetCmd.Run(); err != nil {
		fmt.Println("There was an error fetching the file", err)
		// file exists?
		returnValue = false
	}
	return returnValue
}

// =========== File I/O ===========
func loadServerIndex(serverPublicKey string) bool {
	returnBool := true
	file, err := os.Open(generateServerIndexAbsPath(serverPublicKey))
	defer file.Close()
	defer func() {
		if err := recover(); err != nil {
			CurrentServerIndex = make(map[int][2]string)
		}
	}()

	if err != nil {
		panic(err.Error())
	}

	fileStats, err := file.Stat()

	if err != nil {
		panic(err.Error())
	}

	if fileStats.Size() == 0 {
		returnBool = false
	}

	parseServerIndex(&file)
	return returnBool
}
func clearFile(filename string) {
	os.Remove(filename)
}

func clearServerIndexFile(serverPublicKey string) {
	serverCacheLoc := "/tmp/index." + serverPublicKey
	os.Remove(serverCacheLoc)
}

func clearCacheConfig() {
	configFile := generateConfigAbsPath()
	file, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
		fmt.Println(file)
	}

}

func appendToConfig(friendlyName string, serverPublicKey string) {

	friendlyName = removeNewline(friendlyName)
	serverPublicKey = removeNewline(serverPublicKey)
	rawData := friendlyName + ";" + serverPublicKey + string('\n')

	dataToWrite := []byte(rawData)

	f, err := os.OpenFile(generateConfigAbsPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write(dataToWrite); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func loadCache() bool {
	returnBool := true
	file, err := os.Open(generateConfigAbsPath())
	defer file.Close()
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	if err != nil {
		panic(err.Error())
	}

	fileStats, err := file.Stat()

	if err != nil {
		panic(err.Error())
	}

	if fileStats.Size() == 0 {
		returnBool = false
	}

	// load up map values
	parseConfigFile(&file)
	return returnBool
}

func parseConfigFile(file **os.File) {
	savedServers := make(map[int][2]string)
	friendlyNameIndex := 0
	serverPubKeyIndex := 1
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("Error parsing configuration file.")
		}
	}()

	fileScan := bufio.NewScanner(*file)

	i := 0
	for fileScan.Scan() {
		splitStringArray := [2]string{"", ""}
		tmpString := fileScan.Text()
		tmpSplitString := strings.Split(tmpString, ";")
		splitStringArray[0] = tmpSplitString[friendlyNameIndex]
		splitStringArray[1] = tmpSplitString[serverPubKeyIndex]
		savedServers[i] = splitStringArray
		i++
	}
	SavedServers = savedServers
}

func parseServerIndex(file **os.File) {
	currentServerIndex := make(map[int][2]string)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("Error parsing server index file.")
		}
	}()

	fileScan := bufio.NewScanner(*file)

	i := 0
	for fileScan.Scan() {
		var swapVar [2]string
		inputRow := fileScan.Text()
		if sepIndex := strings.Index(inputRow, ";"); sepIndex != -1 {
			//convert slice back into array
			parsedString := strings.Split(inputRow, ";")
			swapVar[0] = parsedString[0]
			swapVar[1] = parsedString[1]
			currentServerIndex[i] = swapVar
			i++
		}
	}

	CurrentServerIndex = currentServerIndex
}

// =========== Wizards ===========
func browseNow(serverPublicKey string) {
	consoleInput := bufio.NewReader(os.Stdin)

Browse:
	fmt.Print("Would you like to browse this server now? (Y/N): ")
	userAnswer, _ := consoleInput.ReadString('\n')

	switch formattedInput := strings.ToUpper(removeNewline(userAnswer)); formattedInput {
	case "Y":
		refreshServerIndex(serverPublicKey, true)
		serverIndexMenuHandler(serverPublicKey)
		//load server index
	case "N":
		// continue to main menu
	default:
		goto Browse
	}
}
func firstRunWizard() {
	fmt.Println("It looks like this is your frist time running the dmsg-http CLI browser.")
	serverPublicKey := addServer()
	browseNow(serverPublicKey)
	loadCache()
}

func addServer() string {
	keyLength := 66
	consoleInput := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the public key for the dmsg-http server you want to add: ")

PubKeyInput:
	publicKey, _ := consoleInput.ReadString('\n')
	publicKey = removeNewline(publicKey)

	if len(publicKey) == keyLength {
		fmt.Print("Add a friendly name to this public key (default: [public_key]): ")
		friendlyName, _ := consoleInput.ReadString('\n')
		friendlyName = removeNewline(friendlyName)
		friendlyName = removeSemiColon(friendlyName)
		if len(friendlyName) == 0 {
			appendToConfig(publicKey, publicKey)
		} else {
			appendToConfig(friendlyName, publicKey)
		}
		fmt.Println("Entry cached.")

	} else {
		errorInfo := fmt.Sprintf("Provided key has length of %d. Expected length of %d.", len(publicKey), keyLength)
		fmt.Println(errorInfo)
		fmt.Print("Invalid key length please enter public key again: ")
		goto PubKeyInput
	}
	return publicKey
}

func deleteServerIndex(indexToDelete int) {
	clearCacheConfig()

	for index := 0; index < len(SavedServers); index++ {
		if index == indexToDelete-1 {
			continue
		} else {
			appendToConfig(SavedServers[index][0], SavedServers[index][1])
		}
	}

	loadCache()
}

func deleteServerWizard() {
DeletePrompt:
	fmt.Print("Which server do you want to delete? (Enter C to Cancel): ")
	consoleInputWhichServer := bufio.NewReader(os.Stdin)
	userDelete, _ := consoleInputWhichServer.ReadString('\n')
	userDelete = strings.ToUpper(removeNewline(userDelete))

	switch userDelete {
	case "C":
		goto ExitLoop
	default:
		userInt, err := strconv.Atoi(userDelete)
		if err != nil {
			break
		}
	ConfirmDelete:
		if userInt >= 1 && userInt <= len(SavedServers) {
			deleteConfirmPrompt := fmt.Sprintf("Are you sure you want to delete (Y/N)? { %s | %s }", SavedServers[userInt-1][0], SavedServers[userInt-1][1])
			fmt.Println(deleteConfirmPrompt)
			deleteConfirmInput := bufio.NewReader(os.Stdin)

			deleteConfirm, _ := deleteConfirmInput.ReadString('\n')
			deleteConfirm = strings.ToUpper(removeNewline(deleteConfirm))
			//deleteIndex, err := strconv.Atoi(deleteConfirm)

			switch deleteConfirm {
			case "Y":

				deleteServerIndex(userInt)
				goto ExitLoop
			case "N":
				goto ExitLoop
			default:
				goto ConfirmDelete
			}
		} else {
			break
		}
	}
	goto DeletePrompt
ExitLoop:
}

// =========== String formatting functions ===========

func removeNewline(userInput string) string {
	return strings.TrimRight(userInput, "\n")
}

func removeSemiColon(stringToScan string) string {
	semiColonCode := byte(59)
	spaceBarCode := byte(32)
	tmpByteString := []byte(stringToScan)
	for i, v := range tmpByteString {
		if v == semiColonCode {
			tmpByteString[i] = spaceBarCode
		}
	}
	return string(tmpByteString)
}
