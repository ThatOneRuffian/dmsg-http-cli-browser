package dmsg_gui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// MainDownloadsLoc is the location where downloads are stored
var MainDownloadsLoc string

//DownloadListLength is how many entries are displayed on-screen at a time
var DownloadListLength int = 20

//IndexDownloadLoc is where the active server's index is downloaded

func LoadServerIndex(serverPublicKey string) bool {
	returnBool := true
	file, err := os.Open(GenerateServerIndexAbsPath(serverPublicKey))
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

	ParseServerIndex(&file)
	return returnBool
}
func ClearFile(filename string) {
	os.Remove(filename)
}

func ParseServerIndex(file **os.File) {
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

func ClearCacheConfig() {
	configFile := generateConfigAbsPath()
	file, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
		fmt.Println(file)
	}

}

func GenerateServerIndexAbsPath(serverPublicKey string) string {
	indexPath := "/tmp/index." + serverPublicKey

	return indexPath
}

func generateConfigAbsPath() string {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	configPath := "/.config/dmsg-http-browser.config"

	return homeDir + configPath
}

func ParseConfigFile(file **os.File) {
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

func AppendToConfig(friendlyName string, serverPublicKey string) {

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

func InitDownloadsFolder() {
	tmpString, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error initializing downloads location")
	}
	MainDownloadsLoc = tmpString + "/Downloads"
}

func LoadCache() bool {
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
	ParseConfigFile(&file)
	return returnBool
}
