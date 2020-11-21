package dmsg_gui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// navPtr points to the current directory object being browsed
var navPtr *directory = nil

//DownloadBrowserIndex stores the current page number of the server's download list
var DownloadBrowserIndex int = 0

//var rootDir Directory
var rootDir directory = directory{
	files:     make(map[string]float64),
	parentDir: nil,
	dirName:   "/",
	subDirs:   make(map[string]*directory),
}

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

	parseServerIndex2(&file)
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

func parseServerIndex2(file **os.File) {
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
				populateFileSystem(inputRow)

			} else {
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

func populateFileSystem(fullFilePath string) {
	//scraping file name, size, and file structure
	//passes these items to create the filesystem
	var fileNameAndSize = make(map[string]int)
	splitString := strings.Split(fullFilePath, ";")
	fullPath := splitString[0]

	fileSize, err := strconv.Atoi(splitString[1])
	if err != nil {
		fmt.Println("Unable to convert filesize into int while populating directory: ", err)
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

			currentDirPtr = &newDirectory
		}

	}
}
