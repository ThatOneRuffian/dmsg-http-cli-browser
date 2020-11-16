package dmsg_gui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//ConfigFileHomePath stores the path, in the user's home dir, where the server address cache is saved
var ConfigFileHomePath string = "/.config"

//MainDownloadsLoc is the location where downloads are stored
var MainDownloadsLoc string

//DownloadListLength is how many entries are displayed on-screen at a time
var DownloadListLength int = 10

//IndexDownloadLoc is where the active server's index is downloaded
var IndexDownloadLoc string = os.TempDir()

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

func clearCacheConfig() {
	configFile := generateConfigAbsFilePath()
	file, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
		fmt.Println(file)
	}

}

func generateServerIndexAbsPath(serverPublicKey string) string {
	indexPath := "/tmp/index." + serverPublicKey

	return indexPath
}

func generateConfigAbsDirPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}
	return homeDir + ConfigFileHomePath
}

func generateConfigAbsFilePath() string {

	configAbsFilePath := fmt.Sprintf("%s/dmsg-http-browser.config", generateConfigAbsDirPath())

	return configAbsFilePath
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

func appendToConfig(friendlyName string, serverPublicKey string) {

	friendlyName = removeNewline(friendlyName)
	serverPublicKey = removeNewline(serverPublicKey)
	rawData := friendlyName + ";" + serverPublicKey + string('\n')

	dataToWrite := []byte(rawData)

	configDirPath := generateConfigAbsDirPath()
	err := os.Chdir(configDirPath)

	if os.IsNotExist(err) {
		dirErr := os.MkdirAll(configDirPath, 0700)
		if dirErr != nil {
			errorInfo := fmt.Sprintf("There was an error writng the config file to: %s\n%s", configDirPath, dirErr)
			log.Fatal(errorInfo)
		}
	}

	f, err := os.OpenFile(generateConfigAbsFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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

func InitDownloadsFolder() {
	tmpString, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error initializing downloads location")
	}
	MainDownloadsLoc = tmpString + "/Downloads"
}

func LoadCache() bool {
	returnBool := true
	file, err := os.Open(generateConfigAbsFilePath())
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

func deleteServerIndex(indexToDelete int) {
	clearCacheConfig()

	for index := 0; index < len(SavedServers); index++ {
		if index == indexToDelete-1 {
			continue
		} else {
			appendToConfig(SavedServers[index][0], SavedServers[index][1])
		}
	}

	LoadCache()
}
