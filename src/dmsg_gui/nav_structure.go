package dmsg_gui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// navPtr points to the current directory object being browsed
var navPtr *Directory = nil

//var rootDir Directory
var rootDir Directory = Directory{
	parentDir: nil,
	dirName:   "/",
	subDirs:   make(map[string]*Directory),
}

var subDir Directory

type Directory struct {
	files     map[string]int
	subDirs   map[string]*Directory
	parentDir *Directory
	dirName   string
}

func CreateTest() {
	file, err := os.Open(GenerateServerIndexAbsPath("02da372fab7ebb28d77508629dfed017011b12608a8454ab76f1f4359c01519a8d"))
	defer file.Close()
	defer func() {
		if err := recover(); err != nil {
			CurrentServerIndex = make(map[int][2]string)
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
				//root dir level file
				fmt.Println("root dir level file found: ", strings.Split(inputRow, ";")[0])
			}

			i++
		}
	}
	navPtr = &rootDir
	fmt.Println(*navPtr)
	navPtr = rootDir.subDirs["Educational_Media"]
	navPtr = navPtr.subDirs["A+"]
	fmt.Println(navPtr)
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
	//append file to dir.files[]

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
			newDirectory := Directory{
				files:     make(map[string]int),
				subDirs:   make(map[string]*Directory),
				parentDir: currentDirPtr,
				dirName:   currentPathName,
			}
			currentDirPtr.subDirs[currentPathName] = &newDirectory

			currentDirPtr = &newDirectory
		}

	}
}
