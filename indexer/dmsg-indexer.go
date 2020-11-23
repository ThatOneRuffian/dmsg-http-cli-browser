package main

import (
	"dmsggui"
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

func main() {
	//default index interval
	sleepInterval := time.Second
	intervalInputString := ""

	flag.StringVar(&indexPath, "d", ".", "Specify directory to be indexed.")
	flag.StringVar(&intervalInputString, "t", "30", "Specify the index interval in seconds.")

	flag.Parse()

	// set time interval
	sleepIntervalInput, timeParseErr := strconv.Atoi(intervalInputString)
	if timeParseErr != nil {
		dmsggui.ClearScreen()
		fmt.Println("Error interpreting time interval input. Enter the number of seconds between each index as an integer.")
		printUseage()
		os.Exit(0)
	} else if sleepIntervalInput > 0 {
		// if time input okay, then assign to interval
		sleepInterval = time.Duration(sleepIntervalInput) * time.Second
	}

	//set index path
	if indexPath != "" {
		fmt.Println("Setting index path to: ", indexPath)
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Unable to get current working directory.", err)
		}

		pathErr := os.Chdir(indexPath)

		if os.IsNotExist(pathErr) {
			//if index dir not found attempt to create
			createDirErr := os.MkdirAll(indexPath, 0744)
			if createDirErr != nil {
				fmt.Println("Unable to create directory:", indexPath)
				panic(createDirErr)
			}
		}

		cwdError := os.Chdir(currentDir)

		if cwdError != nil {
			fmt.Println("Error changing directory back.")
		}
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

			if directory[entry][0] != "index" {
				appendToIndex(directory[entry])
			}
		}
		time.Sleep(sleepInterval)
	}
}

func printUseage() {
	fmt.Println("Program usage:  indexer [index_interval_in_seconds - (Default=30s)]")
	fmt.Println()
}

// filePathWalk will list all absolute file paths and their sizes
func filePathWalk(root string) ([][2]string, error) {
	var files [][2]string
	var appendData [2]string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			fileInfo, err := os.Stat(path)
			if err != nil {
				fmt.Println(err)
			}
			fileSize := fmt.Sprint(fileInfo.Size())
			appendData[0] = path
			appendData[1] = string(fileSize)

			files = append(files, appendData)
		}
		return nil
	})
	return files, err
}

func clearCurrentIndex() {
	configFile := indexPath + "index"
	file, err := os.Create(configFile)

	if err != nil {

		fmt.Println("Error opening", configFile)

		fmt.Println(err)
		fmt.Println(file)
	}

}

func appendToIndex(filename [2]string) {
	filename[0] = removeNewline(filename[0])
	rawData := fmt.Sprintf("%s;%s\n", filename[0], filename[1])

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

func removeNewline(userInput string) string {
	return strings.TrimRight(userInput, "\n")
}
