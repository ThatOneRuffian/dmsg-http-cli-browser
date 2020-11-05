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
	programArguments := os.Args[0]
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

			if directory[entry] != "index" {
				fmt.Println(directory[entry])
				appendToIndex(directory[entry])
			}
		}
		time.Sleep(sleepInterval)
	}
}

// FilePathWalk will list all files and sub directories in a path
func FilePathWalk(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
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

func appendToIndex(filename string) {
	filename = removeNewline(filename)
	rawData := filename + string('\n')

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
