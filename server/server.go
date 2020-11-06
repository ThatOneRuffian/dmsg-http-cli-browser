package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	//default interval
	sleepInterval := time.Duration(10) * time.Second

	//parse program arguments
	programArguments := ""
	if len(os.Args) > 1 {
		programArguments = os.Args[1]
	}
	intergerValue, err := strconv.Atoi(programArguments)
	if err != nil {
	} else {
		sleepInterval = time.Duration(intergerValue) * time.Second
	}

	// enter main file monitor loop
	for true {
		directory, err := FilePathWalk(".")
		if err != nil {
			fmt.Println("An error occured while reading the directory.")
		}
		clearCurrentIndex()
		for entry := range directory {

			if directory[entry][0] != "index" {
				fmt.Println(directory[entry])
				appendToIndex(directory[entry])
			}
		}
		time.Sleep(sleepInterval)
	}
}

// FilePathWalk will list all files and sub directories in a path
func FilePathWalk(root string) ([][2]string, error) {
	var files [][2]string
	var appendData [2]string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileInfo, err := os.Stat(path)
			if err != nil {
				fmt.Println(err)
			}
			fileSize := fmt.Sprint(fileInfo.Size())
			fmt.Println(fileSize) // bytes
			appendData[0] = path
			appendData[1] = string(fileSize)

			files = append(files, appendData)
		}
		return nil
	})
	return files, err
}

func clearCurrentIndex() {
	configFile := "./index"
	file, err := os.Create(configFile)

	if err != nil {
		fmt.Println(err)
		fmt.Println(file)
	}

}

func appendToIndex(filename [2]string) {
	filename[0] = removeNewline(filename[0])
	rawData := fmt.Sprintf("%s;%s\n", filename[0], filename[1])

	dataToWrite := []byte(rawData)

	f, err := os.OpenFile("./index", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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
