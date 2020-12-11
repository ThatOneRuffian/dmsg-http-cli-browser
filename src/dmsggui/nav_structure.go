package dmsggui

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// navPtr points to the current directory object being browsed
var navPtr *directory = nil

const downloadBrowserStartIndex = 0

//downloadBrowserIndex stores the current page number - 1 of the server's download list
var downloadBrowserIndex int = downloadBrowserStartIndex

//downloadBrowserIndex stores the current page number - 1 of the main menu
var mainMenuBrowserIndex int = 0

var downloadQueueIndex int = 0

var rootDir directory

var subDir directory

type directory struct {
	files     map[string]float64
	subDirs   map[string]*directory
	parentDir *directory
	dirName   string
}

func assembleFileStructure(serverPublicKey string) {
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

	if err != nil {
		panic(err.Error())
	}
	parseServerIndex(&file)
}

func initRootDir() {
	// reinitialize root dir
	rootDir = directory{
		files:     make(map[string]float64),
		parentDir: nil,
		dirName:   "/",
		subDirs:   make(map[string]*directory),
	}
}

func populateFileSystem(fullFilePath string) {
	//scraping file name, size, and file structure
	//passes these items to create the filesystem
	var fileNameAndSize = make(map[string]int)
	splitString := strings.Split(fullFilePath, ";")
	fullPath := splitString[0]

	fileSize, err := strconv.Atoi(splitString[1])
	if err != nil {
		//fmt.Println("Unable to convert filesize into int while populating directory: ", err) //save for logs?
	}
	fileNameSlice := strings.Split(fullPath, "/")[len(strings.Split(fullPath, "/"))-1:]
	//converting filename from slice into string
	fileNameString := strings.Join(fileNameSlice, "")
	fileNameAndSize[fileNameString] = fileSize
	dirStructure := strings.Split(fullPath, "/")[:len(strings.Split(fullPath, "/"))-1] // strip filename from path
	//create dir structure
	createDirPath(dirStructure)
	//insert file into dir
	insertFileIntoDir(dirStructure, fileNameString, fileSize)
}

func getPresentWorkingDirectory() string {
	var workingDir string
	var dirs []string

	if navPtr != nil {
		tmpPtr := navPtr

		for tmpPtr.parentDir != nil {
			dirs = append(dirs, tmpPtr.dirName)
			tmpPtr = tmpPtr.parentDir
		}
		for i := len(dirs) - 1; i >= 0; i-- {
			workingDir = workingDir + "/" + dirs[i]
		}
	}

	return workingDir + "/"
}

func parseServerIndex(file **os.File) {
	initRootDir()
	navPtr = &rootDir
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("Error parsing server index file.")
		}
	}()

	fileScan := bufio.NewScanner(*file)

	i := 0
	for fileScan.Scan() {
		inputRow := fileScan.Text()
		if sepIndex := strings.Index(inputRow, ";"); sepIndex != -1 {
			if len(strings.Split(inputRow, "/")) > 1 {
				//if the entry is a directory
				populateFileSystem(inputRow)

			} else {
				//if the entry is a file
				fileInfo := strings.Split(inputRow, ";")
				fileSize, err := strconv.Atoi(fileInfo[1])
				if err != nil {
					fmt.Println("Unable to convert filesize string into int")
				}
				rootDir.files[fileInfo[0]] = float64(fileSize)
				fmt.Println("root dir level file found: ", strings.Split(inputRow, ";")[0])
			}

			i++
		}
	}
}

func parseConfigFile(file **os.File) {
	_savedServers := make(map[int][2]string)
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
		_savedServers[i] = splitStringArray
		i++
	}
	SavedServers = _savedServers
}

func insertFileIntoDir(filePath []string, fileName string, fileSize int) {
	currentDirPtr := &rootDir

	for _, currentPathName := range filePath {
		_, ok := currentDirPtr.subDirs[currentPathName]
		// if the subdir exist then point to that
		if ok {
			currentDirPtr = currentDirPtr.subDirs[currentPathName]
			if currentPathName == filePath[len(filePath)-1] {
				// final directory reached
				currentDirPtr.files[fileName] = float64(fileSize)
			}
		} else {
			fmt.Println("Dir does not exist cannot create file")
			break
		}
	}
}

func createDirPath(fullDirPath []string) {
	//create the dir path for the passed string
	currentDirPtr := &rootDir
	navPtr = &rootDir
	for _, currentPathName := range fullDirPath {
		_, ok := currentDirPtr.subDirs[currentPathName]
		if ok {
			//if directory already present then move into
			currentDirPtr = currentDirPtr.subDirs[currentPathName]

		} else {
			//create file tree
			newDirectory := directory{
				files:     make(map[string]float64),
				subDirs:   make(map[string]*directory),
				parentDir: currentDirPtr,
				dirName:   currentPathName,
			}
			currentDirPtr.subDirs[currentPathName] = &newDirectory
			//set newly created dir as current dir
			currentDirPtr = &newDirectory
		}

	}
}

func resetDownLoadPageIndex() {
	downloadBrowserIndex = downloadBrowserStartIndex
}

func getCurrentDirMetaData() map[int]map[string]bool {
	var subDirKeys []string
	var fileNames []string
	returnValue := make(map[int]map[string]bool)
	swapDir := make(map[string]bool)
	lengthOfFilterList := len(currentDirFilter)
	//dump dir/file keys in current dir and sort A-Z
	if navPtr != &rootDir {
		//add key for directory back
		subDirKeys = append(subDirKeys, "..")
	}

	for key := range navPtr.subDirs {
		//append subdir keys
		subDirKeys = append(subDirKeys, key)
	}

	for key := range navPtr.files {
		//append file keys
		fileNames = append(fileNames, key)
	}

	sort.Strings(fileNames)
	sort.Strings(subDirKeys)

	//merge metadata
	for key, value := range subDirKeys {
		if lengthOfFilterList > 0 {
			if strings.Contains(strings.ToUpper(value), strings.ToUpper(string(currentDirFilter))) || strings.Contains(strings.ToUpper(value), "..") {
				swapDir[value] = true
				returnValue[key+1] = swapDir
				swapDir = make(map[string]bool)
			}
		} else if lengthOfFilterList == 0 {
			swapDir[value] = true
			returnValue[key+1] = swapDir
			swapDir = make(map[string]bool)
		}
	}

	for key, value := range fileNames {
		if strings.Contains(strings.ToUpper(value), strings.ToUpper(string(currentDirFilter))) {
			swapDir[value] = false
			returnValue[key+1+len(subDirKeys)] = swapDir
			swapDir = make(map[string]bool)
		} else if lengthOfFilterList == 0 {
			swapDir[value] = false
			returnValue[key+1+len(subDirKeys)] = swapDir
			swapDir = make(map[string]bool)
		}
	}
	return returnValue
}

func isMetaDataSorted(directoryMetaData map[int]map[string]bool) bool {
	metaDataLength := len(directoryMetaData)
	isSorted := true
	for index := range directoryMetaData {
		if index > metaDataLength {
			isSorted = false
		}
	}
	return isSorted
}
