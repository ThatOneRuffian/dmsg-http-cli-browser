package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var indexPath string
var fileFilters []string

func main() {
	//default index interval
	sleepInterval := time.Second
	intervalInputString := ""
	filterFile := ""
	flag.StringVar(&indexPath, "d", ".", "Specify directory to be indexed.")
	flag.StringVar(&intervalInputString, "t", "30", "Specify the index interval in seconds.")
	flag.StringVar(&filterFile, "f", "", "Specify a txt file where each line item is a keyword filter to keep files containing those keywords from being indexed.")

	flag.Parse()

	// set time interval
	sleepIntervalInput, timeParseErr := strconv.Atoi(intervalInputString)
	if timeParseErr != nil {
		fmt.Println("Error interpreting time interval input. Enter the number of seconds between each indexing as an integer.")
		os.Exit(0)
	} else if sleepIntervalInput > 0 {
		// if time input okay, then assign to interval
		sleepInterval = time.Duration(sleepIntervalInput) * time.Second
	}

	//set index path
	if indexPath != "" {
		parseIndexPathInput(&indexPath)
	}
	//Init filter file
	if filterFile != "" {
		parseFilterFile(&filterFile)
	}

	fmt.Println("Indexing with an interval of:", sleepInterval)
	// enter main file monitor loop
	for true {
		directory, err := filePathWalk(indexPath)
		if err != nil {
			fmt.Println("An error occured while reading the directory.")
			fmt.Println(err)
		}
		clearCurrentIndex() //todo add diff check before rewriting index?
		for entry := range directory {

			writeToIndex(directory[entry])

		}
		time.Sleep(sleepInterval)
	}
}

func parseFilterFile(filterFileLocation *string) {
	file, err := os.Open(*filterFileLocation)

	if os.IsNotExist(err) {
		fmt.Println(err.Error())
		log.Fatal("Unable to open filter file")
	}

	*filterFileLocation = normalizePath(*filterFileLocation)
	if strings.Contains(*filterFileLocation, "/") {
		//add the rest of the filters from file
		parseConfigFile(&file)
	}

}

func normalizePath(filePath string) string {
	tmpByteString := []byte(filePath)
	tmpString := ""
	var forwardSlashChar byte = 47
	var backSlashChar byte = 92
	for _, value := range tmpByteString {
		byteChar := value
		if byteChar == backSlashChar {
			byteChar = forwardSlashChar
		}
		tmpString += string(byteChar)
	}
	return tmpString
}

func parseConfigFile(file **os.File) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("Error parsing filter file.")
		}
	}()

	fileScan := bufio.NewScanner(*file)

	for fileScan.Scan() {
		tmpString := fileScan.Text()
		fileFilters = append(fileFilters, tmpString)
	}
}

func parseIndexPathInput(indexPath *string) {

	pathByteArray := []byte(*indexPath)
	lastByteChar := pathByteArray[len(pathByteArray)-1]
	const forwardSlash byte = 92
	const backSlash byte = 47
	//append "/" if missing from provided dir
	if lastByteChar != forwardSlash && lastByteChar != backSlash {
		*indexPath = *indexPath + "/"
	}
	fmt.Println("Setting index path to: ", *indexPath)
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to get current working directory.", err)
	}
	// try changing directory to test if provided directory exist
	pathErr := os.Chdir(*indexPath)

	if os.IsNotExist(pathErr) {
		//if provided dir not found, then attempt to create
		createDirErr := os.MkdirAll(*indexPath, 0744)
		if createDirErr != nil {
			fmt.Println("Unable to create directory:", indexPath)
			panic(createDirErr)
		}
	}
	//change  directory back to starting dir
	cwdError := os.Chdir(currentDir)

	if cwdError != nil {
		fmt.Println("Error changing directory back.")
	}
}

// filePathWalk will list all absolute file paths and their sizes
func filePathWalk(root string) ([][2]string, error) {
	var files [][2]string
	var appendData [2]string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			fileInfo, err := os.Stat(path)
			if err != nil {
				fmt.Println("Error walking dir", err)
			}
			fileSize := fmt.Sprint(fileInfo.Size())
			appendData[0] = path
			appendData[1] = string(fileSize)

			//filter file
			if len(fileFilters) > 0 {
				tmpPath := normalizePath(path)
				tmpFileName := strings.Split(tmpPath, "/")

				for _, filter := range fileFilters {
					if strings.Contains(tmpFileName[len(tmpFileName)-1], filter) {
						goto SkipFile
					}
				}
			}
			files = append(files, appendData)

		SkipFile:
		}
		return nil
	})
	return files, err
}

func clearCurrentIndex() {
	configFile := indexPath + "index"
	err := os.Remove(configFile)

	if os.IsNotExist(err) {
		// no index file or empty directory nothing to index
	} else if err != nil {
		fmt.Println("Error removing index file; ", err)
	}

}

func writeToIndex(fileInfo [2]string) {
	fileInfo[0] = normalizePath(removeNewline(fileInfo[0]))

	if len(fileInfo[0]) > 0 {
		if strings.Split(fileInfo[0], "/")[len(fileInfo)-1] != "index" {
			rawData := fmt.Sprintf("%s;%s\n", fileInfo[0], fileInfo[1])
			dataToWrite := []byte(rawData)

			f, err := os.OpenFile(indexPath+"/index", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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
	}
}

func removeNewline(userInput string) string {
	return strings.TrimRight(userInput, "\n")
}
