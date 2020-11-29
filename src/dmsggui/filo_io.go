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
	downloadBrowserIndex = 0
}

func clearCacheConfig() {
	configFile := generateConfigAbsFilePath()
	_, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
	}
}

func generateServerIndexAbsPath(serverPublicKey string) string {
	indexPath := indexDownloadLoc + "/index." + serverPublicKey

	return indexPath
}

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
func InitProgramWorkingDir() string {
	_programCurrentWorkingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error obtaining program's working dir.")
	}
	fmt.Println("Working dir is:", _programCurrentWorkingDir)
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
