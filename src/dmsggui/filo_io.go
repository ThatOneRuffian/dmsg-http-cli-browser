package dmsggui

import (
	"fmt"
	"log"
	"os"
)

//configFileHomePath stores the path, in the user's home dir, where the server address cache is saved
var configFileHomePath string = "/.config"

//MainDownloadsLoc is the location where downloads are stored
var MainDownloadsLoc string

//indexDownloadLoc is where the active server's index is downloaded
var indexDownloadLoc string = os.TempDir()

var programCurrentWorkingDir string = ""

func InitProgramWorkingDir() string {
	_programCurrentWorkingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error obtaining program's working dir.")
	}
	programCurrentWorkingDir = _programCurrentWorkingDir
	return programCurrentWorkingDir
}

func InitDownloadsFolder(customDir string) string {
	if len(customDir) == 0 {
		tmpString, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error initializing downloads location")
		}
		//check if download path exist
		MainDownloadsLoc = tmpString + "/Downloads"

	} else {
		MainDownloadsLoc = customDir
	}

	// change dir back to program working dir for relative paths
	if err := os.Chdir(programCurrentWorkingDir); err != nil {
		fmt.Println("Error navigating program root directory.")
		os.Exit(1)
	}

	dirNotFoundErr := os.Chdir(MainDownloadsLoc)
	// If download location is not found...
	if os.IsNotExist(dirNotFoundErr) {
		// Attempt to create dir if it does not exist
		mkdirErr := os.Mkdir(MainDownloadsLoc, 0744)

		if mkdirErr != nil {
			// could not create downloads location panic
			errorMsg := fmt.Sprintf("Unable to initialize downloads location. Make sure you have permission to write to: %s", MainDownloadsLoc)
			panic(errorMsg)
		}
	}
	return MainDownloadsLoc
}

// server index
func loadServerIndex(serverPublicKey string) bool {
	returnBool := true
	file, err := os.Open(generateServerIndexAbsPath(serverPublicKey))
	defer file.Close()
	defer func() {
		if err := recover(); err != nil {
			currentServerIndexContents = make(map[int][2]string)
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

func clearServerIndexFile(serverPublicKey string) {
	serverCacheLoc := indexDownloadLoc + "/index." + serverPublicKey
	os.Remove(serverCacheLoc)
	resetDownLoadPageIndex()
}

func generateServerIndexAbsPath(serverPublicKey string) string {
	indexPath := indexDownloadLoc + "/index." + serverPublicKey

	return indexPath
}

// server list cache/config
func generateConfigAbsDirPath() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}
	return homeDir + configFileHomePath
}

func generateConfigAbsFilePath() string {

	configAbsFilePath := fmt.Sprintf("%s/dmsg-http-browser.config", generateConfigAbsDirPath())

	return configAbsFilePath
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

func appendToConfig(friendlyName string, serverPublicKey string) {

	friendlyName = stripIllegalChars(friendlyName)
	serverPublicKey = stripIllegalChars(serverPublicKey)
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

func deleteServerIndex(indexToDelete int) {
	clearCacheConfig()
	//filter out entry and rewrite to file
	for index := 0; index < len(SavedServers); index++ {
		if index == indexToDelete-1 {
			continue
		} else {
			appendToConfig(SavedServers[index][0], SavedServers[index][1])
		}
	}

	LoadCache()
}

func refreshServerIndex(serverPublicKey string, clearCache bool) {
	ClearScreen()
	resetDownLoadPageIndex()
	if clearCache {
		fmt.Println("Downloading Server Index...")
		clearServerIndexFile(serverPublicKey)
		dmsggetWrapper(serverPublicKey, indexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
	if loadServerIndex(serverPublicKey) {

	} else {
		ClearScreen()
		fmt.Println("Downloading Server Index...")
		clearServerIndexFile(serverPublicKey)
		dmsggetWrapper(serverPublicKey, indexDownloadLoc, "index", "index."+serverPublicKey, false)
	}
	assembleFileStructure(serverPublicKey)
}

func clearCacheConfig() {
	configFile := generateConfigAbsFilePath()
	_, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
	}
}

func getDownloadFileSize(fileName string) float64 {
	var returnValue float64
	fileInfo, err := os.Stat(MainDownloadsLoc + "/" + fileName)

	if err != nil {
		fmt.Println("Error:", err)
		returnValue = 0.0
	} else {
		returnValue = float64(fileInfo.Size())
	}
	return returnValue
}
